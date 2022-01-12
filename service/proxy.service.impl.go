package service

import (
	"fmt"
	"proxy-tut/gateway"
	"proxy-tut/repository"
)

type ProxyServiceImpl struct {
	dataGateway         gateway.DataGateway
	temporaryRepository repository.TemporaryRepository
	cacheRepository     repository.CacheRepository
}

func (proxyServiceImpl ProxyServiceImpl) New(dataGateway gateway.DataGateway) *ProxyServiceImpl {
	temporaryRepositoryImpl := repository.TemporaryRepositoryImpl{}.New()
	cacheRepositoryImpl := repository.CacheRepositoryImpl{}.New()
	return &ProxyServiceImpl{dataGateway, temporaryRepositoryImpl, cacheRepositoryImpl}
}

func (proxyServiceImpl *ProxyServiceImpl) GetDataWithParam(param string) (string, error) {
	// 동일한 조회 요청이 왔을때 타 서버에 동일한 요청을 보내지 않기위해 동일한 요청은 block
	var responseData string
	responseData = proxyServiceImpl.cacheRepository.GetData(param)
	if len(responseData) != 0 {
		return responseData, nil
	} else {
		isContain, responseModelWithCtx, err := proxyServiceImpl.temporaryRepository.GetResponseModel(param)
		defer responseModelWithCtx.Cancel()
		if err != nil {
			return "", err
		}
		if isContain {
			fmt.Printf("서버: [%s] 는 동일한 요청이므로 타 서버에 요청을 보내지 않습니다. 하나의 응답이 올때까지 블락시킵니다\n", param)
			// 하나의 요청이 온 후 3초동안은 동일한 요청이 와도 문맥종료까지 block후 대기
			<-responseModelWithCtx.Ctx.Done()
			responseData = responseModelWithCtx.Value
		} else {
			// 최초의 요청인 경우
			fmt.Printf("서버: 2초동안 [%s] 라는 요청은 response를 받기전까지 타서버에 중복으로 보내지 않습니다 \n", param)
			responseData, err = proxyServiceImpl.dataGateway.GetDataFromDataServer(param)
			if err != nil {
				return "", err
			}
			responseModelWithCtx.Value = responseData
			// 최초의 요청이 서버에서 데이터를 받으면 나머지 요청들과 데이터를 공유하고 나머지 요청의 block을 푼다.
			responseModelWithCtx.Cancel()

			// 클라이언트에게 응답을 먼저 내려준후 데이터 캐시
			go func() {
				proxyServiceImpl.cacheRepository.SetData(param, responseData)
			}()
		}
	}
	return responseData, nil
}
