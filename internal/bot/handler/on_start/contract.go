package on_start

import "context"

type storage interface {
	CreateUser(ctx context.Context, userId int64) error
}
