package mlog

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"go-war/internal/sys"
	"io/ioutil"
	"math/rand"
	"net"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var Logger *zap.Logger

// config logrus log to local filesystem, with file rotation
func ConfigLocalFilesystemLogger(confDir string) {
	//env := "debug"
	// 日志路径
	appLogPath := `/data0/www/applogs/go-war`

	rawJSON, err := ioutil.ReadFile("config/log/debug.json")
	if err != nil {
		fmt.Printf("%s\n", err)
		panic(err)
	}

	var cfg zap.Config
	if err = json.Unmarshal(rawJSON, &cfg); err != nil {
		panic(err)
	}

	opts := []zap.Option{}

	stackLevel := zap.InfoLevel
	if cfg.Development {
		opts = append(opts, zap.Development())
		stackLevel = zap.DebugLevel
	}
	if !cfg.DisableStacktrace {
		opts = append(opts, zap.AddStacktrace(stackLevel))
	}

	if !cfg.DisableCaller {
		opts = append(opts, zap.AddCaller())
		opts = append(opts, zap.AddCallerSkip(2))
	}
	if cfg.Sampling != nil {
		opts = append(opts, zap.WrapCore(func(core zapcore.Core) zapcore.Core {
			return zapcore.NewSampler(core, time.Second, int(cfg.Sampling.Initial), int(cfg.Sampling.Thereafter))
		}))
	}

	if len(cfg.InitialFields) > 0 {
		fs := make([]zap.Field, 0, len(cfg.InitialFields))
		keys := make([]string, 0, len(cfg.InitialFields))
		for k := range cfg.InitialFields {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		for _, k := range keys {
			fs = append(fs, zap.Any(k, cfg.InitialFields[k]))
		}
		opts = append(opts, zap.Fields(fs...))
	}

	// 自动切分日志的writer
	rollingLogger := zapcore.AddSync(&lumberjack.Logger{
		Filename:   appLogPath + `/applogs.log`,
		MaxSize:    4096, // megabytes
		MaxBackups: 50,
		MaxAge:     7, // days
		LocalTime:  true,
	})

	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoder := zapcore.NewJSONEncoder(encoderConfig)

	logLevel := cfg.Level
	cores := []zapcore.Core{} // 可以多目标打印日志 后续直接添加即可
	// 打印在文件中
	cores = append(cores, zapcore.NewCore(encoder, rollingLogger, logLevel))
	// 打印在控制台
	if logLevel.String() == "debug" {
		consoleDebugging := zapcore.Lock(os.Stdout)
		cores = append(cores, zapcore.NewCore(encoder, consoleDebugging, logLevel))
	}

	core := zapcore.NewTee(cores...)
	Logger = zap.New(core, opts...)
	defer Logger.Sync()
}

func Access(c *gin.Context, cost int) {
	//costMap, ok := c.Get("costMap")
	//costStr := ""
	//if ok {
	//	costByte, err := json.Marshal(costMap)
	//	if err != nil {
	//		costStr = ""
	//	} else {
	//		costStr = string(costByte)
	//	}
	//}
	////这里转换成毫秒
	//cost = cost / 1000
	Logger.Info("",
		// Structured context as strongly typed Field values.
		zap.String("type", "access"),
		zap.String("host", c.Request.URL.Host),
		zap.String("uri", c.Request.RequestURI),
		zap.String("method", c.Request.Method),
		zap.Int("http_code", c.Writer.Status()),
		zap.Int("cost", cost),
		zap.String("request_id", GetUniqid(c)),
		zap.String("user_ip", GetClientIp(c)),
		zap.Int64("timestamp", time.Now().Unix()),
		//zap.Float64("gw_cost", gwCost),
		//zap.String("next_cost", costStr),
		//zap.String("post_form", c.Request.PostForm),
	)
}
func GetClientIp(c *gin.Context) (ip string) {

	ipCache, exists := c.Get("HTTP_USER_IP")
	if exists {
		return ipCache.(string)
	}
	remoteAddr := c.Request.RemoteAddr
	if ip := c.Request.Header.Get("X-REAl-IP"); ip != "" {
		remoteAddr = ip
	} else if ip := c.Request.Header.Get("HTTP-X-FORWARDED-FOR"); ip != "" {
		remoteAddr = ip
	} else if ip = c.Request.Header.Get("X-FORWARDED-FOR"); ip != "" {
		remoteAddr = ip
	} else if ip = c.Request.Header.Get("HTTP-CLIENT-IP"); ip != "" {
		remoteAddr = ip
	} else if ip = c.Request.Header.Get("CLIENT-IP"); ip != "" {
		remoteAddr = ip
	} else {
		remoteAddr, _, _ = net.SplitHostPort(remoteAddr)
	}
	index := strings.Index(remoteAddr, ",") // 如果是 10.33.106.148,10.33.106.147
	if index != -1 {
		remoteAddr = remoteAddr[0:index]
	}
	if remoteAddr == "::1" {
		remoteAddr = "127.0.0.1"
	}
	clientIp := net.ParseIP(remoteAddr)
	if clientIp != nil {
		ip = clientIp.String()
	} else {
		ip = ""
	}

	c.Set("HTTP_USER_IP", ip)
	return ip
}
func GetUniqid(c *gin.Context) (uniqid string) {
	id, exists := c.Get("uniqid")
	if !exists {
		uniqid = ""
	} else {
		uniqid = id.(string)
	}
	return uniqid
}
func DebugCtx(c *gin.Context, str string) {
	Logger.Debug(str,
		// Structured context as strongly typed Field values.
		zap.String("type", "application"),
		zap.String("request_id", GetUniqid(c)),
		zap.String("user_ip", GetClientIp(c)),
	)
}
func DebugCtxf(c *gin.Context, format string, v ...interface{}) {
	str := fmt.Sprintf(format, v...)
	Logger.Debug(str,
		// Structured context as strongly typed Field values.
		zap.String("type", "application"),
		zap.String("request_id", GetUniqid(c)),
		zap.String("user_ip", GetClientIp(c)),
	)
}
func InfoCtx(c *gin.Context, str string) {
	Logger.Info(str,
		zap.String("type", "application"),
		zap.String("request_id", GetUniqid(c)),
		zap.String("user_ip", GetClientIp(c)),
	)
}
func InfoCtxf(c *gin.Context, format string, v ...interface{}) {
	str := fmt.Sprintf(format, v...)
	Logger.Info(str,
		zap.String("type", "application"),
		zap.String("request_id", GetUniqid(c)),
		zap.String("user_ip", GetClientIp(c)),
	)
}
func WarnCtx(c *gin.Context, str string) {
	Logger.Warn(str,
		zap.String("type", "application"),
		zap.String("request_id", GetUniqid(c)),
		zap.String("user_ip", GetClientIp(c)),
	)
}
func ErrorCtxf(c *gin.Context, format string, v ...interface{}) {
	configFormat, ok := sys.ErrorTypes[format]
	if ok {
		format = configFormat.LogInfo
	}
	str := fmt.Sprintf(format, v...)
	Logger.Error(str,
		zap.String("type", "application"),
		zap.String("request_id", GetUniqid(c)),
		zap.String("user_ip", GetClientIp(c)),
	)
}
func ErrorCtx(c *gin.Context, str string) {
	Logger.Error(str,
		zap.String("type", "application"),
		zap.String("request_id", GetUniqid(c)),
		zap.String("user_ip", GetClientIp(c)),
	)
}

func Debug(str string) {
	Logger.Debug(str,
		zap.String("type", "common"),
	)
}

func Info(str string) {
	Logger.Info(str,
		zap.String("type", "common"),
	)
}
func Warn(str string) {
	Logger.Warn(str,
		zap.String("type", "common"),
	)
}
func Error(str string) {
	Logger.Error(str,
		zap.String("type", "common"),
	)
}

func GenLogId() string {
	h2 := md5.New()
	rand.Seed(time.Now().Unix())
	str := fmt.Sprintf("%d%d%d", os.Getpid(), time.Now().UnixNano(), rand.Int())

	h2.Write([]byte(str))
	uniqid := hex.EncodeToString(h2.Sum(nil))

	return uniqid
}
