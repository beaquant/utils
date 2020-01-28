package logger

import (
	"github.com/beaquant/utils/file"
	"github.com/lestrrat-go/file-rotatelogs"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
	"os"
	"time"
)

func newRotateIO(writerDir, logFileName string, ageDay uint, rotationDuration time.Duration) (*rotatelogs.RotateLogs, error) {
	_ = file.IsNotExistMkDir(writerDir)
	logName := writerDir + "/" + logFileName
	return rotatelogs.New(
		logName+".%Y%m%d%H.txt",
		// WithLinkName为最新的日志建立软连接，以方便随着找到当前日志文件
		rotatelogs.WithLinkName(logName),

		// WithRotationTime设置日志分割的时间，这里设置为一小时分割一次
		rotatelogs.WithRotationTime(rotationDuration),

		// WithMaxAge和WithRotationCount二者只能设置一个，
		// WithMaxAge设置文件清理前的最长保存时间，
		// WithRotationCount设置文件清理前最多保存的个数。
		rotatelogs.WithMaxAge(time.Hour*24*time.Duration(ageDay)),
		//rotatelogs.WithRotationCount(24 * 3),
	)
}

func NewLoggerWithRotate(dir, logFileName string, ageDay uint, rotationDuration time.Duration, formatter logrus.Formatter, level logrus.Level, exitHandler func()) *logrus.Logger {
	var Logger = &logrus.Logger{
		Out:       os.Stdout,
		Formatter: formatter,
		Hooks:     make(logrus.LevelHooks),
		// Minimum level to log at (5 is most verbose (debug), 0 is panic)
		Level: level,
	}
	// 设置logrus实例的输出到任意io.writer
	Logger.Out = os.Stdout

	// 为当前logrus实例设置消息输出格式为text格式。
	if formatter == nil {
		Logger.Formatter = &logrus.TextFormatter{}
	} else {
		Logger.Formatter = formatter
	}

	// 设置日志级别
	Logger.Level = level

	// 添加 hook
	Logger.Hooks.Add(newLfsHook(dir, logFileName, ageDay, rotationDuration, level))

	// 让logrus在执行os.Exit(1)之前进行相应的处理。fatal handler可以在系统异常时调用一些资源释放api等，让应用正确的关闭。
	logrus.RegisterExitHandler(exitHandler)

	return Logger
}

// 日志本地文件分割的HOOK
func newLfsHook(dir, logFileName string, ageDay uint, rotationDuration time.Duration, level logrus.Level) logrus.Hook {
	writer, err := newRotateIO(dir, logFileName, ageDay, rotationDuration)
	if err != nil {
		logrus.Errorf("config local file system for logger error: %v", err)
	}

	writeMap := make(lfshook.WriterMap)
	levels := LevelThreshold(level)
	for _, v := range levels {
		writeMap[v] = writer
	}
	lfsHook := lfshook.NewHook(writeMap, &logrus.TextFormatter{TimestampFormat: SimpleDateTimeFormat, FullTimestamp: true})

	return lfsHook
}
