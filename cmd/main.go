package main

import (
	"net"

	"github.com/NajmiddinAbdulhakim/post-service/config"
	pb "github.com/NajmiddinAbdulhakim/post-service/genproto"
	"github.com/NajmiddinAbdulhakim/post-service/pkg/db"
	"github.com/NajmiddinAbdulhakim/post-service/pkg/logger"
	"github.com/NajmiddinAbdulhakim/post-service/service"
	grpcClient "github.com/NajmiddinAbdulhakim/post-service/service/grpc_client"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	cfg := config.Load()

	log := logger.New(cfg.LogLevel, "post-service")
	defer logger.Cleanup(log)

	log.Info("main: sqlxConfig",
		logger.String("host", cfg.PostgresHost),
		logger.Int("port", cfg.PostgresPort),
		logger.String("database", cfg.PostgresDatabase))

	connDB, err := db.ConnectToDB(cfg)
	if err != nil {
		log.Fatal("sqlx connection to postgres error", logger.Error(err))
	}

	grpcC, err := grpcClient.New(cfg)
	if err != nil {
		log.Fatal("grpc connection to userservice error", logger.Error(err))
		return
	}

	postService := service.NewPostService(connDB, log, grpcC)

	lis, err := net.Listen("tcp", cfg.RPCPort)
	if err != nil {
		log.Fatal("Error while listening: %v", logger.Error(err))
	}

	s := grpc.NewServer()
	pb.RegisterPostServiceServer(s, postService)
	log.Info("main: server running",
		logger.String("port", cfg.RPCPort))
	reflection.Register(s)
	if err := s.Serve(lis); err != nil {
		log.Fatal("Error while listening: %v", logger.Error(err))
	}
}
