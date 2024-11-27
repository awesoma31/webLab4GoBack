package model

import "context"

type AuthService interface {
	Authorize(ctx context.Context) error
}
