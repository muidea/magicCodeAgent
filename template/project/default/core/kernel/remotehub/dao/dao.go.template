package dao

import "github.com/muidea/magicBatis/client"

type RemoteHub interface {
}

func New(clnt client.Client) RemoteHub {
	return &remotehub{batisClient: clnt}
}

type remotehub struct {
	batisClient client.Client
}
