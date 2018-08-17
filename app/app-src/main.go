package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/siskinc/todolist/app/app-src/handler"
)

func main() {
	route := gin.New()
	allRouter(route, fmt.Sprintf("%s", "v1"))

	//起一个goroutine 跑http服务
	s := &http.Server{
		Addr:         fmt.Sprintf("%s:%d", "0.0.0.0", 80),
		Handler:      route,
		ReadTimeout:  time.Duration(30) * time.Second,
		WriteTimeout: time.Duration(30) * time.Second,
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
