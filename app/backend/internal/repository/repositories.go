package repository

import "github.com/apk471/go-boilerplate/internal/server"

type Repositories struct {
	User            *UserRepository
	FinancialRecord *FinancialRecordRepository
}

func NewRepositories(s *server.Server) *Repositories {
	return &Repositories{
		User:            NewUserRepository(s),
		FinancialRecord: NewFinancialRecordRepository(s),
	}
}
