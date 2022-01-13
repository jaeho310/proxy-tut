package main

import (
	"fmt"
	"proxy-tut/gateway"
	"proxy-tut/service"
	"time"
)

func main() {
	dataGateway := gateway.DataGatewayImpl{}.New()
	proxyService := service.ProxyDataServiceImpl{}.New(dataGateway)

	runClient(proxyService)

	<-time.After(time.Second * 10)
}

func runClient(proxyService *service.ProxyDataServiceImpl) {
	for i := 0; i < 5; i++ {
		go func() {
			data, err := proxyService.GetDataWithParam("hello")
			if err != nil {
				panic(err)
			}
			fmt.Println("클라이언트: 결과 데이터:", data)
		}()
		<-time.After(time.Millisecond * 400)
	}
}
