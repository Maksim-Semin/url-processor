package api

import (
	"context"
	"fmt"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	api "main/pkg/api/proto"
	lg "main/pkg/logger"
	"main/pkg/storage"
	"main/pkg/urlProcess"
	"net"
)

const grpcPort = 50051

type server struct {
	api.UnimplementedApiServer
}

func (s *server) ChangeURL(_ context.Context, req *api.ChangeURLMessage) (*api.ChangeURLResponse, error) {
	url := req.URL

	lg.Logger.Info("Received request to ChangeURL", zap.String("req", url))

	newURL, _ := urlProcess.GetNewLink(url)

	return &api.ChangeURLResponse{
		URL: newURL}, nil
}

func (s *server) GetSourceURL(_ context.Context, req *api.ChangeURLMessage) (*api.ChangeURLResponse, error) {
	url := req.GetURL()

	lg.Logger.Info("Received request to GetSourceURL", zap.String("key", url))

	sourceURL, _ := storage.LinkManager("", url, "get")

	return &api.ChangeURLResponse{
		URL: sourceURL}, nil
}

func StartServer() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", grpcPort))
	if err != nil {
		lg.Logger.Fatal("failed GRPC serve", zap.Error(err))
	}

	s := grpc.NewServer()

	reflection.Register(s)

	api.RegisterApiServer(s, &server{})

	lg.Logger.Info("starts successfully")

	if err = s.Serve(lis); err != nil {
		lg.Logger.Fatal("error on grpc serve", zap.Error(err))
	}

}
