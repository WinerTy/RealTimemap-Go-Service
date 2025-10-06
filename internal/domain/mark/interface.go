package mark

import "context"

type Repository interface {
	GetByOwner(ctx context.Context, ownerID int) ([]*Mark, error)
}
