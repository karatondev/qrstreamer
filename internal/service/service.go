package service

import (
	"context"
	"encoding/json"
	"io"
	"qrstreamer/internal/handler"
	"qrstreamer/internal/provider"
	"qrstreamer/model"
	"time"

	proto "qrstreamer/model/pb"

	"github.com/redis/go-redis/v9"
)

type QRStreamer interface {
	StreamWhatsappQR(ctx context.Context, deviceID string) error
}

type service struct {
	logger provider.ILogger
	hub    *handler.Hub
	app    *handler.App
	redis  *redis.Client
}

func NewService(logger provider.ILogger, hub *handler.Hub, app *handler.App, redis *redis.Client) QRStreamer {
	return &service{
		logger: logger,
		hub:    hub,
		app:    app,
		redis:  redis,
	}
}

func (s *service) StreamWhatsappQR(ctx context.Context, deviceID string) error {
	req := &proto.ConnectDeviceRequest{
		Name: deviceID,
	}
	stream, err := s.app.StreamConnectDevice(ctx, req)
	if err != nil {
		s.logger.Errorfctx(provider.AppLog, ctx, false, "Error calling GenerateNumbers: %v", err)
	}

	for {
		resp, err := stream.Recv()
		if err == io.EOF {
			s.logger.Infofctx(provider.AppLog, ctx, "Stream closed by server")
			break
		}
		if err != nil {
			s.logger.Errorfctx(provider.AppLog, ctx, false, "Error receiving stream: %v", err)
		}

		if err := s.emitQRCodeToClient(deviceID, resp.Qr); err != nil {
			s.logger.Errorfctx(provider.AppLog, ctx, false, "Error emitting QR code: %v", err)
		}
	}

	return nil
}

func (s *service) emitQRCodeToClient(deviceID string, qrData string) error {
	message := model.WSMessage{
		Type:      "qr_code",
		DeviceId:  deviceID,
		Data:      qrData,
		Timestamp: time.Now(),
	}

	msgBytes, err := json.Marshal(message)
	if err != nil {
		return err
	}

	s.logger.Infofctx(provider.AppLog, context.Background(), "Emitting QR code  %s", msgBytes)

	// Emit to Websocket client
	s.hub.EmitToClient(deviceID, msgBytes)

	return nil
}
