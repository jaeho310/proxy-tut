package service

import "proxy-tut/gateway"

type DataServiceImpl struct {
	gateway.DataGateway
}

func (dataServiceImpl DataServiceImpl) New(dataGateway gateway.DataGateway) *DataServiceImpl {
	return &DataServiceImpl{dataGateway}
}

func (dataServiceImpl *DataServiceImpl) GetDataWithParam(param string) (string, error) {
	// 동일한 조회 요청이 왔을때도 gateway를 통해 본서버로 요청이 들어가게 된다.
	dataFromSomeServer, err := dataServiceImpl.DataGateway.GetDataFromDataServer(param)
	if err != nil {
		return "", err
	}
	return dataFromSomeServer, nil
}
