# OpenFx-cli
OpenFx를 사용하기 위한 Command Line Interface 도구이다. 이를 통해 OpenFx 프레임 워크 위 서비스들을 배포하여 사용할 수 있다. 설치 방법은 다음과 같다. 
### 

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

