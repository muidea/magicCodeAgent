package model

import "github.com/muidea/magicBatis/client"

var modelList = []interface{}{
	&Log{},
	&Totalizer{},
}

// InitializeModel initialize model
func InitializeModel(clnt client.Client) (err error) {
	for _, val := range modelList {
		err = clnt.RegisterModel(val)
		if err != nil {
			break
		}
	}
	if err != nil {
		return
	}

	err = CreateModel(clnt)
	if err != nil {
		return err
	}

	return
}

// CreateModel create model
func CreateModel(clnt client.Client) (err error) {
	for _, val := range modelList {
		err = clnt.CreateSchema(val)
		if err != nil {
			break
		}
	}

	return
}

// DropModel drop model
func DropModel(clnt client.Client) (err error) {
	for _, val := range modelList {
		err = clnt.DropSchema(val)
		if err != nil {
			break
		}
	}

	return
}

// UninitializeModel uninitialize model
func UninitializeModel(clnt client.Client) (err error) {
	err = DropModel(clnt)
	if err != nil {
		return
	}

	for _, val := range modelList {
		err = clnt.UnregisterModel(val)
		if err != nil {
			break
		}
	}

	return
}
