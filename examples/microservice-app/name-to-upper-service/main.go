package main

import (
	context "context"
	"fmt"
	"github.com/caarlos0/env"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	pb "name-service/protocol"
	"net"
	"strings"
)

type NameServiceConfig struct {
	Host string `env:"NAME_LISTEN_HOST"`
	Port int    `env:"NAME_PORT,required"`
}

type server struct {
	pb.UnimplementedNameServiceServer
}

func (s server) NameToUpperCase(ctx context.Context, req *pb.NameReq) (*pb.NameResp, error) {
	name := req.GetName()

	name = strings.ToUpper(name)

	return &pb.NameResp{Name: name}, nil
}

func main() {
	var nameCfg NameServiceConfig

	err := env.Parse(&nameCfg)
	if err != nil {
		log.Fatalf("Couldn't parse env: %v", err)
	}

	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", nameCfg.Host, nameCfg.Port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	log.Infof("Server listens at port: %d", nameCfg.Port)

	grpcServer := grpc.NewServer()
	pb.RegisterNameServiceServer(grpcServer, &server{})
	if err = grpcServer.Serve(lis); err != nil {
		log.Fatalf("unable to start listener: %v", err)
	}
}
