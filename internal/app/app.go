package app

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"qrstreamer/internal/handler"
	"qrstreamer/internal/provider"
	"qrstreamer/internal/routes"
	"qrstreamer/internal/service"
	"qrstreamer/model/constant"
	"qrstreamer/util"
	"syscall"
	"time"

	"google.golang.org/grpc/connectivity"
)

func Run(cfg *util.Config) {
	ctx := context.WithValue(context.Background(), constant.CtxReqIDKey, "MAIN")

	logger := provider.NewLogger()

	redis, err := provider.NewRedisConnection(ctx)
	if err != nil {
		logger.Errorfctx(provider.AppLog, ctx, false, "Failed connect to Redis: %v", err)
		return
	}

	logger.Infofctx(provider.AppLog, ctx, "Application started")

	app := handler.NewApp(logger)
	hub := handler.NewHub(logger)
	svc := service.NewService(logger, hub, app, redis)

	go hub.Run()

	go func() {
		// Setup gRPC client connection
		grpcServerAddr := fmt.Sprintf("localhost:%d", cfg.Server.Port)

		logger.Infofctx(provider.AppLog, ctx, "Starting gRPC client for server at %s", grpcServerAddr)
		conn, err := app.GRPCClient(grpcServerAddr)
		if err != nil {
			logger.Errorfctx(provider.AppLog, ctx, false, "Failed to create gRPC connection: %v", err)
			return
		}
		defer app.CloseGRPCConnection()

		logger.Infofctx(provider.AppLog, ctx, "gRPC client connected successfully to %s", grpcServerAddr)

		// Monitor connection state
		for {
			select {
			case <-ctx.Done():
				logger.Infofctx(provider.AppLog, ctx, "gRPC client shutting down")
				return
			default:
				// Check connection state
				state := conn.GetState()

				switch state {
				case connectivity.Ready:
					logger.Debugfctx(provider.AppLog, ctx, "gRPC connection is ready")
				case connectivity.Connecting:
					logger.Infofctx(provider.AppLog, ctx, "gRPC connection is connecting...")
				case connectivity.TransientFailure:
					logger.Errorfctx(provider.AppLog, ctx, false, "gRPC connection in transient failure state")
					// Wait for state change or timeout
					if !conn.WaitForStateChange(ctx, state) {
						logger.Errorfctx(provider.AppLog, ctx, false, "Context cancelled while waiting for state change")
						return
					}
				case connectivity.Idle:
					logger.Infofctx(provider.AppLog, ctx, "gRPC connection is idle")
				case connectivity.Shutdown:
					logger.Errorfctx(provider.AppLog, ctx, false, "gRPC connection is shutdown")
					return
				}

				// Wait before next state check
				time.Sleep(10 * time.Second)
			}
		}
	}()

	go func() {
		// Start WS HTTP server
		routes.RegisterRoutes(hub, svc)
		logger.Infofctx(provider.AppLog, ctx, "Websocket Server started on :%d", cfg.Websocket.Port)
		if err := http.ListenAndServe(fmt.Sprintf(":%d", cfg.Websocket.Port), nil); err != nil {
			logger.Errorfctx(provider.AppLog, ctx, false, "Failed to start Websocket Server: %v", err)
		}
	}()

	shutdownCh := make(chan os.Signal, 1)
	signal.Notify(shutdownCh, os.Interrupt, syscall.SIGTERM)

	sig := <-shutdownCh
	logger.Infofctx(provider.AppLog, ctx, "Receiving signal: %s", sig)

	func(logger provider.ILogger) {

		logger.Infofctx(provider.AppLog, ctx, "Successfully stop Application.")
	}(logger)

}
