package socialservice

import (
	uow "github.com/ZaiiiRan/messenger/backend/social-service/internal/repositories/unitofwork/postgres"
	"github.com/ZaiiiRan/messenger/backend/social-service/internal/transport/postgres"
)

type socialDataProvider struct {
	pg *postgres.PostgresClient
}

func newSocialDataProvider(pg *postgres.PostgresClient) *socialDataProvider {
	return &socialDataProvider{
		pg: pg,
	}
}

func (udp *socialDataProvider) newUOW() *uow.UnitOfWork {
	return uow.New(udp.pg)
}
