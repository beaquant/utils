package logger

import (
	"fmt"
	"os"
	//"github.com/json-iterator/go"
	"github.com/beaquant/utils/wx"
	"github.com/sirupsen/logrus"
)

var (
	log *logrus.Logger
	//json   = jsoniter.ConfigCompatibleWithStandardLibrary
)

// Supported log levels
var AllLevels = []logrus.Level{
	logrus.DebugLevel,
	logrus.InfoLevel,
	logrus.WarnLevel,
	logrus.ErrorLevel,
	logrus.FatalLevel,
	logrus.PanicLevel,
}

func init() {
}

func NewLogger() *logrus.Logger {
	if log == nil {
		log = &logrus.Logger{
			Out:       os.Stdout,
			Formatter: &logrus.TextFormatter{ForceColors: true, TimestampFormat: "2006-01-02 15:04:05", FullTimestamp: true},
			Hooks:     make(logrus.LevelHooks),
			// Minimum level to log at (5 is most verbose (debug), 0 is panic)
			Level: logrus.DebugLevel,
		}
		fileHook, err := NewLogrusFileHook("log.txt", os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666, logrus.DebugLevel)
		if err == nil {
			log.Hooks.Add(fileHook)
		}
		//
		//wxHook, err := NewWxHook(logrus.WarnLevel, "https://sc.ftqq.com/", "xxxx")
		//if err == nil {
		//	log.Hooks.Add(wxHook)
		//}
	}
	return log
}

type LogrusFileHook struct {
	file           *os.File
	flag           int
	chmod          os.FileMode
	formatter      *logrus.TextFormatter
	acceptedLevels []logrus.Level
}

// Fire event
func (hook *LogrusFileHook) Fire(entry *logrus.Entry) error {

	plainformat, err := hook.formatter.Format(entry)
	line := string(plainformat)
	_, err = hook.file.WriteString(line)
	if err != nil {
		fmt.Fprintf(os.Stderr, "unable to write file on filehook(entry.String)%v", err)
		return err
	}

	return nil
}

func (hook *LogrusFileHook) Levels() []logrus.Level {
	if hook.acceptedLevels == nil {
		return AllLevels
	}
	return hook.acceptedLevels
}

// LevelThreshold - Returns every logging level above and including the given parameter.
func LevelThreshold(l logrus.Level) []logrus.Level {
	for i := range AllLevels {
		if AllLevels[i] == l {
			return AllLevels[i:]
		}
	}
	return []logrus.Level{}
}

func NewLogrusFileHook(filename string, flag int, chmod os.FileMode, l logrus.Level) (*LogrusFileHook, error) {
	var f *os.File
	var err1 error
	plainFormatter := &logrus.TextFormatter{DisableColors: true}
	if checkFileIsExist(filename) {
		f, err1 = os.OpenFile(filename, flag, chmod)
	} else {
		f, err1 = os.Create(filename)
	}
	if err1 != nil {
		fmt.Fprintf(os.Stderr, "unable to write file on filehook %v", err1)
		return nil, err1
	}

	return &LogrusFileHook{f, flag, chmod, plainFormatter, LevelThreshold(l)}, err1
}

func checkFileIsExist(filename string) bool {
	var exist = true
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		exist = false
	}
	return exist
}

type WxHook struct {
	wxpush         *wx.WxPush
	formatter      *logrus.TextFormatter
	acceptedLevels []logrus.Level
}

// Fire event
func (hook *WxHook) Fire(entry *logrus.Entry) error {
	subject := ""
	if len(entry.Data) > 0 {
		for k, v := range entry.Data {
			subject += fmt.Sprintf("%s=%s", k, v)
		}
	} else {
		subject += entry.Level.String()
	}
	hook.wxpush.SendWxString(subject, entry.Message)

	return nil
}

func (hook *WxHook) Levels() []logrus.Level {
	if hook.acceptedLevels == nil {
		return AllLevels
	}
	return hook.acceptedLevels
}

func NewWxHook(l logrus.Level, url, key string) (*WxHook, error) {
	plainFormatter := &logrus.TextFormatter{DisableColors: true}
	w := wx.NewWxPush(url, key)
	return &WxHook{w, plainFormatter, LevelThreshold(l)}, nil
}
