package server

import (
	"context"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"user/proto"
)

type UserServer struct {
	proto.UnimplementedUserServer
}

var Inter grpc.UnaryServerInterceptor

func Medata() {

	Inter = func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
		// 获取 metadata 里面的值

		md, ok := metadata.FromIncomingContext(ctx)

		if ok {
			zap.S().Info("metadata携带的数据", md)
		}
		resp, err = handler(ctx, req)

		if err != nil {
			return nil, err
		}
		return resp, nil

	}
}
