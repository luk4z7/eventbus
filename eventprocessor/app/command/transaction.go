package command

import (
	"context"
	"encoding/json"
	"errors"
	"eventhandler/entities"
	"eventprocessor/domain/transaction"

	"github.com/ThreeDotsLabs/watermill/message"
)

func NewTransactionHandler(repo transaction.Repository, name string) *TransactionHandler {
	return &TransactionHandler{
		repo: repo,
		name: name,
	}
}

type TransactionHandler struct {
	repo transaction.Repository
	name string
}

func (t *TransactionHandler) HandlerName() string {
	return t.name
}

func (t *TransactionHandler) NewEvent() interface{} {
	return &message.Message{}
}

func (t *TransactionHandler) Handle(ctx context.Context, event any) error {
	msg, ok := event.(*message.Message)
	if !ok {
		return errors.New("this is not a *message.Message")
	}

	var data entities.Transaction
	if err := json.Unmarshal(msg.Payload, &data); err != nil {
		return err
	}

	return t.repo.Save(data)
}
