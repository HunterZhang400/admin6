package httpcros

import (
	"github.com/astaxie/beego/plugins/cors"
	"github.com/gin-gonic/gin"
	"net/http"
)

func CROS(ctx *gin.Context) {
	opt := &cors.Options{
		AllowAllOrigins:  true,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowHeaders:     []string{"X-Token", "X-Session-ID", "X-Admin-User-ID", "Origin", "Authorization", "Access-Control-Allow-Origin", "Access-Control-Allow-Headers", "Content-Type"},
		ExposeHeaders:    []string{"X-Token", "X-Session-ID", "Content-Length", "Access-Control-Allow-Origin", "Access-Control-Allow-Headers", "Content-Type"},
		AllowCredentials: true,
	}
	var (
		headerOrigin         = "Origin"
		headerRequestMethod  = "Access-Control-Request-Method"
		headerRequestHeaders = "Access-Control-Request-Headers"
		origin               = ctx.Request.Header.Get(headerOrigin)
		requestedMethod      = ctx.Request.Header.Get(headerRequestMethod)
		requestedHeaders     = ctx.Request.Header.Get(headerRequestHeaders)
		// additional headers to be added
		// to the response.
		headers map[string]string
	)

	if ctx.Request.Method == "OPTIONS" {
		headers = opt.PreflightHeader(origin, requestedMethod, requestedHeaders)
		if opt.AllowAllOrigins {
			headers["Access-Control-Allow-Origin"] = origin
		}
		for key, value := range headers {
			ctx.Writer.Header().Add(key, value)
			//ctx.Output.Header(key, value)
		}
		ctx.Writer.WriteHeader(http.StatusOK)
		return
	}
	headers = opt.Header(origin)

	if opt.AllowAllOrigins {
		headers["Access-Control-Allow-Origin"] = origin
	}

	for key, value := range headers {
		ctx.Writer.Header().Add(key, value)
	}
}
