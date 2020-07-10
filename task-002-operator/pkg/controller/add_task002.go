package controller

import (
	"github.com/keunlee/task-002-operator/pkg/controller/task002"
)

func init() {
	// AddToManagerFuncs is a list of functions to create controllers and add them to a manager.
	AddToManagerFuncs = append(AddToManagerFuncs, task002.Add)
}
