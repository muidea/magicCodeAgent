package dao

import (
	"github.com/muidea/magicBatis/client"
)

type Setting interface {
}

// New 新建Base
func New(clnt client.Client) Setting {
	return &setting{batisClient: clnt}
}

type setting struct {
	batisClient client.Client
}
