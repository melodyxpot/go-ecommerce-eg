package main

import (
	"context"
	"os"
	"time"

	rocketmqLib "github.com/apache/rocketmq-client-go/v2"
	"github.com/labstack/echo/v4"

	"github.com/smiletrl/micro_ecommerce/pkg/accesslog"
	"github.com/smiletrl/micro_ecommerce/pkg/config"
	"github.com/smiletrl/micro_ecommerce/pkg/constants"
	"github.com/smiletrl/micro_ecommerce/pkg/errors"
	"github.com/smiletrl/micro_ecommerce/pkg/healthcheck"
	"github.com/smiletrl/micro_ecommerce/pkg/logger"
	"github.com/smiletrl/micro_ecommerce/pkg/postgres"
	"github.com/smiletrl/micro_ecommerce/pkg/redis"
	"github.com/smiletrl/micro_ecommerce/pkg/rocketmq"
	"github.com/smiletrl/micro_ecommerce/pkg/tracing"
	"github.com/smiletrl/micro_ecommerce/service.payment/internal/payment"
)

type provider struct {
	config   config.Config
	tracing  tracing.Provider
	logger   logger.Provider
	rocketmq rocketmqLib.Producer
	pdb      postgres.Provider
	rdb      redis.Provider
}

func main() {
	// stage
	stage := os.Getenv(constants.Stage)
	if stage == "" {
		stage = constants.StageLocal
	}

	// init config
	cfg, err := config.Load(stage)
	if err != nil {
		panic(err)
	}

	// init logger
	logger := logger.NewProvider(cfg.Logger)
	defer logger.Close()

	// init tracing
	tracing, err := tracing.NewProvider(constants.TracingCart, cfg)
	if err != nil {
		panic(err)
	}
	defer tracing.Close()

	// init postgres
	pdb, err := postgres.NewProvider(cfg, tracing)
	if err != nil {
		panic(err)
	}
	defer pdb.Close()

	// init rocketmq
	rocketmqProvider := rocketmq.NewProvider(cfg.RocketMQ, pdb, constants.ServiceOrder)
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	producer, err := rocketmqProvider.CreateProducer(ctx, constants.RocketMQGroupPayment)
	if err != nil {
		panic(err)
	}
	defer rocketmqProvider.ShutdownProducer(producer)

	// redis
	rdb, err := redis.NewProvider(cfg, tracing)
	if err != nil {
		panic(err)
	}

	p := provider{
		config:   cfg,
		tracing:  tracing,
		logger:   logger,
		rocketmq: producer,
		pdb:      pdb,
		rdb:      rdb,
	}
	buildRegisters(p)
}

func buildRegisters(p provider) {
	// echo instance
	e := echo.New()

	// middleware
	e.Use(accesslog.Middleware(p.logger))
	e.Use(p.tracing.Middleware(p.logger))
	e.Use(errors.Middleware(p.logger))

	// initialize health
	healthcheck.RegisterHandlers(e.Group(""))

	group := e.Group("/api/v1")

	paymentRepo := payment.NewRepository(p.pdb, p.rdb)
	paymentService := payment.NewService(paymentRepo, p.rocketmq, p.tracing, p.logger)
	payment.RegisterHandlers(group, paymentService)

	// Start rest server
	panic(e.Start(constants.RestPort))
}
