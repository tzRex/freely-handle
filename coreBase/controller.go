package corebase

import (
	"net/http"
	"reflect"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/tzRex/freely-handle/before"
	corecode "github.com/tzRex/freely-handle/coreCode"
)

type IController interface {
	Save(*gin.Context, bool) // 新增和修改
	Del(*gin.Context)
	Info(*gin.Context)
	Table(*gin.Context, bool) // 分页查询和全部查询
	GetPrefix() string
}

type Controller struct {
	Prefix string
	Module string
	Source ISource
}

type delParam struct {
	Ids []uint `json:"ids"`
}

var apiGet = "GET"
var apiPost = "POST"
var apiPut = "PUT"

type apiObj struct {
	From string
	Call func(ctx *gin.Context)
}

// 增加和修改
func (c *Controller) Save(ctx *gin.Context, isAdd bool) {
	if objIsNil(c.Source, ctx) {
		return
	}

	entity := c.Source.GetModel()
	params := reflect.New(reflect.TypeOf(entity).Elem()).Interface()

	if err := ctx.ShouldBindJSON(params); err != nil {
		ctx.JSON(http.StatusBadRequest, corecode.ReqBad(err.Error()))
		return
	}

	var result interface{}
	uid, err := c.Source.SourceSave(params, isAdd)

	if err != nil {
		ctx.JSON(http.StatusOK, corecode.ReqFail(err.Error()))
		return
	}

	if isAdd {
		result = uid
	} else {
		result = "修改成功"
	}

	ctx.JSON(http.StatusOK, corecode.ReqOk(result))
}

// 删除
func (c *Controller) Del(ctx *gin.Context) {
	if objIsNil(c.Source, ctx) {
		return
	}

	stu := delParam{}

	if err := ctx.ShouldBindJSON(&stu); err != nil {
		ctx.JSON(http.StatusBadRequest, corecode.ReqBad(err.Error()))
		return
	}

	if len(stu.Ids) == 0 {
		// 等于零直接返回正确
		ctx.JSON(http.StatusOK, corecode.ReqOk(0))
		return
	}
	rows, err := c.Source.SourceDel(stu.Ids)
	if err != nil {
		ctx.JSON(http.StatusOK, corecode.ReqFail(err.Error()))
		return
	}
	ctx.JSON(http.StatusOK, corecode.ReqOk(rows))
}

// 详情
func (c *Controller) Info(ctx *gin.Context) {
	if objIsNil(c.Source, ctx) {
		return
	}

	strId := ctx.Query("id")
	if strId == "" {
		ctx.JSON(http.StatusBadRequest, corecode.ReqBad("must.be.id"))
		return
	}
	id, err := strconv.ParseUint(strId, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, corecode.ReqBad("id.must.to.number"))
		return
	}
	result, err := c.Source.SourceInfo(uint(id))
	if err != nil {
		ctx.JSON(http.StatusOK, corecode.ReqFail(err.Error()))
		return
	}
	ctx.JSON(http.StatusOK, corecode.ReqOk(result))
}

// 列表
func (c *Controller) Table(ctx *gin.Context, isPage bool) {
	if objIsNil(c.Source, ctx) {
		return
	}

	// params := makeParams(ctx)
	var params = &BaseSearch{}

	if err := ctx.ShouldBindJSON(params); err != nil {
		ctx.JSON(http.StatusBadRequest, corecode.ReqBad(err.Error()))
		return
	}

	var result interface{}
	var err error

	if isPage {
		result, err = c.Source.SourcePage(params)
	} else {
		result, err = c.Source.SourceList(params)
	}

	if err != nil {
		ctx.JSON(http.StatusOK, corecode.ReqFail(err.Error()))
		return
	}
	ctx.JSON(http.StatusOK, corecode.ReqOk(result))
}

// 获取分组
func (c *Controller) GetPrefix() string {
	return c.Prefix
}

// func makeParams(ctx *gin.Context) *BaseSearch {
// 	var params = &BaseSearch{}
// 	var equal = ctx.QueryArray("equal")
// 	var like = ctx.QueryArray("like")
// 	var less = ctx.QueryArray("less")
// 	var greater = ctx.QueryArray("greater")
// }

func objIsNil(model ISource, ctx *gin.Context) bool {
	var isNil = model == nil
	if isNil {
		ctx.JSON(http.StatusOK, corecode.ReqBad(corecode.ErrNilClass.Error()))
		return true
	}
	return false
}

func RegisterController(c IController, apis []string, extendRouter func(router *gin.RouterGroup), middle ...func(ctx *gin.Context)) error {
	group := before.Gin.Group(c.GetPrefix())

	if len(middle) > 0 {
		for _, fun := range middle {
			group.Use(fun)
		}
	}

	var apiMap = map[string]*apiObj{
		"/add":    {apiPost, func(ctx *gin.Context) { c.Save(ctx, true) }},
		"/update": {apiPost, func(ctx *gin.Context) { c.Save(ctx, false) }},
		"/delete": {apiPost, c.Del},
		"/page":   {apiPost, func(ctx *gin.Context) { c.Table(ctx, true) }},
		"/list":   {apiPost, func(ctx *gin.Context) { c.Table(ctx, false) }},
		"/info":   {apiGet, c.Info},
	}

	for _, api := range apis {
		var apiType = apiMap[api]
		if apiType != nil {
			switch apiType.From {
			case apiPut:
				group.PUT(api, apiType.Call)
			case apiPost:
				group.POST(api, apiType.Call)
			case apiGet:
				group.GET(api, apiType.Call)
				// default:
			}
		}
	}

	if extendRouter != nil {
		extendRouter(group)
	}

	return nil
}
