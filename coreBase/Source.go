package corebase

import (
	"reflect"
	"strings"

	"github.com/tzRex/freely-handle/before"
	corecode "github.com/tzRex/freely-handle/coreCode"
	"github.com/tzRex/freely-handle/utils"
	"gorm.io/gorm"
)

type ISource interface {
	SourceSave(interface{}, bool) (uint, error)
	SourceDel([]uint) (int64, error)
	SourceInfo(id uint) (interface{}, error)
	SourcePage(*BaseSearch) (interface{}, error)
	SourceList(*BaseSearch) (interface{}, error)
	GetModel() IModel
}

var (
	BasePageNum   = 1
	BasePageSize  = 10
	BaseListLimit = 1000
)

type Source struct {
	Model           IModel
	TableOmitFields []string
	// SelectFields []string
}

type PageResult struct {
	Page *BasePage   `json:"page"`
	List interface{} `json:"list"`
}

/**
 * 新增和修改
 */
func (s *Source) SourceSave(data interface{}, isAdd bool) (uint, error) {
	if err := ValidStuct(data); err != nil {
		return 0, err
	}
	var table = before.GetDB()
	var query *gorm.DB

	// data 都是指针类型，需要获取对应地址的值
	refID := reflect.ValueOf(data).Elem().FieldByName("ID")
	ID := refID.Interface().(uint)

	if isAdd {
		query = table.Omit("id").Create(data)
	} else {
		if ID == 0 {
			return 0, corecode.ErrIdFail
		}
		query = table.Omit("id", "username").Save(data)
	}

	if query.Error != nil {
		return 0, query.Error
	}

	return ID, nil
}

/**
 * 删除
 */
func (s *Source) SourceDel(ids []uint) (int64, error) {
	if len(ids) == 0 {
		return 0, nil
	}

	if s.Model == nil {
		return 0, corecode.ErrNilClass
	}

	query := before.GetDB().Where("id IN ?", ids).Delete(s.Model)
	if query.Error != nil {
		return 0, query.Error
	}

	return query.RowsAffected, nil
}

/**
 * 详情（各自实例的详情需自行实现）
 */
func (s *Source) SourceInfo(id uint) (interface{}, error) {
	return nil, nil
}

/**
 * 分页
 */
func (s *Source) SourcePage(search *BaseSearch) (interface{}, error) {
	if search.Page == 0 {
		search.Page = BasePageNum
	}

	if search.Size == 0 {
		search.Size = BasePageSize
	}

	if s.Model == nil {
		return nil, corecode.ErrNilClass
	}

	var model = reflect.TypeOf(s.Model)
	sliceT := reflect.SliceOf(model)
	list := reflect.MakeSlice(sliceT, 0, 0).Interface()

	offset := (search.Page - 1) * search.Size

	m := before.GetDB().Model(s.Model)
	wh := SourceSearchWhere(m, search)
	fields := SourceFieldsFiltter(wh, s.TableOmitFields)
	query := fields.Limit(search.Size).Offset(offset).Find(&list)

	if query.Error != nil {
		return nil, query.Error
	}

	var result = &PageResult{
		Page: &BasePage{
			PageSize: search.Page,
			PageNum:  search.Size,
		},
		List: list,
	}

	query = SourceSearchWhere(m, search).Count(&result.Page.PageTotal)
	if query.Error != nil {
		return nil, query.Error
	}

	return result, nil
}

/**
 * 列表
 */
func (s *Source) SourceList(search *BaseSearch) (interface{}, error) {
	if s.Model == nil {
		return nil, corecode.ErrNilClass
	}

	var model = reflect.TypeOf(s.Model)
	sliceT := reflect.SliceOf(model)
	list := reflect.MakeSlice(sliceT, 0, 0).Interface()

	m := before.GetDB().Model(s.Model)
	wh := SourceSearchWhere(m, search)
	fields := SourceFieldsFiltter(wh, s.TableOmitFields)
	result := fields.Limit(BaseListLimit).Find(&list)

	if result.Error != nil {
		return nil, result.Error
	}
	return list, nil
}

func (s *Source) GetModel() IModel {
	return s.Model
}

func SourceSearchWhere(db *gorm.DB, search *BaseSearch) *gorm.DB {
	var conds []interface{}
	sql := ""
	if search == nil {
		return db
	}
	// 等于，可多个
	if len(search.Equal) > 0 {
		for _, q := range search.Equal {
			field := utils.StrToSnake(q.Field)
			sql += " AND " + field + " IN (?)"
			conds = append(conds, q.Value)
		}
	}
	// 模糊查询
	if len(search.Like) > 0 {
		for _, q := range search.Like {
			val, ok := q.Value.(string)
			if !ok {
				val = ""
			}
			field := utils.StrToSnake(q.Field)
			sql += " AND " + field + " like ?"
			conds = append(conds, "%"+val+"%")
		}
	}
	// 小于
	if len(search.Less) > 0 {
		for _, q := range search.Less {
			val, ok := q.Value.(int)
			if !ok {
				continue
			}
			field := utils.StrToSnake(q.Field)
			sql += " AND " + field + " < ?"
			conds = append(conds, val)
		}
	}
	// 大于
	if len(search.Greater) > 0 {
		for _, q := range search.Greater {
			val, ok := q.Value.(int)
			if !ok {
				continue
			}
			field := utils.StrToSnake(q.Field)
			sql += " AND " + field + " > ?"
			conds = append(conds, val)
		}
	}

	sql = strings.TrimPrefix(sql, " AND")

	db = db.Where(sql, conds...)

	var od string

	if search.Order != "" {
		od = search.Order
	} else {
		od = "updated_at" // 默认以修改时间排序
	}

	if search.Sort != "" {
		od += " " + search.Sort
	} else {
		// DESC and ASC
		od += " DESC" // 默认升序 DESC
	}

	return db.Order(od)
}

func SourceFieldsFiltter(db *gorm.DB, filter []string) *gorm.DB {
	if filter != nil {
		db.Omit(filter...)
	}
	return db
}
