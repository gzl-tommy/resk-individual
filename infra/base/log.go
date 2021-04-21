package base

import (
	"fmt"
	rotatelogs "github.com/lestrrat/go-file-rotatelogs"
	"github.com/mattn/go-colorable"
	"github.com/rifflock/lfshook"
	log "github.com/sirupsen/logrus"
	"github.com/tietang/go-utils"
	"github.com/tietang/props/v3/kvs"
	prefixed "github.com/x-cray/logrus-prefixed-formatter"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strings"
	"time"
	"io"
)

var formatter *prefixed.TextFormatter

func init() {
	// 定义日志格式
	// 不使用 logrus 默认的日志格式配置，采用一个第三方的日志格式配置（prefixed），完全兼容 logrus，并扩展功能
	formatter = &prefixed.TextFormatter{}

	// 开启完整时间戳输出和时间戳格式
	formatter.FullTimestamp = true
	formatter.TimestampFormat = "2006-01-02 15:04:05.000000"

	// 控制台高亮显示
	formatter.ForceFormatting = true
	formatter.DisableColors = false
	formatter.ForceColors = true
	// 设置高亮显示的色彩样式
	formatter.SetColorScheme(&prefixed.ColorScheme{
		InfoLevelStyle:  "green",
		WarnLevelStyle:  "yellow",
		ErrorLevelStyle: "red",
		FatalLevelStyle: "41",
		PanicLevelStyle: "41",
		DebugLevelStyle: "blue",
		PrefixStyle:     "cyan",
		TimestampStyle:  "37",
	})

	//设置日志formatter
	log.SetFormatter(formatter)
	log.SetOutput(colorable.NewColorableStdout())

	// 日志级别，通过环境表里来设置
	// 后期可以变更到配置中来设置
	level := os.Getenv("log.debug")
	level = "true"
	if level == "true" {
		log.SetLevel(log.DebugLevel)
	}

	// 开启调用函数、文件、代码行信息的输出
	log.SetReportCaller(true)

	// 设置函数、文件、代码行信息的输出 hook 
	SetLineNumLogrusHook()

	// 日志文件和滚动配置
	// logrus 默认不提供日志文件的输出功能，需要使用第三方库通过 hook 来实现
	//github.com/lestrrat/go-file-rotatelogs
	//logFileSettings()
}

var lfh *utils.LineNumLogrusHook

func SetLineNumLogrusHook() {
	lfh = utils.NewLineNumLogrusHook()
	lfh.EnableFileNameLog = true
	lfh.EnableFuncNameLog = true
	log.AddHook(lfh)
}

// 将滚动日志writer共享给iris glog output
var log_writer io.Writer

// 初始化 log 配置，配置 logrus 日志文件滚动生成和
func InitLog(conf kvs.ConfigSource) {
	//设置日志输出级别
	level, err := log.ParseLevel(conf.GetDefault("log.level", "info"))
	if err != nil {
		level = log.InfoLevel
	}
	log.SetLevel(level)
	if conf.GetBoolDefault("log.enableLineLog", true) {
		lfh.EnableFileNameLog = true
		lfh.EnableFuncNameLog = true
	} else {
		lfh.EnableFileNameLog = false
		lfh.EnableFuncNameLog = false
	}

	//配置日志输出目录
	logDir := conf.GetDefault("log.dir", "./logs")
	logTestDir, err := conf.Get("log.test.dir")
	if err == nil {
		logDir = logTestDir
	}
	logPath := logDir //+ "/logs"
	logFilePath, _ := filepath.Abs(logPath)
	log.Infof("log dir: %s", logFilePath)
	logFileName := conf.GetDefault("log.file.name", "red-envelop")
	maxAge := conf.GetDurationDefault("log.max.age", time.Hour*24)
	rotationTime := conf.GetDurationDefault("log.rotation.time", time.Hour*1)
	os.MkdirAll(logPath, os.ModePerm)

	baseLogPath := path.Join(logPath, logFileName)
	//设置滚动日志输出writer
	writer, err := rotatelogs.New(
		strings.TrimSuffix(baseLogPath, ".log")+".%Y%m%d%H.log",
		rotatelogs.WithLinkName(baseLogPath),      // 生成软链，指向最新日志文件
		rotatelogs.WithMaxAge(maxAge),             // 文件最大保存时间
		rotatelogs.WithRotationTime(rotationTime), // 日志切割时间间隔
	)
	if err != nil {
		log.Errorf("config local file system logger error. %+v", err)
	}

	//设置日志文件输出的日志格式
	formatter := &log.TextFormatter{}
	formatter.CallerPrettyfier = func(frame *runtime.Frame) (function string, file string) {
		function = frame.Function
		dir, filename := path.Split(frame.File)
		f := path.Base(dir)
		return function, fmt.Sprintf("%s/%s:%d", f, filename, frame.Line)
	}

	lfHook := lfshook.NewHook(lfshook.WriterMap{
		log.DebugLevel: writer, // 为不同级别设置不同的输出目的
		log.InfoLevel:  writer,
		log.WarnLevel:  writer,
		log.ErrorLevel: writer,
		log.FatalLevel: writer,
		log.PanicLevel: writer,
	}, formatter)

	log.AddHook(lfHook)
	log_writer = writer
}
