# go cache example

## 1. 캐시 서버의 기능
- 클라이언트에게 요청이 오면 캐시되었는지 확인한후, 캐시되어있다면 본서버로 요청을 보내지 않고 데이터를 내려줍니다. 
- 동일한 요청이 여러개 들어올 시(본서버에서 응답을 받기전이라 캐시되지 않은경우) 동일한 요청은 본서버에 한번만 보내고 나머지 요청은 블락시키고 응답을 받으면 공유하여 내려줍니다.

## 2. 캐시 저장소 구성 및 저장규칙
- CacheRepository <br>
go cache를 사용하여 관리하여 본서버에서 데이터를 받으면 1분간 데이터를 캐시합니다.
캐시를 저장하느라 클라이언트에게 데이터를 늦게 주면 안되므로 저장로직은 고루틴(스레드)으로 돌아갑니다.
스케일 아웃을 고려한다면 구현체로 global nosql을 사용합니다.

- TemporaryRepository <br>
동일한 요청이 본서버의 응답을 받았다면 데이터가 캐시되어있겠지만, 
본서버에서 응답이 오기전(캐시하기전) 동일한 요청이 또 들어왔을때를 처리하기 위함입니다.
(본서버에서 응답을 내려주는데 1초가 걸린다면 1초안에 온 100개의 동일한 요청은 본서버에 들어갈 필요 x) 
본서버에 들어간 요청과 동일한 요청은 block시킨후 하나의 응답이 오면 모두 response를 내려줍니다.
본서버에서 응답이 오면 고루틴(스레드)를 사용하여 데이터를 캐시하므로 고루틴이 돌기전 들어온 동일한 요청의 데이터를 캐시(2초간)하는 역할도 포함됩니다.


## 3. output
```
프록시서버: [hello] 요청을 본서버에 보냈습니다. 2초동안 동일한 요청은 본서버에 보내지 않습니다 
프록시서버: [hello] 요청은 이미 본서버에 보낸 요청입니다. 응답이 오면 동일한 결과를 내립니다. 
프록시서버: [hello] 요청은 이미 본서버에 보낸 요청입니다. 응답이 오면 동일한 결과를 내립니다. 
본서버: [hello] 요청을 받았습니다 결과값을 내려줍니다. 
클라이언트: 결과 데이터: world
클라이언트: 결과 데이터: world
클라이언트: 결과 데이터: world
서버: 본서버에서 response가 왔으므로 {hello: world}의 데이터를 캐시합니다
프록시서버: [hello]는 캐시되어있으므로 [world]를 내려줍니다
클라이언트: 결과 데이터: world
프록시서버: [hello]는 캐시되어있으므로 [world]를 내려줍니다
클라이언트: 결과 데이터: world
```
