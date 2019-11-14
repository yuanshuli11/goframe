package sys

import (
	"encoding/json"
	"fmt"

	"github.com/gin-gonic/gin"
)

type JSONResponse struct {
	Errno int         `json:"errno"`
	Msg   string      `json:"msg"`
	Data  interface{} `json:"data"`
}

func ReturnSuccess(c *gin.Context, data interface{}) {
	resp := JSONResponse{0, "success", data}
	j, _ := json.Marshal(resp)
	fmt.Fprintf(c.Writer, "%s", j)
}
func ReturnError(c *gin.Context, errorType string) {
	resp := JSONResponse{ErrorTypes[errorType].ErrNum, ErrorTypes[errorType].ErrTpl, nil}
	j, _ := json.Marshal(resp)
	fmt.Fprintf(c.Writer, "%s", j)
}

func ReturnErrorMsg(c *gin.Context, errorType string, errorMsg string) {
	resp := JSONResponse{ErrorTypes[errorType].ErrNum, errorMsg, nil}
	j, _ := json.Marshal(resp)
	fmt.Fprintf(c.Writer, "%s", j)
}
