package merror

import (
	"fmt"
	"go-war/internal/mlog"
	"go-war/internal/sys"

	"github.com/gin-gonic/gin"
)

type Error struct {
	ctx     *gin.Context
	section string
	Context []interface{} // 用于填充%s占位符的值
}

func (e *Error) Error() string {
	_, errstr := e.GetErrNumAndErrStr()
	return errstr
}

func (e *Error) GetErrNumAndErrStr() (errnum int, errstr string) {
	t, ok := sys.ErrorTypes[e.section]
	if !ok {
		if e.ctx != nil {
			mlog.ErrorCtx(e.ctx, "系统错误，未匹配到错误类型")
		} else {
			mlog.Error("系统错误，未匹配到错误类型")
		}
		return 500000, "系统错误"
	}
	errnum = t.ErrNum
	errstr = fmt.Sprintf(t.ErrTpl, e.Context...)
	return
}
func (e *Error) ErrName() string {
	return e.section
}

func (e *Error) ErrMsg() string {

	t, ok := sys.ErrorTypes[e.section]
	if !ok {
		return ""
	}
	return fmt.Sprintf(t.ErrTpl, e.Context...)
}
func New(ctx *gin.Context, section string, context ...interface{}) *Error {
	return &Error{ctx, section, context}
}
