package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/ilyakaznacheev/cleanenv"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	wordspb "yadro.com/course/proto/words"
	"yadro.com/course/words/words"
)

const phraseSizeLimit = 4 * 1024

type Config struct {
	Address string `yaml:"words_address" env:"WORDS_ADDRESS" env-default:"80"`
}

type server struct {
	wordspb.UnimplementedWordsServer
	logger *slog.Logger
}

func (s *server) Ping(_ context.Context, in *emptypb.Empty) (*emptypb.Empty, error) {
	return nil, nil
}

func (s *server) Norm(_ context.Context, in *wordspb.WordsRequest) (*wordspb.WordsReply, error) {
	s.logger.Debug("norm", "phrase", in.GetPhrase())

	if len(in.GetPhrase()) > phraseSizeLimit {
		s.logger.Error("input phrase exceeds 4KiB", "len(in.Phrase)", len(in.Phrase))

		return nil, status.Error(codes.ResourceExhausted, "input phrase exceeds 4KiB")
	}

	stemsOfWords := words.Norm(in.GetPhrase(), s.logger)

	return &wordspb.WordsReply{
		Words: stemsOfWords,
	}, nil
}

func main() {
	logger := slog.New(
		slog.NewTextHandler(
			os.Stdout,
			&slog.HandlerOptions{
				AddSource: true,
				Level:     slog.LevelDebug,
			},
		),
	)

	var configPath string
	flag.StringVar(&configPath, "config", "config.yaml", "configuration file")
	flag.Parse()

	logger.Debug("read config flag", "configPath", configPath)

	var cfg Config
	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		panic(err)
	}

	if err := run(cfg, logger); err != nil {
		logger.Error("failed to run server", "error", err)
		os.Exit(1)
	}
}

func run(cfg Config, logger *slog.Logger) error {
	listener, err := net.Listen("tcp", cfg.Address)
	if err != nil {
		return fmt.Errorf("failed to listen %s: %v", cfg.Address, err)
	}

	s := grpc.NewServer()
	wordspb.RegisterWordsServer(s, &server{logger: logger})
	reflection.Register(s)

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	group, ctx := errgroup.WithContext(ctx)

	group.Go(
		func() error {
			<-ctx.Done()
			logger.Info("shutdown started")
			s.GracefulStop()
			logger.Info("shutdown stopped")
			return nil
		},
	)

	group.Go(
		func() error {
			logger.Debug("starting server", "adress", cfg.Address)
			return s.Serve(listener)
		},
	)

	return group.Wait()
}
