package services

import (
	"github.com/pkg/errors"
	"gorm.io/gorm"
	"orderin-server/internal/dto"
	"orderin-server/internal/models"
	mcontext "orderin-server/pkg/common/context"
	"orderin-server/pkg/common/db/relation"
	"orderin-server/pkg/common/errs"
	"orderin-server/pkg/common/log"
	"orderin-server/pkg/common/service"
	"sort"
	"time"
)

type SysDict struct {
	service.Service
}

func (e SysDict) Insert(dict models.SysDict) error {
	var err error
	dbRecord := models.SysDict{}
	tx := e.Orm.Model(&dict).Where("name = ?", dict.Name).First(&dbRecord)
	err = tx.Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			err = e.Orm.Create(&dict).Error
			return err
		} else {
			log.ZError(e.Context, "db error", err)
			return err
		}
	}
	if tx.RowsAffected > 0 {
		if dbRecord.Deleted {
			dbRecord.Deleted = false
			dbRecord.DeletedAt = nil
			dbRecord.DeletedBy = nil
			dbRecord.Description = dict.Description

			currentUserId := mcontext.GetOpUserID(e.Context)
			dbRecord.CreatedBy = &currentUserId
			e.UpdateSelectiveById(dbRecord.ID, dbRecord)
		} else {
			//log.ZError(e.Context, "dict has been exist", err)
			err := errs.NewCodeError(errs.DuplicateKeyError, "字典已存在")
			return err
		}
	}
	err = e.Orm.Create(&dict).Error
	return err
}

func (e SysDict) GetByName(name string) (*models.SysDict, error) {
	var err error
	dict := models.SysDict{}
	tx := e.Orm.Model(&dict).Where("name = ?", name).First(&dict)
	err = tx.Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		} else {
			log.ZError(e.Context, "db error", err)
			return nil, err
		}
	}
	if tx.RowsAffected == 0 {
		return nil, err
	}
	return &dict, nil
}

func (e SysDict) UpdateSelectiveById(id int64, dict models.SysDict) error {
	dict.ID = id
	return e.Orm.Model(&dict).Updates(&dict).Error
}

func (e SysDict) Delete(id int64) {
	e.Orm.Model(&models.SysDict{}).Where("id = ?", id).Updates(map[string]interface{}{
		"deleted":    true,
		"deleted_by": mcontext.GetOpUserID(e.Context),
		"deleted_at": time.Now(),
	})
}

func (e SysDict) PageQuery(req *dto.SysDictPageQueryReq, list *[]models.SysDict, count *int64) error {
	err := e.Orm.
		Scopes(
			relation.MakeCondition(*req),
			relation.NotDeleted(),
			relation.Paginate(req.GetPageSize(), req.GetPageIndex()),
		).
		Find(list).Limit(-1).Offset(-1).
		Count(count).Error
	if err != nil {
		log.ZError(e.Context, "page query sys_dicts fail", err)
		return err
	}
	return nil
}

func (e SysDict) PageQueryItem(req *dto.SysDictItemPageQueryReq, list *[]models.SysDictItem, count *int64) error {
	err := e.Orm.
		Scopes(
			relation.MakeCondition(*req),
			relation.Paginate(req.GetPageSize(), req.GetPageIndex()),
		).
		Find(list).Limit(-1).Offset(-1).
		Count(count).Error
	if err != nil {
		log.ZError(e.Context, "page query sys_dict_items fail", err)
		return err
	}
	return nil
}

func (e SysDict) AddItem(dictItem models.SysDictItem) error {
	var err error
	var i int64
	err = e.Orm.Model(&dictItem).Where("dict_name = ? and item_value = ?", dictItem.DictName, dictItem.ItemValue).Count(&i).Error
	if err != nil {
		log.ZError(e.Context, "db error", err)
		return err
	}
	if i > 0 {
		//log.ZError(e.Context, "item has been exist", err)
		err = errs.NewCodeError(errs.DuplicateKeyError, "字典项已存在")
		return err
	}
	dictItem.Status = 1
	return e.Orm.Create(&dictItem).Error
}

func (e SysDict) UpdateItem(dictItem models.SysDictItem) error {
	return e.Orm.Model(&dictItem).Updates(&dictItem).Error
}

func (e SysDict) GetItems(name string) ([]models.SysDictItem, error) {
	var items []models.SysDictItem
	err := e.Orm.Where("dict_name = ? and status = 1", name).Order("sort asc").Find(&items).Error
	return items, err
}

func (e SysDict) BatchGetItems(names []string) (*map[string][]models.SysDictItem, error) {
	var list []models.SysDictItem
	err := e.Orm.Where("dict_name in (?) and status = 1", names).Find(&list).Error
	if err != nil {
		return nil, err
	}
	mapping := make(map[string][]models.SysDictItem, 0)
	for _, item := range list {
		itemList, exists := mapping[item.DictName]
		if exists {
			itemList = append(itemList, item)
			mapping[item.DictName] = itemList
		} else {
			mapping[item.DictName] = []models.SysDictItem{item}
		}
	}
	for _, name := range names {
		itemList, exists := mapping[name]
		if !exists {
			mapping[name] = make([]models.SysDictItem, 0)
		} else {
			sort.Slice(itemList, func(i, j int) bool {
				if i > 0 && j > 0 {
					return itemList[i].Sort < itemList[j].Sort
				} else {
					return itemList[i].ID < itemList[j].ID
				}
			})
		}
	}
	return &mapping, err
}

func (e SysDict) BatchAddDictAndItems(name string, label string, labels []string, values []string) {
	e.Insert(models.SysDict{Name: name, Label: label})
	for i, label := range labels {
		e.AddItem(models.SysDictItem{
			DictName:  name,
			ItemLabel: label,
			ItemValue: values[i],
			Sort:      i + 1,
		})
	}
}

func (e SysDict) GetItemById(id int64) (*models.SysDictItem, error) {
	var item models.SysDictItem
	result := e.Orm.First(&item, id)
	err := result.Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errs.NewCodeError(errs.RecordNotFoundError, "字典项不存在")
		} else {
			return nil, err
		}
	}
	if result.RowsAffected == 0 {
		return nil, nil
	}
	return &item, nil
}
