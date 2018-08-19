package mylog

import (
	"context"
	"log"
	"sync"

	"github.com/cihub/seelog"
	"github.com/siskinc/todolist/app/app-src/module/common"
	"github.com/siskinc/todolist/app/app-src/module/configure"
)

type mylog struct {
	File       string // ordinary log file
	FileAccess string // todolist access log file

	Log       seelog.LoggerInterface
	LogAccess seelog.LoggerInterface

	StructName string
}

var (
	myLog     *mylog
	myLogOnce sync.Once
)

// MyLog func
func MyLog(ctx context.Context) *mylog {
	myLogOnce.Do(func() {
		myLog := &mylog{}
		if err := myLog.loadConfigure(); err != nil {
			log.Fatalln(err)
		}
		myLog.waitStop(ctx)
	})
	return myLog
}

// loadConfigure 从file中读取seelog配置
func (l *mylog) loadConfigure() error {
	conf := configure.GetConfigure()
	file := conf.Log.File
	fileAccess := conf.Log.Access

	logger, err := seelog.LoggerFromConfigAsFile(file)
	if err != nil {
		return err
	}
	l.Log = logger
	loggerAccess, err := seelog.LoggerFromConfigAsFile(fileAccess)
	if err != nil {
		return err
	}
	l.LogAccess = loggerAccess

	l.StructName = common.GetStructName(l)
	log.Printf("[%s] Start\n", l.StructName)
	return nil
}

func (l *mylog) waitStop(ctx context.Context) {
	go func() {
		for {
			select {
			case <-ctx.Done():
				l.LogAccess.Flush()
				l.Log.Flush()
				log.Printf("[%s] Stoped\n", l.StructName)
				return
			}
		}
	}()
}

// Infof 输出info信息
func (l *mylog) Access() seelog.LoggerInterface {
	return l.LogAccess
}

func (l *mylog) Regular() seelog.LoggerInterface {
	return l.Log
}

//得到日志实例
func GetMylog() *mylog {
	return myLog
}
