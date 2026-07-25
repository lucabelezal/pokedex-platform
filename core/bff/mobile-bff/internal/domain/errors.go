package domain

import "errors"

var (
	ErrPokemonNotFound       = errors.New("pokemon nao encontrado")
	ErrUserAlreadyExists     = errors.New("usuario ja existe")
	ErrFavoriteNotFound      = errors.New("favorito nao encontrado")
	ErrFavoriteAlreadyExists = errors.New("favorito ja existe")
	ErrUnauthorized          = errors.New("nao autorizado")
	ErrInvalidCredentials    = errors.New("credenciais invalidas")
	ErrInvalidInput          = errors.New("dados invalidos")
	ErrInvalidToken          = errors.New("token invalido")
	ErrAuthUnavailable       = errors.New("auth service indisponivel")
	ErrServiceUnavailable    = errors.New("servico temporariamente indisponivel")
)
