package service

import (
	job "github.com/apk471/go-boilerplate/internal/lib/jobs"
	"github.com/apk471/go-boilerplate/internal/repository"
	"github.com/apk471/go-boilerplate/internal/server"
)

type Services struct {
	Auth *AuthService
	Job  *job.JobService
	User *UserService
}

func NewServices(s *server.Server, repos *repository.Repositories) (*Services, error) {
	authService := NewAuthService(s)
	userService := NewUserService(s, repos.User)

	return &Services{
		Job:  s.Job,
		Auth: authService,
		User: userService,
	}, nil
}
