package user

import (
	"github.com/gin-gonic/gin"
	"goframe/internal/sys"
)

func GetInfo() gin.HandlerFunc {
	return func(c *gin.Context) {

		data := make(map[string]interface{})

		data["user_info"] = "userinfo"
		sys.ReturnSuccess(c, data)
	}
}
