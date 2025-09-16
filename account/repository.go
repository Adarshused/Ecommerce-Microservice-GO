package account

import (
	"context"
	"database/sql"
	 _ "github.com/lib/pq" // protocol buffer
)

type Respository interface {
	close()
	PutAccount(ctx context.Context, a Account) error
	GetAccountByID(ctx context.Context, id string) (*Account, error)
}