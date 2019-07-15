package service

//mockery -name=UserRepository
type UserRepository interface {
	Repository
	GetByEmail(email string) (*User, error)
	GetByLogin(login string) (*User, error)
	FindOneByLogin(login string) (*User, error)
	FindOneByEmail(login string) (*User, error)
}
