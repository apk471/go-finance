package repository

import "github.com/apk471/go-boilerplate/internal/server"

type Repositories struct {
	User *UserRepository
}

func NewRepositories(s *server.Server) *Repositories {
	return &Repositories{
		User: NewUserRepository(s),
	}
}
