package main

import (
	"context"
	"flag"
	"log"
	"net"
	"strconv"

	petname "github.com/go-petname/golang-petname"
	"github.com/ilyakaznacheev/cleanenv"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"

	petnamepb "yadro.com/course/proto"
)

type Config struct {
	Port int `yaml:"port" env:"PETNAME_GRPC_PORT" env-default:"8080"`
}

type server struct {
	petnamepb.UnimplementedPetnameGeneratorServer
}

func (s *server) Ping(_ context.Context, in *emptypb.Empty) (*emptypb.Empty, error) {
	return nil, nil
}

func (s *server) Generate(_ context.Context, in *petnamepb.PetnameRequest) (*petnamepb.PetnameResponse, error) {
	if in.Words <= 0 {
		return nil, status.Error(codes.InvalidArgument, "words must be greater than 0")
	}

	name := petname.Generate(int(in.Words), in.Separator)

	return &petnamepb.PetnameResponse{Name: name}, nil
}

func (s *server) GenerateMany(in *petnamepb.PetnameStreamRequest, stream petnamepb.PetnameGenerator_GenerateManyServer) error {
	if in.Words <= 0 {
		return status.Error(codes.InvalidArgument, "words must be greater than 0")
	}

	if in.Names <= 0 {
		return status.Error(codes.InvalidArgument, "names must be greater than 0")
	}

	for i := int64(0); i < in.Names; i++ {
		name := petname.Generate(int(in.Words), in.Separator)
		if err := stream.Send(&petnamepb.PetnameResponse{Name: name}); err != nil {
			return err
		}
	}

	return nil
}

func main() {
	var configPath string
	flag.StringVar(&configPath, "config", "./config.yaml", "configuration file")
	flag.Parse()

	var cfg Config
	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		panic(err)
	}

	listener, err := net.Listen("tcp", ":"+strconv.Itoa(cfg.Port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	petnamepb.RegisterPetnameGeneratorServer(s, &server{})
	reflection.Register(s)

	if err := s.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
