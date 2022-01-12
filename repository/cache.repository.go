package repository

type CacheRepository interface {
	GetData(key string) string
	SetData(key string, value string)
}
