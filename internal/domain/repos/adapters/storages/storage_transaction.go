package storages

import "context"

type StorageTransaction interface {
	InvokeTransactionMechanism(ctx context.Context) (interface{}, error)
	ShadowTransactionMechanism(ctx context.Context, transaction interface{}) error
}
