package relation

import (
	"fmt"
	"reflect"
	"strings"
)

const (
	// FromQueryTag tag标记
	FromQueryTag = "search"
	// Mysql 数据库标识
	Mysql = "mysql"
	// Postgres 数据库标识
	Postgres = "postgres"
)

// ResolveSearchQuery 解析
/**
 * 	exact / iexact 等于
 * 	contains / icontains 包含
 *	gt / gte 大于 / 大于等于
 *	lt / lte 小于 / 小于等于
 *	startswith / istartswith 以…起始
 *	endswith / iendswith 以…结束
 *	in
 *	isnull
 *  order 排序		e.g. order[key]=desc     order[key]=asc
 */
func ResolveSearchQuery(driver string, q interface{}, condition Condition) {
	qType := reflect.TypeOf(q)
	qValue := reflect.ValueOf(q)
	var tag string
	var ok bool
	var t *resolveSearchTag
	for i := 0; i < qType.NumField(); i++ {
		tag, ok = "", false
		tag, ok = qType.Field(i).Tag.Lookup(FromQueryTag)
		if !ok {
			//递归调用
			ResolveSearchQuery(driver, qValue.Field(i).Interface(), condition)
			continue
		}
		switch tag {
		case "-":
			continue
		}
		t = makeTag(tag)
		if qValue.Field(i).IsZero() {
			continue
		}
		//解析
		switch t.Type {
		case "or":
			//获取属性类型
			kind := qType.Field(i).Type.Kind()
			if kind == reflect.Array || kind == reflect.Slice {
				fieldValue := qValue.Field(i)
				if fieldValue.Len() > 0 {
					chilrenCondition := make([]*GormCondition, 0)
					for j := 0; j < fieldValue.Len(); j++ {
						childCondition := &GormCondition{
							GormPublic: GormPublic{},
							Join:       make([]*GormJoin, 0),
						}
						chilrenCondition = append(chilrenCondition, childCondition)

						elem := fieldValue.Index(j)
						ResolveSearchQuery(driver, elem.Interface(), childCondition)
					}
					query := ""
					args := make([]interface{}, 0)
					for k := 0; k < len(chilrenCondition); k++ {
						itemQuery := ""
						for key, value := range chilrenCondition[k].Where {
							itemQuery = itemQuery + " AND " + key
							args = append(args, value...)
						}
						itemQuery = "(" + itemQuery[4:] + ")"
						if i > 0 {
							query = query + " OR "
						}
						query = query + itemQuery
					}
					if len(query) > 0 {
						condition.SetWhere(query[3:], args)
					}
				}
			}
		case "left":
			join := condition.SetJoinOn(t.Type, fmt.Sprintf(
				"left join `%s` on `%s`.`%s` = `%s`.`%s`",
				t.Join,
				t.Join,
				t.On[0],
				t.Table,
				t.On[1],
			))
			ResolveSearchQuery(driver, qValue.Field(i).Interface(), join)
		case "right":
			//右连接
			join := condition.SetJoinOn(t.Type, fmt.Sprintf(
				"right left join `%s` on `%s`.`%s` = `%s`.`%s`",
				t.Join,
				t.Join,
				t.On[0],
				t.Table,
				t.On[1],
			))
			ResolveSearchQuery(driver, qValue.Field(i).Interface(), join)
		case "inner":
			//内连接
			join := condition.SetJoinOn(t.Type, fmt.Sprintf(
				"join `%s` on `%s`.`%s` = `%s`.`%s`",
				t.Join,
				t.Join,
				t.On[0],
				t.Table,
				t.On[1],
			))
			ResolveSearchQuery(driver, qValue.Field(i).Interface(), join)
		case "exact", "iexact":
			condition.SetWhere(fmt.Sprintf("`%s`.`%s` = ?", t.Table, t.Column), []interface{}{qValue.Field(i).Interface()})
		case "contains", "icontains":
			//fixme mysql不支持ilike
			if driver == Postgres && t.Type == "icontains" {
				condition.SetWhere(fmt.Sprintf("`%s`.`%s` ilike ?", t.Table, t.Column), []interface{}{"%" + qValue.Field(i).String() + "%"})
			} else {
				columns := strings.Split(t.Column, ",")
				if len(columns) > 1 {
					clauses := make([]string, len(columns))
					values := make([]interface{}, len(columns))

					currentValue := "%" + qValue.Field(i).String() + "%"
					for i, columnName := range columns {
						clauses[i] = fmt.Sprintf("`%s`.`%s` like ?", t.Table, columnName)
						values[i] = currentValue
					}
					condition.SetWhere(strings.Join(clauses, " or "), values)
				} else {
					//获取属性类型,如果是数组或者切片，需要拼接query语句
					kind := qType.Field(i).Type.Kind()
					if kind == reflect.Array || kind == reflect.Slice {
						fieldValue := qValue.Field(i)
						if fieldValue.Len() > 0 {
							query := ""
							args := make([]interface{}, 0)
							for j := 0; j < fieldValue.Len(); j++ {
								elem := fieldValue.Index(j)
								query = query + " OR " + fmt.Sprintf("`%s`.`%s` like ?", t.Table, t.Column)
								args = append(args, "%"+elem.String()+"%")
							}
							if len(query) > 0 {
								condition.SetWhere(query[3:], args)
							}
						}
					} else {
						condition.SetWhere(fmt.Sprintf("`%s`.`%s` like ?", t.Table, t.Column), []interface{}{"%" + qValue.Field(i).String() + "%"})
					}

				}
			}
		case "gt":
			condition.SetWhere(fmt.Sprintf("`%s`.`%s` > ?", t.Table, t.Column), []interface{}{qValue.Field(i).Interface()})
		case "gte":
			condition.SetWhere(fmt.Sprintf("`%s`.`%s` >= ?", t.Table, t.Column), []interface{}{qValue.Field(i).Interface()})
		case "lt":
			condition.SetWhere(fmt.Sprintf("`%s`.`%s` < ?", t.Table, t.Column), []interface{}{qValue.Field(i).Interface()})
		case "lte":
			condition.SetWhere(fmt.Sprintf("`%s`.`%s` <= ?", t.Table, t.Column), []interface{}{qValue.Field(i).Interface()})
		case "not":
			condition.SetWhere(fmt.Sprintf("`%s`.`%s` != ?", t.Table, t.Column), []interface{}{qValue.Field(i).Interface()})
		case "startswith", "istartswith":
			if driver == Postgres && t.Type == "istartswith" {
				condition.SetWhere(fmt.Sprintf("`%s`.`%s` ilike ?", t.Table, t.Column), []interface{}{qValue.Field(i).String() + "%"})
			} else {
				condition.SetWhere(fmt.Sprintf("`%s`.`%s` like ?", t.Table, t.Column), []interface{}{qValue.Field(i).String() + "%"})
			}
		case "endswith", "iendswith":
			if driver == Postgres && t.Type == "iendswith" {
				condition.SetWhere(fmt.Sprintf("`%s`.`%s` ilike ?", t.Table, t.Column), []interface{}{"%" + qValue.Field(i).String()})
			} else {
				condition.SetWhere(fmt.Sprintf("`%s`.`%s` like ?", t.Table, t.Column), []interface{}{"%" + qValue.Field(i).String()})
			}
		case "in":
			field := qValue.Field(i)
			kind := field.Kind()
			// 判断是否为切片或数组类型
			if kind == reflect.Slice || kind == reflect.Array {
				// 获取切片或数组长度
				length := field.Len()
				// 判断切片或数组长度是大于零
				if length > 0 {
					condition.SetWhere(fmt.Sprintf("`%s`.`%s` in (?)", t.Table, t.Column), []interface{}{field.Interface()})
				}
			}
		case "notin":
			field := qValue.Field(i)
			kind := field.Kind()
			// 判断是否为切片或数组类型
			if kind == reflect.Slice || kind == reflect.Array {
				// 获取切片或数组长度
				length := field.Len()
				// 判断切片或数组长度是大于零
				if length > 0 {
					condition.SetWhere(fmt.Sprintf("`%s`.`%s` not in (?)", t.Table, t.Column), []interface{}{field.Interface()})
				}
			}
		case "isnull":
			if !(qValue.Field(i).IsZero() && qValue.Field(i).IsNil()) {
				condition.SetWhere(fmt.Sprintf("`%s`.`%s` is null", t.Table, t.Column), make([]interface{}, 0))
			}
		case "order":
			switch strings.ToLower(qValue.Field(i).String()) {
			case "desc", "asc":
				condition.SetOrder(fmt.Sprintf("`%s`.`%s` %s", t.Table, t.Column, qValue.Field(i).String()))
			}
		}
	}
}
