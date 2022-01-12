package repository

import (
	"context"
	"github.com/patrickmn/go-cache"
	"proxy-tut/model"
	"sync"
	"time"
)

type TemporaryRepositoryImpl struct {
	cache *cache.Cache
	mu    *sync.Mutex
}

func (TemporaryRepositoryImpl) New() *TemporaryRepositoryImpl {
	c := cache.New(2*time.Second, 1*time.Second)
	return &TemporaryRepositoryImpl{c, &sync.Mutex{}}
}

func (temporaryRepositoryImpl *TemporaryRepositoryImpl) setData(key string, value *model.ResponseCtxModel) {
	temporaryRepositoryImpl.cache.Set(key, value, cache.DefaultExpiration)
}

func (temporaryRepositoryImpl *TemporaryRepositoryImpl) GetResponseModel(key string) (bool, *model.ResponseCtxModel, error) {
	temporaryRepositoryImpl.mu.Lock()
	defer temporaryRepositoryImpl.mu.Unlock()
	data, found := temporaryRepositoryImpl.cache.Get(key)
	if found {
		return true, data.(*model.ResponseCtxModel), nil
	} else {
		ctx, cancel := context.WithCancel(context.Background())
		temporaryModel := &model.ResponseCtxModel{Value: "", Cancel: cancel, Ctx: ctx}
		temporaryRepositoryImpl.setData(key, temporaryModel)
		return false, temporaryModel, nil
	}

}
