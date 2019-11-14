package main

import (
	"flag"
	"go-war/app"
	"go-war/internal/mlog"
)

var (
	p = "8080" //端口
)

func init() {
	flag.StringVar(&p, `p`, "", `服务端口 例如: -p 80 (默认读取环境变量PORT的值)`)
	flag.Parse()
	mlog.ConfigLocalFilesystemLogger("")
}
func main() {
	router := app.GetRouters()
	router.Run(":" + p)
}
