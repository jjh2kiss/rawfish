# rawfish
[![Build Status](https://travis-ci.org/jjh2kiss/rawfish.png?branch=master)](https://travis-ci.org/jjh2kiss/rawfish)  

HTTP Server for RAW HTTP Message

RAWFISH는 HTTP Message 자체를 서비스하는 경량 웹 서버 입니다.

### 주요 기능
* BandWidth 조절
* Read/Write Timtout 설정
* HTTP/HTTPS 지원 
* 강제 200 OK 지정(요청한 주소가 존재하지 않을 경우 200 OK 전송)
* RAW/NORMAL 의 선택적 적용 및 상속 

### Download
> git clone https://github.com/jjh2kiss/rawfish.git

### Build & Install
rawfish를 Golang를 이용해 개발되었습니다. 소스를 이용해 빌드하려면 Golang(>=1.7.0)가 설치되어 있어햐 합니다.
Golang의 설치를 아래 주소를 참조하세요.
https://golang.org/doc/install#install

golang가 설치되어 있다면 아래 명령을 실행하여 rawfish를 빌드할 수 있습니다.
> cd [RAWFISHDIR]  
> go get  
> go build  
> go install  

### 서버 구동
Rawfish는 http, https를 모두 서비스합니다.
아래 명령을 이용해 현재 디렉토리의 파일을 서비스 하는 rawfish server를 구동 할 수 있습니다.
> rawfish -r ./

wget를 이용해 rawfish서버가 정상적으로 구동되었는지 확인할 수 있습니다.
> wget http://localhost

### RAW 모드와 Normal 모드
rawfish는 normal, raw 두가지 모드를 지원합니다.  
normal은 apache, nginx의 Static File 서비스와 동일합니다. 사용자가 요청한 파일이 존재할 경우 HTTP 프로토코을 이용해 파일을 전송합니다.  
raw 모드의 경우 HTTP Message 서비스합니다. HTTP Message는 보통 아래와 같이 구성됩니다.  

HTTP Message | Example
------------ | ------
STATUS LINE  | HTTP/1.1 200 OK
HEADERS(General, Reponse, Entity) | Connection: Close
Message Body | &lt;HTML&gt;...&lt;/HTML&gt;

RESTFul API, Header 값에 따른  클라이언트 동작 테스트 등에 유용하게 사용하 수 있습니다.  

Mode 판정은 아래 규직을 따릅니다.  

1. 요청한 파일의 디렉토리에 ".raw"파일이 존재하는지 검사한다.  
2. ".raw"파일이 존재할 경우 RAWMODE로 서비스한다.  
3. ".raw"파일이 존재하지 않을 경우 ".normal"파일을 검사한다.  
4. ".normal"파일이 존재할 경우 NORMAL MODE로 서비스한다.  
5. ".raw", ".normal"파일이 모두 존재하지 않을 경우 상위 디렉토리 이동후 1번 과정을 수행한다.  
6. 최상위 디렉토리에 ".raw", ".normal"모두 없을 경우 NORMAL 모드로 서비스한다.  

모드는 상속 가능한 속성입니다.  
최상위 디렉토리에 ".raw"를 생성해 놓을 경우, 하위디렉토는 자신의 모드를 지정하지 않는이상 RAWMODE를 사용하게 됩니다.  

rawfish는 아래와 같은 모드 예제를 포함하고 있습니다. 사용에 참고하세요  
[RAWFISHDIR]/samples  
> .  
> ./normal  
> ./normal/www.youtube.com  
> ./raw  
> ./raw/www.youtube.com  
> ./raw/.raw  


1. http://127.0.0.1 -> Directory Listing
2. http://127.0.0.1/normal/www.youtube.com -> Normal 모드, www.youtube.com 파일 다운로드
3. http://127.0.0.1/rawfish/www.youtube.com -> Raw 모드, http://www.youtube.com으로 이동 

### Options
* --addr, -a
서비스에 사용할 주소를 지정합니다. 기본값은 "0.0.0.0"으로 모든 주소에 바인딩됩니다.
* --port, -p
서비스에 사용할 TCP Port를 지정합니다. 기본값은 "80"입니다.
* --root, -r
서비스할 디렉토리를 지정합니다.
* --read-timeout
데이터를 클라이언트로 부터 읽을때 사용할 시간제한 값입니다. 기본값은 10초 입니다.
* --write-timeout
데이터를 클라이언트에 전송할 때 사용할 시간제한 값입니다. 기본값은 10초 입니다.
* --force-200-ok, -f
지정한 파일이 없거나, 기타 다른 이류로 서비스가 불가능 할 경우 5XX, 4XX 에러 대신 200 OK를 강제로 전송합니다.
* --force-200-ok-content-size
force-200-ok 옵션에 의해 강제로 200 OK를 전송할때 사용할 컨텐츠의 사이즈 입니다.
바이트 단위로 지정할 수 있습니다. 기본값은 0입니다. 
* --https
HTTPS 모드를 활성화 합니다. HTTPS 모드를 사용하기 위해서는 반드시 --pemfile을 이용해 pemfile을 지정해야 합니다. 또한 --port=443으로 지정할 것을 권장합니다.
* verbose
HTTP Request, Response에 대한 정보를 확인할 수 있습니다.
* process
서비스에 사용할 물리적인 CPU 코어개수를 지정합니다. 기본값은 1입니다.
* rate
Bandwitdh를 지정합니다. Byte단위로 지정할 수 있습니다. 기본값은 0이며 무제한을 나타냅니다.


