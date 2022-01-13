package service

import (
	"fmt"
	"proxy-tut/gateway"
	"proxy-tut/repository"
)

type ProxyDataServiceImpl struct {
	dataGateway         gateway.DataGateway
	temporaryRepository repository.TemporaryRepository
	cacheRepository     repository.CacheRepository
}

func (proxyServiceImpl ProxyDataServiceImpl) New(dataGateway gateway.DataGateway) *ProxyDataServiceImpl {
	temporaryRepositoryImpl := repository.TemporaryRepositoryImpl{}.New()
	cacheRepositoryImpl := repository.CacheRepositoryImpl{}.New()
	return &ProxyDataServiceImpl{dataGateway, temporaryRepositoryImpl, cacheRepositoryImpl}
}

func (proxyServiceImpl *ProxyDataServiceImpl) GetDataWithParam(param string) (string, error) {
	var responseData string
	var err error
	// 캐시된 요청이면 바로 응답을 내려준다.
	responseData = proxyServiceImpl.cacheRepository.GetData(param)
	if len(responseData) != 0 {
		fmt.Printf("프록시서버: [%s]는 캐시되어있으므로 [%s]를 내려줍니다\n", param, responseData)
		return responseData, nil
	}
	// 캐시되지 않은 요청이면 최초의 요청인지 확인
	first, responseModelWithCtx := proxyServiceImpl.temporaryRepository.GetResponseModel(param)
	defer responseModelWithCtx.Cancel()
	if first {
		// 최초의 요청인 경우
		fmt.Printf("프록시서버: [%s] 요청을 본서버에 보냈습니다. "+
			"2초동안 동일한 요청은 본서버에 보내지 않습니다 \n", param)
		responseData, err = proxyServiceImpl.dataGateway.GetDataFromDataServer(param)
		if err != nil {
			return "", err
		}
		responseModelWithCtx.Value = responseData
		// 최초의 요청이 본서버에서 데이터를 받으면 나머지 요청들과 데이터를 공유해 내려준다.
		responseModelWithCtx.Cancel()

		// 클라이언트에게 응답은 바로 가야하므로 캐시하는 로직은 고루틴으로 처리
		go func() {
			proxyServiceImpl.cacheRepository.SetData(param, responseData)
		}()
	} else {
		fmt.Printf("프록시서버: [%s] 요청은 이미 본서버에 보낸 요청입니다."+
			" 응답이 오면 동일한 결과를 내립니다. \n", param)
		// 최초의 요청이 아니라면 2초동안은 본서버에 request를 보내지않고 block
		<-responseModelWithCtx.Ctx.Done()
		responseData = responseModelWithCtx.Value
	}

	return responseData, nil
}
