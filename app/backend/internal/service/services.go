package service

import (
	job "github.com/apk471/go-boilerplate/internal/lib/jobs"
	"github.com/apk471/go-boilerplate/internal/repository"
	"github.com/apk471/go-boilerplate/internal/server"
)

type Services struct {
	Auth            *AuthService
	Job             *job.JobService
	User            *UserService
	FinancialRecord *FinancialRecordService
}

func NewServices(s *server.Server, repos *repository.Repositories) (*Services, error) {
	authService := NewAuthService(s)
	userService := NewUserService(s, repos.User)
	financialRecordService := NewFinancialRecordService(repos.FinancialRecord)

	return &Services{
		Job:             s.Job,
		Auth:            authService,
		User:            userService,
		FinancialRecord: financialRecordService,
	}, nil
}
