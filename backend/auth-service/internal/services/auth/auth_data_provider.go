package authservice

import (
	uow "github.com/ZaiiiRan/messenger/backend/auth-service/internal/repositories/unitofwork/postgres"
	"github.com/ZaiiiRan/messenger/backend/auth-service/internal/transport/postgres"
)

type authDataProvider struct {
	pg *postgres.PostgresClient
}

func newAuthDataProvider(pg *postgres.PostgresClient) *authDataProvider {
	return &authDataProvider{
		pg: pg,
	}
}

func (adp *authDataProvider) newUOW() *uow.UnitOfWork {
	return uow.New(adp.pg)
}
