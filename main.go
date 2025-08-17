package main

import (
	"context"
	"errors"
	"github.com/fatih/color"
	"github.com/spf13/viper"
	"github.com/trancecho/mundo-prd-manager/initialize/gin"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	handler := gin.GinInit()
	srv := &http.Server{
		Addr:    viper.GetString("app.URL"),
		Handler: handler,
	}
	// 启动服务
	go func() {
		if err := srv.ListenAndServe(); err != nil && errors.Is(err, http.ErrServerClosed) {
			color.Red("Server Error: %s", err.Error())
		}
	}()
	color.Green("Server Run At: http://localhost:%s", viper.GetString("app.Port"))
	// 等待中断信号
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	color.Blue("Shutdown Server ...")
	// 关闭服务
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		color.Red("Server Shutdown Error: %s", err.Error())
	} else {
		color.Green("Server Shutdown Gracefully")
	}
}
