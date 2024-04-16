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

func (s *server) ChangeURL(_ context.Context, req *api.URLRequest) (*api.URLResponse, error) {
	url := req.URL

	lg.Logger.Info("Received request to ChangeURL", zap.String("req", url))

	if len(url) == 0 {
		return &api.URLResponse{
			URL:   "",
			Error: "field URL should not be empty",
		}, nil
	}

	newURL, _ := urlProcess.GetNewLink(url)

	return &api.URLResponse{
		URL:   newURL,
		Error: "",
	}, nil
}

func (s *server) GetSourceURL(_ context.Context, req *api.URLRequest) (*api.URLResponse, error) {
	url := req.GetURL()

	lg.Logger.Info("Received request to GetSourceURL", zap.String("key", url))

	sourceURL, _ := storage.LinkManager("", url, "get")

	if len(sourceURL) == 0 {
		return &api.URLResponse{
			URL:   "",
			Error: "there is no data on such a code",
		}, nil
	}

	return &api.URLResponse{
		URL:   sourceURL,
		Error: "",
	}, nil
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
