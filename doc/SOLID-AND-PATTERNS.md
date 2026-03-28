# SOLID E PATTERNS

## O Que Vamos Ver

- Como os princípios SOLID aparecem em Go.
- Quais patterns ajudam a aplicar cada princípio na prática.
- Exemplos curtos em Go.
- Links diretos para exemplos canônicos no Refactoring Guru.

## Ideia Central

SOLID e patterns não são a mesma coisa.

- `SOLID` ajuda a perceber a pressão de design.
- `patterns` ajudam a resolver essa pressão de forma conhecida.

Em Go, isso precisa ser feito com cuidado. A linguagem favorece:

- interfaces pequenas
- composição
- dependências explícitas
- menos cerimônia

Então o objetivo não é encher o projeto de abstrações. O objetivo é usar o pattern certo quando ele realmente melhora a clareza.

Os exemplos abaixo seguem uma linha mais direta e didática, no estilo do material que você trouxe sobre `Applying SOLID Principles in Golang`, e depois conectam cada princípio a patterns que ajudam na prática.

## 1. SRP: Single Responsibility Principle

Uma unidade de código deve ter um motivo claro para mudar.

### Exemplo base em Go

```go
package user

import "fmt"

type User struct {
	ID   int
	Name string
}

type UserService struct {
	users []User
}

func (us *UserService) AddUser(user User) {
	us.users = append(us.users, user)
}

func (us *UserService) GetUserByID(id int) (User, error) {
	for _, user := range us.users {
		if user.ID == id {
			return user, nil
		}
	}
	return User{}, fmt.Errorf("user not found")
}
```

Aqui o `UserService` tem uma responsabilidade clara: gerenciar usuários.

### O que quebraria SRP

Se o mesmo tipo também:

- enviasse e-mail
- escrevesse em banco diretamente
- formatasse resposta HTTP

então ele passaria a ter mais de um motivo para mudar.

### Exemplo melhor

```go
package user

import "fmt"

type User struct {
	ID   int
	Name string
}

type UserRepository struct {
	users []User
}

func (r *UserRepository) Save(user User) {
	r.users = append(r.users, user)
}

type UserService struct {
	repo *UserRepository
}

func NewUserService(repo *UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) AddUser(user User) {
	s.repo.Save(user)
}

func (s *UserService) GetUserByID(id int) (User, error) {
	for _, user := range s.repo.users {
		if user.ID == id {
			return user, nil
		}
	}
	return User{}, fmt.Errorf("user not found")
}
```

### Pattern que ajuda

- `Adapter`
- `Facade`

### Por quê

- `Adapter` separa transporte, banco e integrações do núcleo da aplicação.
- `Facade` ajuda a deixar uma interface simples para o cliente, enquanto a coordenação fica em outro lugar.

### Links úteis

- Adapter: https://refactoring.guru/pt-br/design-patterns/adapter/go/example
- Facade: https://refactoring.guru/pt-br/design-patterns/facade/go/example

## 2. OCP: Open/Closed Principle

Código deve estar aberto para extensão e fechado para modificação desnecessária.

### Exemplo em Go

```go
package user

type User struct {
	ID   int
	Name string
}

type UserService struct {
	users []User
}

func (us *UserService) AddUser(user User) {
	us.users = append(us.users, user)
}

type Authenticator interface {
	Authenticate(username, password string) bool
}

type AuthenticatedUserService struct {
	*UserService
	Authenticator
}

func (aus *AuthenticatedUserService) AddUserWithAuth(user User, username, password string) {
	if aus.Authenticate(username, password) {
		aus.AddUser(user)
	}
}
```

Aqui nós estendemos o comportamento sem modificar o `UserService` original.

### Exemplo mais idiomático com interfaces pequenas

```go
package payment

type Gateway interface {
	Charge(amount int64) error
}

type StripeGateway struct{}

func (StripeGateway) Charge(amount int64) error { return nil }

type PixGateway struct{}

func (PixGateway) Charge(amount int64) error { return nil }

type Service struct {
	gateway Gateway
}

func NewService(gateway Gateway) *Service {
	return &Service{gateway: gateway}
}

func (s *Service) Checkout(amount int64) error {
	return s.gateway.Charge(amount)
}
```

Você consegue adicionar um novo gateway sem reescrever o fluxo principal.

### Patterns que ajudam

- `Strategy`
- `Factory Method`
- `Abstract Factory`

### Por quê

- `Strategy` permite trocar algoritmos ou implementações por interface.
- `Factory Method` ajuda quando a criação muda por configuração.
- `Abstract Factory` ajuda quando você cria famílias de dependências relacionadas.

### Links úteis

- Strategy: https://refactoring.guru/pt-br/design-patterns/strategy/go/example
- Factory Method: https://refactoring.guru/pt-br/design-patterns/factory-method/go/example
- Abstract Factory: https://refactoring.guru/pt-br/design-patterns/abstract-factory/go/example

## 3. LSP: Liskov Substitution Principle

Se algo implementa uma abstração, deve poder substituí-la sem comportamento surpreendente.

### Exemplo em Go

```go
package shape

import "math"

type Shape interface {
	Area() float64
}

type Rectangle struct {
	Width  float64
	Height float64
}

func (r Rectangle) Area() float64 {
	return r.Width * r.Height
}

type Circle struct {
	Radius float64
}

func (c Circle) Area() float64 {
	return math.Pi * c.Radius * c.Radius
}
```

Quem usa `Shape` não precisa saber se recebeu `Rectangle` ou `Circle`.

### Patterns que ajudam

- `Strategy`
- `Adapter`

### Por quê

- `Strategy` só funciona bem quando as implementações realmente são substituíveis.
- `Adapter` evita jogar semânticas incompatíveis direto no domínio.

### Links úteis

- Strategy: https://refactoring.guru/pt-br/design-patterns/strategy/go/example
- Adapter: https://refactoring.guru/pt-br/design-patterns/adapter/go/example

## 4. ISP: Interface Segregation Principle

Clientes não devem depender de métodos que não usam.

### Exemplo ruim

```go
type OfficeMachine interface {
	Print()
	Scan()
	Fax()
}
```

Se um cliente só precisa imprimir, ele acaba dependendo de coisas que não usa.

### Exemplo melhor

```go
type Printer interface {
	Print()
}

type Scanner interface {
	Scan()
}

type FaxMachine interface {
	Fax()
}
```

Agora cada cliente escolhe apenas a interface que realmente precisa.

### Patterns que ajudam

- `Strategy`
- `Repository`

### Por quê

- interfaces menores ficam mais fáceis de trocar e testar
- `Repository` funciona melhor quando representa uma capacidade real, e não um balde genérico

`Repository` não é um pattern do Gang of Four, mas é extremamente útil em DDD e arquitetura hexagonal.

## 5. DIP: Dependency Inversion Principle

Módulos de alto nível devem depender de abstrações, não de detalhes concretos.

### Exemplo em Go

```go
package notification

import "fmt"

type Notifier interface {
	Notify(message string)
}

type EmailNotifier struct{}

func (EmailNotifier) Notify(message string) {
	fmt.Println("email:", message)
}

type SMSNotifier struct{}

func (SMSNotifier) Notify(message string) {
	fmt.Println("sms:", message)
}

type Service struct {
	notifier Notifier
}

func NewService(notifier Notifier) *Service {
	return &Service{notifier: notifier}
}

func (s *Service) Send(message string) {
	s.notifier.Notify(message)
}
```

Aqui o módulo de alto nível depende de `Notifier`, e não de `EmailNotifier` ou `SMSNotifier` diretamente.

### Patterns que ajudam

- `Adapter`
- `Strategy`
- `Repository`

### Por quê

- `Adapter` esconde detalhes técnicos atrás de uma interface
- `Strategy` permite trocar implementações sem mudar o fluxo
- `Repository` protege o domínio dos detalhes de persistência

### Links úteis

- Adapter: https://refactoring.guru/pt-br/design-patterns/adapter/go/example
- Strategy: https://refactoring.guru/pt-br/design-patterns/strategy/go/example

## Mapa Prático

### SOLID -> Patterns

- `SRP` -> `Adapter`, `Facade`
- `OCP` -> `Strategy`, `Factory Method`, `Abstract Factory`
- `LSP` -> `Strategy`, `Adapter`
- `ISP` -> `Strategy`, `Repository`
- `DIP` -> `Adapter`, `Strategy`, `Repository`

## Como Pensar Isso No Projeto

No contexto deste repositório:

- handlers HTTP funcionam como `Adapter`
- clients de banco e serviços remotos também funcionam como `Adapter`
- ports e implementações trocáveis se aproximam de `Strategy`
- o `mobile-bff` se aproxima de `Facade`
- repositórios protegem os casos de uso dos detalhes de persistência

## Referência Principal

- Refactoring Guru, catálogo Go: https://refactoring.guru/pt-br/design-patterns/go

## Recapitulação

Patterns não devem ser usados como enfeite. Em Go, eles valem a pena quando tornam fronteiras mais claras, facilitam troca de implementação e reduzem acoplamento real. O melhor uso de `SOLID` com `patterns` é pragmático: menos cerimônia, mais clareza.
