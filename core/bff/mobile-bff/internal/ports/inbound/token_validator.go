package inbound

import "context"

// TokenValidator verifica se um token de acesso ainda está ativo no auth-service.
type TokenValidator interface {
	IsTokenActive(ctx context.Context, token string) (bool, error)
}
