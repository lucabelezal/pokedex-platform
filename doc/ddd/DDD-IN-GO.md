# DDD Em Go

## O Que Vamos Ver

- Como os conceitos de DDD podem ser expressos em Go.
- Exemplos de entities, value objects, repositories e services.
- Por que Go se beneficia de uma interpretação mais leve de DDD.
- O que evitar ao misturar DDD com Go idiomático.

## Introdução

DDD não exige um estilo pesado de orientação a objetos. Suas ideias centrais são conceituais, então funcionam em Go desde que preservemos o significado:

- conceitos de negócio devem ser explícitos
- o modelo deve carregar significado real de negócio
- as fronteiras devem ser claras
- preocupações técnicas não devem dominar o modelo central

Go normalmente incentiva um estilo de implementação mais leve:

- interfaces menores
- dependências explícitas
- composição em vez de herança
- menos camadas cerimoniais

Isso costuma combinar bem com um DDD pragmático.

## Exemplo De Entity

Uma entity tem identidade e comportamento relevante para o negócio.

```go
package domain

type OrderStatus string

const (
	OrderPending  OrderStatus = "pending"
	OrderPaid     OrderStatus = "paid"
	OrderCanceled OrderStatus = "canceled"
)

type Order struct {
	ID     string
	Status OrderStatus
	Total  int64
}

func (o *Order) Pay() error {
	if o.Status != OrderPending {
		return ErrInvalidStatusTransition
	}
	o.Status = OrderPaid
	return nil
}
```

O que importa aqui não são apenas os dados. A entity também protege regras de negócio.

## Exemplo De Value Object

Um value object é definido por seu valor e frequentemente carrega validação ou invariantes.

```go
package domain

type Email string

func NewEmail(value string) (Email, error) {
	if value == "" {
		return "", ErrInvalidEmail
	}
	return Email(value), nil
}
```

Outro exemplo clássico é `Money`, normalmente modelado com valor e moeda.

## Exemplo De Repository

O repository é um contrato voltado ao domínio para persistência.

```go
package ports

import (
	"context"
	"myapp/domain"
)

type OrderRepository interface {
	FindByID(ctx context.Context, id string) (*domain.Order, error)
	Save(ctx context.Context, order *domain.Order) error
}
```

Isso é útil porque o caso de uso pode falar em termos de negócio sem conhecer SQL, HTTP ou detalhes de armazenamento.

## Exemplo De Service De Aplicação

O service de aplicação coordena um caso de uso.

```go
package service

import (
	"context"
	"myapp/ports"
)

type OrderService struct {
	repo ports.OrderRepository
}

func NewOrderService(repo ports.OrderRepository) *OrderService {
	return &OrderService{repo: repo}
}

func (s *OrderService) PayOrder(ctx context.Context, id string) error {
	order, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}

	if err := order.Pay(); err != nil {
		return err
	}

	return s.repo.Save(ctx, order)
}
```

O service orquestra. A regra de negócio continua pertencendo ao modelo.

## Pensando Em Aggregate

Um aggregate é uma fronteira de consistência, e não apenas um grupo de structs.

Modelo mental simples:

- uma raiz controla as mudanças
- as invariantes são protegidas por essa raiz
- código externo não deve mutar livremente as partes internas

Isso importa mais em domínios transacionais ricos do que em código CRUD simples.

## O Que Evitar Em Go

### Modelos Anêmicos Em Todo Lugar

Se todo tipo de domínio for apenas um saco de campos, a lógica de negócio frequentemente vaza para services, handlers ou repositories.

### Interfaces Gigantes

Go funciona melhor com contratos menores e focados.

### Forçar Todo Artefato De DDD

Nem todo projeto precisa de aggregates elaborados, factories, specifications e domain events todos de uma vez.

### Abstrair Demais Antes De Existir Dor

DDD deve ajudar a tornar a complexidade visível. Não deve fabricar complexidade.

## Uma Boa Heurística Prática

Em Go, um modelo saudável inspirado em DDD costuma se parecer com isto:

- nomes com significado
- invariantes explícitas
- repositories focados
- services de aplicação coordenando casos de uso
- transporte e persistência mantidos fora do modelo central

Isso já é suficiente para obter valor real sem perder clareza.

## Recapitulação

DDD em Go funciona melhor quando permanece focado em significado e fronteiras, e não em cerimônia. Entities devem proteger comportamento de negócio, value objects devem carregar significado invariável, repositories devem descrever necessidades de persistência na linguagem do domínio, e services de aplicação devem orquestrar casos de uso sem engolir o domínio.
