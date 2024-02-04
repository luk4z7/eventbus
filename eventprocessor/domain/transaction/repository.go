package transaction

import "eventhandler/entities"

type Repository interface {
	Save(entry entities.Transaction) error
}
