package app

import "eventprocessor/app/command"

type Application struct {
	Commands Commands
	Queries  Queries
}

type Commands struct {
	Transaction *command.TransactionHandler
}

type Queries struct{}
