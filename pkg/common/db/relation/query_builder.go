package relation

import "gorm.io/gorm"

func MakeCondition(q interface{}) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		condition := &GormCondition{
			GormPublic: GormPublic{},
			Join:       make([]*GormJoin, 0),
		}
		ResolveSearchQuery(Driver, q, condition)
		for _, join := range condition.Join {
			if join == nil {
				continue
			}

			var query = join.JoinOn
			args := make([]interface{}, 0)
			for k, v := range join.Where {
				query = query + " AND " + k
				args = append(args, v...)
			}
			db = db.Joins(query, args...)

			//db = db.Joins(join.JoinOn)
			//for k, v := range join.Where {
			//	db = db.Where(k, v...)
			//}
			//for k, v := range join.Or {
			//	db = db.Or(k, v...)
			//}
			//for _, o := range join.Order {
			//	db = db.Order(o)
			//}
		}
		for k, v := range condition.Where {
			db = db.Where(k, v...)
		}
		for k, v := range condition.Or {
			db = db.Or(k, v...)
		}
		for _, o := range condition.Order {
			db = db.Order(o)
		}
		return db
	}
}

func Paginate(pageSize, pageIndex int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		offset := (pageIndex - 1) * pageSize
		if offset < 0 {
			offset = 0
		}
		return db.Offset(offset).Limit(pageSize)
	}
}

func OrderBy(sort string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		switch sort {
		case "1":
			//id降序
			return db.Order("id desc")
		case "2":
			//id降序
			return db.Order("id asc")
		case "3":
			//更新时间降序
			return db.Order("updated_at desc")
		case "4":
			//更新时间升序
			return db.Order("updated_at asc")
		case "5":
			//先置顶降序，再更新时间降序
			return db.Order("is_pin desc, updated_at desc")
		default:
			return db.Order("id desc")
		}
	}
}

func NotDeleted() func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("deleted = 0")
	}
}
