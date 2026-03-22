package domain

import "errors"

var (
	ErrPokemonNotFound       = errors.New("pokemon nao encontrado")
	ErrUserNotFound          = errors.New("usuario nao encontrado")
	ErrFavoriteNotFound      = errors.New("favorito nao encontrado")
	ErrFavoriteAlreadyExists = errors.New("favorito ja existe")
	ErrInvalidPagination     = errors.New("parametros de paginacao invalidos")
	ErrUnauthorized          = errors.New("nao autorizado")
	ErrInvalidToken          = errors.New("token invalido")
)
