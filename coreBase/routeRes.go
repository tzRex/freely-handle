package corebase

type BasePage struct {
	PageSize  int   `json:"pageSize"`
	PageNum   int   `json:"pageNum"`
	PageTotal int64 `json:"pageTotal"`
}

type BaseSearch struct {
	Equal   []*QueryFild `json:"equal"`   // 精确查询
	Like    []*QueryFild `json:"like"`    // 模糊查询
	Less    []*QueryFild `json:"less"`    // 小于
	Greater []*QueryFild `json:"greater"` // 大于
	Page    int          `json:"page"`    // 当前页
	Size    int          `json:"size"`    // 一页多少条数据
	Order   string       `json:"order"`   // 排序字段
	Sort    string       `json:"sort"`    // 排序方式
}

type QueryFild struct {
	Field string      `json:"field"`
	Value interface{} `json:"value"` // 可以是数字、数字切片、字符串
}
