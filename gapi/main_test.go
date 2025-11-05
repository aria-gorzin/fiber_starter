package gapi

import (
	"testing"
	"time"

	db "github.com/aria/app/db/sqlc"
	"github.com/aria/app/util"
	"github.com/aria/app/worker"
	"github.com/stretchr/testify/require"
)

func newTestServer(t *testing.T, store db.Store, taskDistributor worker.TaskDistributor) *Server {
	config := util.Config{
		TokenSymmetricKey:   util.RandomString(32),
		AccessTokenDuration: time.Minute,
	}

	server, err := NewServer(config, store, taskDistributor)
	require.NoError(t, err)

	return server
}

// func newContextWithBearerToken(t *testing.T, tokenMaker token.PasetoMaker, username string, role string, duration time.Duration, tokenType string) context.Context {
// 	accessToken, _, err := tokenMaker.CreateToken(username, role, duration, tokenType)
// 	require.NoError(t, err)

// 	bearerToken := fmt.Sprintf("%s %s", authorizationBearer, accessToken)
// 	md := metadata.MD{
// 		authorizationHeader: []string{
// 			bearerToken,
// 		},
// 	}

// 	return metadata.NewIncomingContext(context.Background(), md)
// }
