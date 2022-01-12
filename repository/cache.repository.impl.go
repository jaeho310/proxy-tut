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
	c := cache.New(60*time.Second, 1*time.Second)
	return &CacheRepositoryImpl{c}
}

func (cacheRepositoryImpl *CacheRepositoryImpl) GetData(key string) string {
	data, found := cacheRepositoryImpl.cache.Get(key)
	if found {
		fmt.Printf("서버: [%s]는 캐시되어있으므로 [%s]를 내려줍니다\n", key, data)
		return data.(string)
	}
	return ""
}

func (cacheRepositoryImpl *CacheRepositoryImpl) SetData(key string, value string) {
	fmt.Printf("서버: 타서버에서 response가 왔으므로 {%s: %s}의 데이터를 캐시합니다\n", key, value)
	cacheRepositoryImpl.cache.Set(key, value, cache.DefaultExpiration)

}
