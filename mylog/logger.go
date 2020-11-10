package mylog

import (
	"fmt"
	"github.com/rs/zerolog"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
	"path/filepath"
	"time"
)

var (
	Frontend zerolog.Logger
	Backend  zerolog.Logger
)

func SgNow() time.Time {
	location, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		panic(fmt.Sprintf("load location failed, err=%v", err))
	}

	return time.Now().In(location)
}

func ConfigLoggers() {
	var perm = os.ModePerm
	dir, err := filepath.Abs(filepath.Dir("./log"))
	if err != nil {
		panic(err)
	}
	err = os.MkdirAll(dir, perm)
	if err != nil {
		panic(err)
	}
	//zerolog.CallerSkipFrameCount = 3
	zerolog.TimestampFunc = SgNow

	Frontend = zerolog.New(&lumberjack.Logger{
		Filename:   "./log/frontend.log",
		MaxSize:    500,
		MaxBackups: 10,
		MaxAge:     10,
		Compress:   true,
	}).With().Caller().Timestamp().Logger()
	fmt.Println("Frontend log init succeed.")

	Backend = zerolog.New(&lumberjack.Logger{
		Filename:   "./log/backend.log",
		MaxSize:    500,
		MaxBackups: 10,
		MaxAge:     10,
		Compress:   true,
	}).With().Caller().Timestamp().Logger()
	fmt.Println("Backend log init succeed.")
}
