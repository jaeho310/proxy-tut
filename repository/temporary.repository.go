package repository

import "proxy-tut/model"

type TemporaryRepository interface {
	setData(key string, value *model.ResponseCtxModel)
	GetResponseModel(key string) (bool, *model.ResponseCtxModel)
}
