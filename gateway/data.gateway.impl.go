package gateway

import (
	"errors"
	"fmt"
	"time"
)

type DataGatewayImpl struct {
}

func (DataGatewayImpl) New() *DataGatewayImpl {
	return &DataGatewayImpl{}
}

func (*DataGatewayImpl) GetDataFromDataServer(param string) (string, error) {
	// 해당 작업은 Data서버에 부하를 많이 준다고 가정
	fmt.Printf("타서버: [%s] 라는 요청을 받았습니다, 데이터를 내려주기까지 1초 소요됩니다. \n", param)
	<-time.After(time.Second * 1)
	if param == "hello" {
		return "world", nil
	} else if param == "foo" {
		return "bar", nil
	}
	return "", errors.New("잘못된 파라미터입니다")
}
