package log

import (
	"os"
	"path"
	"path/filepath"
	"payment_demo/internal/config"
	"time"

	logrus_stack "github.com/Gurpartap/logrus-stack"
	rotatelogs "github.com/lestrrat/go-file-rotatelogs"
	"github.com/pkg/errors"
	"github.com/rifflock/lfshook"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

//设置日志存放路径
func Init() {
	logPath := filepath.Join(config.GetInstance().GetString("app_path"), "logs")
	if _, err := os.Stat(logPath); os.IsNotExist(err) {
		os.Mkdir(logPath, os.ModePerm)
	}
	filename := "log"
	ConfigLocalFilesystemLogger(logPath, filename, 10*86400*time.Second, 86400*time.Second)
}

//设置日志文件属性信息
func ConfigLocalFilesystemLogger(logPath string, logFileName string, maxAge time.Duration, rotationTime time.Duration) {
	baseLogPath := path.Join(logPath, logFileName)
	writer, err := rotatelogs.New(
		baseLogPath+".%Y-%m-%d.log",
		rotatelogs.WithLinkName(baseLogPath),      // 生成软链，指向最新日志文件
		rotatelogs.WithMaxAge(maxAge),             // 文件最大保存时间
		rotatelogs.WithRotationTime(rotationTime), // 日志切割时间间隔
	)
	if err != nil {
		log.Errorf("config local file system logger error. %+v", errors.WithStack(err))
	}
	lfHook := lfshook.NewHook(lfshook.WriterMap{
		log.DebugLevel: writer, // 为不同级别设置不同的输出目的
		log.InfoLevel:  writer,
		log.WarnLevel:  writer,
		log.ErrorLevel: writer,
		log.FatalLevel: writer,
		log.PanicLevel: writer,
	}, &log.JSONFormatter{
		TimestampFormat: "2006-01-02 15:04:05.000",
	})
	log.AddHook(lfHook)

	// output caller stack under dev environment
	if viper.GetString("service.runmode") != "release" {
		log.AddHook(logrus_stack.StandardHook())
	}
}
