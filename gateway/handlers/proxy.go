package handlers

import (
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/gin-gonic/gin"
)

func ReverseProxy(target string) gin.HandlerFunc {
	return func(c *gin.Context) {
		remote, err := url.Parse(target)
		if err != nil {
			c.AbortWithStatus(http.StatusBadGateway)
			return
		}
		proxy := httputil.NewSingleHostReverseProxy(remote)
		proxy.ErrorHandler = func(rw http.ResponseWriter, r *http.Request, err error) {
			rw.WriteHeader(http.StatusBadGateway)
		}
		proxy.ServeHTTP(c.Writer, c.Request)
	}
}
