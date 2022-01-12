package service

type DataService interface {
	GetDataWithParam(param string) (string, error)
}
