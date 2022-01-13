package repository

import (
	"fmt"
	"github.com/patrickmn/go-cache"
	"time"
)

type CacheRepositoryImpl struct {
	cache *cache.Cache
}

func (CacheRepositoryImpl) New() *CacheRepositoryImpl {
	c := cache.New(60*time.Second, 10*time.Second)
	return &CacheRepositoryImpl{c}
}

func (cacheRepositoryImpl *CacheRepositoryImpl) GetData(key string) string {
	data, found := cacheRepositoryImpl.cache.Get(key)
	if found {
		return data.(string)
	}
	return ""
}

func (cacheRepositoryImpl *CacheRepositoryImpl) SetData(key string, value string) {
	fmt.Printf("서버: 본서버에서 response가 왔으므로 {%s: %s}의 데이터를 캐시합니다\n", key, value)
	cacheRepositoryImpl.cache.Set(key, value, cache.DefaultExpiration)

}
