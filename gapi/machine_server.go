package gapi

import (
	"fmt"

	db "github.com/aria/app/db/sqlc"
	"github.com/aria/app/pb"
	"github.com/aria/app/token"
	"github.com/aria/app/util"
	"github.com/aria/app/worker"
)

// Server serves gRPC requests for our banking service.
type MachineServer struct {
	pb.UnimplementedMachinesServer
	config          util.Config
	store           db.Store
	tokenMaker      token.PasetoMaker
	taskDistributor worker.TaskDistributor
}

// NewServer creates a new gRPC server.
func NewMachineServer(config util.Config, store db.Store, taskDistributor worker.TaskDistributor) (*MachineServer, error) {
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %w", err)
	}

	server := &MachineServer{
		config:          config,
		store:           store,
		tokenMaker:      *tokenMaker,
		taskDistributor: taskDistributor,
	}

	return server, nil
}
