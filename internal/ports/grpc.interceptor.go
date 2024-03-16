package ports

import (
	"context"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

func ServerInterceptor(logger *logrus.Logger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {

		entry := logrus.NewEntry(logger).WithField("method", info.FullMethod)
		entry.Info("Request Received")

		h, err := handler(ctx, req)
		if err == nil {
			entry.Info("Request Complete")
		} else {
			entry.WithError(err).Error("Request Failed")
		}

		return h, err

	}
}
