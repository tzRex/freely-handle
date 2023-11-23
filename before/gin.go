package before

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

var Gin *gin.Engine

var limiter = rate.NewLimiter(rate.Limit(120), 10)

func CreateGin(isProd bool) {
	Gin = gin.New()

	if !isProd {
		Gin.Use(gin.Logger())
	}

	Gin.Use(crosHeader)
	Gin.Use(limitRequest)
	Gin.Use(exceptionHandler)
}

/**
 * 配置请求头，允许跨域访问
 */
func crosHeader(ctx *gin.Context) {
	header := ctx.Writer.Header()

	header.Add("Access-Control-Allow-Origin", "*")
	header.Set("Access-Control-Allow-Methods", "*")
	header.Set("Access-Control-Allow-Headers", "*")

	if ctx.Request.Method == "OPTIONS" {
		ctx.JSON(http.StatusOK, gin.H{"code": 1000, "msg": "request.options"})
	}

	ctx.Next()
}

/**
 * 全局捕获panic的中间件
 */
func exceptionHandler(ctx *gin.Context) {
	defer func() {
		if err := recover(); err != nil {
			var msg string
			if r, ok := err.(error); ok {
				msg = r.Error()
			} else {
				msg = fmt.Sprint(r)
			}
			ctx.JSON(http.StatusInternalServerError, gin.H{"code": 2002, "msg": msg})
			ctx.Abort()
		}
	}()
	ctx.Next()
}

/**
 * 限流接口：一秒钟最多120个接口，并行10
 */
func limitRequest(ctx *gin.Context) {
	if !limiter.Allow() {
		ctx.String(http.StatusTooManyRequests, "too.many.request")
		ctx.Abort()
	} else {
		ctx.Next()
	}
}
