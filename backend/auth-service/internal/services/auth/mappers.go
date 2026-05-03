package authservice

import (
	pb "github.com/ZaiiiRan/messenger/backend/auth-service/gen/go/auth/v1"
	"github.com/ZaiiiRan/messenger/backend/auth-service/internal/domain/token"
	"github.com/ZaiiiRan/messenger/backend/auth-service/internal/utils"
)

func toSessionPb(refreshToken *token.Token) *pb.Session {
	return &pb.Session{
		Id:        refreshToken.GetID(),
		Ip:        refreshToken.GetIP(),
		Country:   refreshToken.GetCountry(),
		City:      refreshToken.GetCity(),
		Os:        refreshToken.GetOS(),
		Browser:   refreshToken.GetBrowser(),
		CreatedAt: utils.FormatTimestamp(refreshToken.GetCreatedAt()),
	}
}
