## How to test
> make test

## How to run
> make run
> Open modern browser & enter "http://localhost:9000"


## 진행 과정
* Front-end
    * jquery는 너무 무겁고 대부분의 method 사용하지 않으므로 jquery와 인터페이스가 거의 같고 헐씬 크기가 작은 zepto.js 사용
    * UI Element를 Component 단위로 나누어 javascript class로 매핑시켜 작은 작업으로 나눔
* Back-end
    * 언어는 Golang
    * web server 구현은 기본 라이브러리인 net/http 사용
    * 크게 handler, service, data access layer로 나누어 작업함
    * todo에 child/parent 계층 관계가 있으므로 todo model에 parent id를 저장하는 property 추가
    * 참조가 걸린 todo를 완료처리하기 위해 Service layer의 "checkIfParentsAreComplete" 함수를 재귀적으로 호출하여 직관적으로 해결함.
