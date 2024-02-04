package service

import (
	"context"
	"eventprocessor/adapters"
	"eventprocessor/app"
	"eventprocessor/app/command"
	"os"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-redisstream/pkg/redisstream"
	"github.com/ThreeDotsLabs/watermill/components/cqrs"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/redis/go-redis/v9"
)

func NewApplication(
	ctx context.Context,
	router *message.Router,
	logger watermill.LoggerAdapter,

) (app.Application, func()) {

	rdb := redis.NewClient(&redis.Options{
		Addr: os.Getenv("REDIS_ADDR"),
	})

	ep, err := cqrs.NewEventProcessorWithConfig(
		router,
		cqrs.EventProcessorConfig{
			SubscriberConstructor: func(params cqrs.EventProcessorSubscriberConstructorParams) (message.Subscriber, error) {
				return redisstream.NewSubscriber(redisstream.SubscriberConfig{
					Client:        rdb,
					ConsumerGroup: "issue-receipt",
				}, logger)
			},
			GenerateSubscribeTopic: func(params cqrs.EventProcessorGenerateSubscribeTopicParams) (string, error) {
				return "TransactionConfirmed", nil
			},
			Marshaler: cqrs.JSONMarshaler{},
			Logger:    logger,
		},
	)
	if err != nil {
		panic(err)
	}

	repo := adapters.NewTransactionRepository()

	application := app.Application{
		Commands: app.Commands{
			Transaction: command.NewTransactionHandler(repo, "TransactionConfirmed"),
		},
		Queries: app.Queries{},
	}

	if err := ep.AddHandlers(application.Commands.Transaction); err != nil {
		panic(err)
	}

	return application, func() {}
}
