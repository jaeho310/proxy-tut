package model

import "context"

type ResponseCtxModel struct {
	Value  string
	Ctx    context.Context
	Cancel context.CancelFunc
}
