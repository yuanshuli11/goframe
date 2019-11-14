package sys

type ErrorType struct {
	ErrNum  int
	LogInfo string
	ErrTpl  string
}

var ErrorTypes = map[string]ErrorType{
	/// 400xx 登陆注册相关
	"ERR_LOGIN_NO_TOKEN":      {400001, "token不存在", "token不存在"},

}
