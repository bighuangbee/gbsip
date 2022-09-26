package log

import (
	"fmt"
	rotatelogs "github.com/lestrrat/go-file-rotatelogs"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

var(
	fileRootPath string
)

type Log struct {
	*logrus.Logger
}

func New(savePath string)*logrus.Logger{
	fileRootPath = savePath

	if _, err := os.Stat(fileRootPath); os.IsNotExist(err) {
		err := os.MkdirAll(fileRootPath, os.ModePerm)
		if err != nil {
			panic("loger Setup" + err.Error())
		}
	}

	src, err := os.OpenFile(os.DevNull, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	if err != nil {
		panic("loger Setup" + err.Error())
	}

	log := logrus.New()
	log.SetReportCaller(true)

	//输出到文件
	log.Out = src

	logerFormatter := new(LogerFormatter)
	log.SetFormatter(logerFormatter)

	//log.AddHook(newLocalFileLogHook(logrus.ErrorLevel, logerFormatter))
	log.AddHook(newLocalFileLogHook(logrus.InfoLevel, logerFormatter))

	log.SetOutput(os.Stdout)
	return log
}

/**
	写入本地日志文件， 按日期、日志级别分割为不同的文件
*/
func newLocalFileLogHook(level logrus.Level, formatter logrus.Formatter) logrus.Hook {

	fileName := filepath.Join(fileRootPath, level.String() + "_%Y%m%d.txt")

	//文件分割
	writer, err := rotatelogs.New(
		fileName,
		// 最大保存时间(30天)
		rotatelogs.WithMaxAge(30*24*time.Hour),
		// 日志分割间隔(1天)
		rotatelogs.WithRotationTime(24*time.Hour),
	)

	if err != nil {
		fmt.Errorf("config local file system for Loger error: %v", err)
	}

	return lfshook.NewHook(lfshook.WriterMap{
		level: writer,
	}, formatter)

}

func (this *Log)Trace(format string, args ...interface{}){
	this.Trace(format, setPrefix("Trace"))
}

func (this *Log) Infof(format string, args ...interface{}){
	this.Infof(format, setPrefix("Infof"))
}

func (this *Log) Info(args ...interface{}){
	this.Info(setPrefix("Info"), args)
}

func (this *Log) Error(args ...interface{}){
	this.Error(setPrefix("Error"), args)
}


// setPrefix set the prefix of the log output
func setPrefix(level string) string{

	pc, file, line, ok := runtime.Caller(2)
	if ok {
		loc, _ := time.LoadLocation("Asia/Shanghai")
		funcName := runtime.FuncForPC(pc).Name()
		funcName = strings.TrimPrefix(filepath.Ext(funcName), ".")
		timestamp := time.Now().In(loc).Format("2006-01-02 15:04:05.000")

		return fmt.Sprintf("[%s][%s][%s:%d:%s]", strings.ToUpper(level), timestamp, filepath.Base(file), line, funcName)
	}
	return ""
}

//日志输出格式
type LogerFormatter struct{}
func (s *LogerFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	msg := fmt.Sprintf("%s \n", entry.Message)
	return []byte(msg), nil
}

////
type Writer struct{
	Log
}

func (w Writer) Printf(format string,args ...interface{}) {
	fmt.Printf(format, args...)
}
