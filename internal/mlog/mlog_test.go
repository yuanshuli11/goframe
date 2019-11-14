package mlog

import (
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"net/http/httptest"
	"testing"
)

func TestConfigLocalFilesystemLogger(t *testing.T) {

	ConfigLocalFilesystemLogger("/data0/www/gopath/src/Heimdallr/configs")
}

func TestContext(t *testing.T) {
	router := gin.Default()

	uri := "/test1"
	router.GET(uri, func(c *gin.Context) {
		Access(c, 100)
	})
	req := httptest.NewRequest("GET", uri, nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)


	uri = "/test0"
	router.GET(uri, func(c *gin.Context) {
		assert.Equal(t, "1.2.3.11", GetClientIp(c))
	})
	req = httptest.NewRequest("GET", uri, nil)
	req.Header.Add("X-REAl-IP", "1.2.3.11")
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)


	uri = "/test2"
	router.GET(uri, func(c *gin.Context) {
		GetUniqid(c)
		DebugCtx(c, "/tmp")
		InfoCtx(c, "/tmp")
		WarnCtx(c, "/tmp")
		ErrorCtx(c, "/tmp")
		Debug("/tmp")
		Info("/tmp")
		Warn("tmp")
		assert.Equal(t, "1.2.3.0", GetClientIp(c))

	})
	req = httptest.NewRequest("GET", uri, nil)
	req.Header.Add("HTTP-X-FORWARDED-FOR", "1.2.3.0")
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	uri = "/test3"
	router.GET(uri, func(c *gin.Context) {
		assert.Equal(t, "1.2.3.1", GetClientIp(c))
	})
	req = httptest.NewRequest("GET", uri, nil)
	req.Header.Add("X-FORWARDED-FOR", "1.2.3.1")
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	uri = "/test4"
	router.GET(uri, func(c *gin.Context) {
		assert.Equal(t, "1.2.3.2", GetClientIp(c))
	})
	req = httptest.NewRequest("GET", uri, nil)
	req.Header.Add("HTTP-CLIENT-IP", "1.2.3.2")
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	uri = "/test5"
	router.GET(uri, func(c *gin.Context) {
		assert.Equal(t, "1.2.3.3", GetClientIp(c))
	})
	req = httptest.NewRequest("GET", uri, nil)
	req.Header.Add("CLIENT-IP", "1.2.3.3")
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)


	uri = "/test6"
	router.GET(uri, func(c *gin.Context) {
		assert.Equal(t, "127.0.0.1", GetClientIp(c))
	})
	req = httptest.NewRequest("GET", uri, nil)
	req.Header.Add("CLIENT-IP", "::1")
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
}

func TestGenLogId(t *testing.T) {
	id := GenLogId()
	assert.NotEmpty(t, id)
}
