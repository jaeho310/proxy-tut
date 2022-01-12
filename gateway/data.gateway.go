package gateway

type DataGateway interface {
	GetDataFromDataServer(param string) (string, error)
}
