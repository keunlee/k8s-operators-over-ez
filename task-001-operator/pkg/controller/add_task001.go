package controller

import (
	"github.com/keunlee/task-001-operator/pkg/controller/task001"
)

func init() {
	// AddToManagerFuncs is a list of functions to create controllers and add them to a manager.
	AddToManagerFuncs = append(AddToManagerFuncs, task001.Add)
}
