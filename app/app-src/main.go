package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"time"

	"github.com/siskinc/todolist/app/app-src/module/mylog"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/siskinc/todolist/app/app-src/handler"
	"github.com/siskinc/todolist/app/app-src/module/configure"
)

func Usage(program string) {
	fmt.Printf("\nusage: %s conf/cf.json\n", program)
	fmt.Printf("\nconf/cf.json      configure file\n")
}

func main() {
	if len(os.Args) != 2 {
		Usage(os.Args[0])
		os.Exit(-1)
	}
	//设置官方日志包log输出格式
	log.SetFlags(log.LstdFlags)
	log.Println("[Main] Starting program")
	defer log.Println("[Main] Exit program successful.")
	//创建一个context
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel() //send cancel operator when exit
	//在启动http服务前,加载单实例资源
	Init(ctx, os.Args[1])

	//配置todolist
	conf := configure.GetConfigure()
	todolist := conf.TODOList
	mode, host, url, port, timeOutRead, timeOutWrite := todolist.Mode, todolist.Host, todolist.URL, todolist.Port, todolist.TimeOutReadS, todolist.TimeOutWriteS
	gin.SetMode(mode)
	router := gin.New()
	useMiddleware(router)                     //配置使用中间件
	allRouter(router, fmt.Sprintf("%s", url)) //配置路由

	//起一个goroutine 跑http服务
	s := &http.Server{
		Addr:         fmt.Sprintf("%s:%d", host, port),
		Handler:      router,
		ReadTimeout:  time.Duration(timeOutRead) * time.Second,
		WriteTimeout: time.Duration(timeOutWrite) * time.Second,
	}
	go func(s *http.Server) {
		log.Printf("[Main] http server start\n")
		err := s.ListenAndServe()
		log.Printf("[Main] http server stop (%+v)\n", err)
	}(s)
	// Trap SIGINT to trigger a shutdown.
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt, os.Kill)
	for {
		select {
		case sig := <-signals:
			log.Println("[Main] Catch signal", sig)
			//平滑关闭server
			err := s.Shutdown(context.Background())
			log.Printf("[Main] start gracefully shuts down http serve %+v", err)
			return
		}
	}
}

func allRouter(route *gin.Engine, prefix string) {
	handler.AllRoute(route, prefix)
}

// Init func
func Init(ctx context.Context, filename string) {
	SetupCPU()
	configure.Configure(ctx, filename)
	mylog.MyLog(ctx)
}

// SetupCPU 配置程序使用几个cpu
func SetupCPU() {
	num := runtime.NumCPU()
	runtime.GOMAXPROCS(num)
}

//配置todolist 使用哪些中间件
func useMiddleware(router *gin.Engine) {
	//输出访问日志
	router.Use(mylog.Logger())
	//添加session管理
	store := cookie.NewStore([]byte("secret"))
	router.Use(sessions.Sessions("TODOLISTSESSID", store))
	/*
	  这里添加其他中间件,这个请放在最下面
	*/
}
