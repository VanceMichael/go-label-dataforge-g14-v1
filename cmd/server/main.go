package main

import (
	"context"
	"github.com/VanceMichael/go-label-dataforge-g14-v1/internal/auth"
	"github.com/VanceMichael/go-label-dataforge-g14-v1/internal/clock"
	"github.com/VanceMichael/go-label-dataforge-g14-v1/internal/config"
	"github.com/VanceMichael/go-label-dataforge-g14-v1/internal/httpapi"
	"github.com/VanceMichael/go-label-dataforge-g14-v1/internal/service"
	"github.com/VanceMichael/go-label-dataforge-g14-v1/internal/storage/sqlite"
	"github.com/VanceMichael/go-label-dataforge-g14-v1/internal/worker"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	cfg := config.Load()
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()
	db, e := sqlite.Open(ctx, cfg.DatabasePath)
	if e != nil {
		panic(e)
	}
	defer db.Close()
	repo := sqlite.NewRepo(db)
	clk := clock.Real{}
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	wctx, wcancel := context.WithCancel(ctx)
	defer wcancel()
	go worker.Worker{Store: repo, Clock: clk, Interval: cfg.WorkerInterval, Log: logger}.Run(wctx)
	h := httpapi.New(httpapi.Handler{Auth: auth.Manager{Store: repo, Clock: clk, TTL: cfg.SessionTTL}, Registry: service.Registry{Store: repo, Clock: clk}, Authorizer: service.Authorizer{Store: repo, Clock: clk}, Sandbox: service.Sandbox{Store: repo, Clock: clk}, Products: service.ProductService{Store: repo, Clock: clk}})
	srv := &http.Server{Addr: cfg.ListenAddr, Handler: h}
	go func() { <-ctx.Done(); _ = srv.Shutdown(context.Background()) }()
	logger.Info("dataforge listening", "addr", cfg.ListenAddr)
	if e := srv.ListenAndServe(); e != nil && e != http.ErrServerClosed {
		logger.Error("server stopped", "err", e)
	}
}
