package server

import (
	"context"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"time"
)

func loggingUnaryInterceptor(logger *zap.Logger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
		start := time.Now()
		resp, err = handler(ctx, req)
		if err != nil {
			logger.Error("request", zap.String("method", info.FullMethod), zap.String("duration", time.Since(start).String()), zap.Error(err))
			return resp, err
		}
		logger.Info("request", zap.String("method", info.FullMethod), zap.String("duration", time.Since(start).String()))
		return resp, err
	}
}

func loggingStreamInterceptor(logger *zap.Logger) grpc.StreamServerInterceptor {
	return func(srv any, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		start := time.Now()
		logger.Info("stream request started", zap.String("method", info.FullMethod))
		err := handler(srv, ss)
		if err != nil {
			logger.Error("stream request failed", zap.String("method", info.FullMethod), zap.String("duration", time.Since(start).String()), zap.Error(err))
			return err
		}
		logger.Info("stream request finished", zap.String("method", info.FullMethod), zap.String("duration", time.Since(start).String()))
		return err
	}
}
