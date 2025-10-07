package mark

import (
	"context"
)

type Repository interface {
	GetByOwner(ctx context.Context, ownerID int) ([]*Mark, error)
	GetNearestMarks(ctx context.Context, filter Filter) ([]*Mark, error)
}

type Service interface {
}
