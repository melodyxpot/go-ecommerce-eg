package main

import (
	"context"
	"os"
	"time"

	rocketmqLib "github.com/apache/rocketmq-client-go/v2"
	rocketConsumer "github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/labstack/echo/v4"

	"github.com/smiletrl/micro_ecommerce/pkg/accesslog"
	"github.com/smiletrl/micro_ecommerce/pkg/config"
	"github.com/smiletrl/micro_ecommerce/pkg/constants"
	"github.com/smiletrl/micro_ecommerce/pkg/errors"
	"github.com/smiletrl/micro_ecommerce/pkg/healthcheck"
	"github.com/smiletrl/micro_ecommerce/pkg/logger"
	"github.com/smiletrl/micro_ecommerce/pkg/postgres"
	"github.com/smiletrl/micro_ecommerce/pkg/rocketmq"
	"github.com/smiletrl/micro_ecommerce/pkg/tracing"
	"github.com/smiletrl/micro_ecommerce/service.order/internal/order"
)

type provider struct {
	config   config.Config
	consumer rocketmqLib.PushConsumer
	rocketmq rocketmq.Provider
	tracing  tracing.Provider
	pdb      postgres.Provider
	logger   logger.Provider
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

	consumer, err := rocketmqProvider.CreatePushConsumer(ctx, constants.RocketMQGroupPayment, rocketConsumer.Clustering)
	if err != nil {
		panic(err)
	}
	defer rocketmqProvider.ShutdownPushConsumer(consumer)

	p := provider{
		config:   cfg,
		tracing:  tracing,
		logger:   logger,
		consumer: consumer,
		pdb:      pdb,
		rocketmq: rocketmqProvider,
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

	// order message
	orderRepo := order.NewRepository(p.pdb)
	orderService := order.NewService(orderRepo, p.logger)
	orderMessage := order.NewMessage(p.consumer, orderRepo, p.rocketmq, p.tracing, p.logger)
	// start the subscribe
	if err := orderMessage.Subscribe(); err != nil {
		panic(err)
	}
	// start the consumer. Consumer needs to start after subscribe.
	if err := p.rocketmq.StartPushConsumer(p.consumer); err != nil {
		panic(err)
	}
	group := e.Group("/api/v1")
	order.RegisterHandlers(group, orderService)

	// Start rest server
	panic(e.Start(constants.RestPort))
}
