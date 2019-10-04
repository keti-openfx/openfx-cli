# OpenFx-cli
OpenFx를 사용하기 위한 Command Line Interface 도구이다. 이를 통해 OpenFx 프레임 워크 위 서비스들을 배포하여 사용할 수 있다. 설치 방법은 다음과 같다. 
### 



# Requirements

`openfx-cli`를 사용하여 서비스들을 배포하기 위해서는 미니쿠베를 통해 구동된 쿠버네티스 클러스터 내에 `openfx-gateway`가 컨테이너로 실행 중이어야 한다. 이는 다음의 [링크](<https://github.com/keti-openfx/openfx/blob/master/documents/3.Compile_OpenFx.md>)를 통해 진행할 수 있다. 



# Compile OpenFx-cli

`openfx-cli`를 클론하여 컴파일을 진행한다. 

```
$ git clone https://github.com/keti-openfx/openfx-cli.git
$ cd openfx-cli
```

`make`명령을 실행하여 컴파일을 진행한다.

```
$ make build
```

`$GOPATH/bin`을 확인해보면 `openfx-cli`가 컴파일 되어있는 것을 확인할 수 있다.

```
$ cd $GOPATH/bin
$ ls
openfx-cli
```

