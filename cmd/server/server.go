package server

import (
	"context"
	"github.com/spf13/cobra"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"
	xiangqin_backend "xiangqin-backend"
	"xiangqin-backend/pkg"
)

var (
	StartCmd = &cobra.Command{
		Use:          "server",
		Short:        "start server",
		Example:      "start server",
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return run()
		},
	}
)

func run() error {
	cfg, err := xiangqin_backend.GetConfig("config.yaml")
	if err != nil {
		return err
	}
	db, err := xiangqin_backend.NewPostgres(cfg)
	if err != nil {
		return err
	}
	router := pkg.NewRouter(db)
	server := &http.Server{
		Addr:    cfg.Listen.Host + ":" + strconv.Itoa(cfg.Listen.Port),
		Handler: router,
	}
	go func() {
		log.Printf("server start,listen: %s", server.Addr)
		err = server.ListenAndServe()
		if err != nil {
			log.Println("server shutdown")
		}
	}()
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	// 等待关闭信号
	<-signalChan

	// 创建一个上下文对象，设置超时时间
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 关闭服务器
	err = server.Shutdown(ctx)
	if err != nil {
		log.Fatalln("Shutdown error:", err)
	}
	return nil
}
