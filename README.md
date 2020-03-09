# OpenFx-cli
OpenFx를 사용하기 위한 Command Line Interface 도구이다. 이를 통해 OpenFx 프레임 워크 위 서비스들을 배포하여 사용할 수 있다. 설치 방법은 다음과 같다. 
### 



# Requirements

### Compile OpenFx

`openfx-cli`를 사용하여 서비스들을 배포하기 위해서는 `openfx-gateway`와 `openfx-executor`의 컴파일이 완료되어야 하고, 미니쿠베를 통해 구동된 쿠버네티스 클러스터 내에 `openfx-gateway`가 컨테이너로 실행 중이어야 한다. 이는 다음의 [링크](<https://github.com/keti-openfx/openfx/blob/master/documents/3.Compile_OpenFx.md>)를 통해 진행할 수 있다. 

###

### Installing Go

- `openfx-cli`를 사용하기 위해서는 Go 언어가 설치되어 있어야 한다.  Go 언어 설치를 진행하기 위해 [공식 홈페이지](<https://golang.org/doc/install>)에서 **호스트 OS에 맞게** 원하는 버전의 설치파일을 다운로드 받는다. 아래와 같은 명령어를 이용하여 설치파일을 압축해제하고, 압축 해제된 디렉토리를 `/usr/local`로 위치를 옮긴다. 

  ```
  $ wget https://dl.google.com/go/go[version].[Host OS].tar.gz
  $ sudo tar -xvf go[version].[Host OS].tar.gz
  $ sudo mv go /usr/local
  ```

- `.bashrc` 파일을 수정하여 go와 관련된 환경변수를 설정한다.

  ```
  $ vim ~/.bashrc
  >>
  # add this lines
  export GOROOT=/usr/local/go
  export GOPATH=$HOME/go
  export PATH=$GOPATH/bin:$GOROOT/bin:$PATH
  ```

- 변경한 `.bashrc` 파일을 적용한다.

  ```
  $ source ~/.bashrc
  ```

- Go가 설치되었는지를 확인한다.

  ```
  $ go version
  >>>
  go version go1.12.3 linux/amd64
  
  $ go env
  >>>
  GOARCH="amd64"
  GOBIN="/root/workspace/go/bin"
  GOCACHE="/root/.cache/go-build"
  GOEXE=""
  GOFLAGS=""
  GOHOSTARCH="amd64"
  GOHOSTOS="linux"
  GOOS="linux"
  GOPATH="/root/workspace/go"
  GOPROXY=""
  GORACE=""
  GOROOT="/usr/local/go"
  GOTMPDIR=""
  GOTOOLDIR="/usr/local/go/pkg/tool/linux_amd64"
  GCCGO="gccgo"
  CC="gcc"
  CXX="g++"
  CGO_ENABLED="1"
  GOMOD=""
  CGO_CFLAGS="-g -O2"
  CGO_CPPFLAGS=""
  CGO_CXXFLAGS="-g -O2"
  CGO_FFLAGS="-g -O2"
  CGO_LDFLAGS="-g -O2"
  PKG_CONFIG="pkg-config"
  GOGCCFLAGS="-fPIC -m64 -pthread -fmessage-length=0 -fdebug-prefix-map=/tmp/go-build323300119=/tmp/go-build -gno-record-gcc-switches"
  ```

  - `go version`과 `go env`를 입력하여 출력되는 정보는 위와 상이할 수 있다.



### Installing Docker

`openfx-cli`를 통해 빌드된 함수 이미지를 레지스트리에 저장하기 위해서는 해당 레지스트리에 로그인하는 과정이 필요하다. 이를 위해 먼저 도커를 설치해주어야 한다. 

```
$ sudo apt-get remove docker docker-engine docker.io containerd runc
$ sudo apt-get update

$ sudo apt-get install apt-transport-https \
    ca-certificates curl gnupg-agent software-properties-common
$ curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo apt-key add -
$ sudo apt-key fingerprint 0EBFCD88
$ sudo add-apt-repository \
   "deb [arch=amd64] https://download.docker.com/linux/ubuntu \
   $(lsb_release -cs) stable"

$ sudo apt-get update
$ sudo apt-get install docker-ce docker-ce-cli containerd.io
$ docker --version
>>
Docker version 18.06.1-ce, build e68fc7a

$ sudo usermod -aG docker $USER 
$ reboot

$ sudo systemctl status docker
```

- `docker version`을 입력하여 출력되는 정보는 위와 상이할 수 있다.



### Setting insecure registries

도커 레지스트리는 SSL 인증서 없이 이용할 수 없다. SSL 인증서 없이 도커 레지스트리를 사용하기 위해서는 `insecure registries`에 대한 설정이 필요한데, 이는 다음과 같이 진행한다.

```
$ sudo vim /etc/docker/daemon.json

>>
{"insecure-registries": ["YOUR PRIVATE REGISTRY SERVER IP:PORT"]}

$ service docker restart
```

- YOUR PRIVATE REGISTRY SERVER IP:PORT에 사용하고자하는 도커 레지스트리의 IP 주소와 Port 번호를 기입하면 된다. 



# Compile OpenFx-cli

- `openfx-cli`를 클론하여 컴파일을 진행한다. `openfx-cli` 클론은 **keti-openfx** 디렉토리 밑에서 진행한다.

  ```
  $ cd $GOPATH/src/github.com/keti-openfx
  $ git clone https://github.com/keti-openfx/openfx-cli.git
  $ cd openfx-cli
  ```

- `make`명령을 실행하여 컴파일을 진행한다.

  ```
  $ make build
  ```

- `$GOPATH/bin`을 확인해보면 `openfx-cli`가 컴파일 되어있는 것을 확인할 수 있다.

  ```
  $ cd $GOPATH/bin
  $ ls
  openfx-cli
  ```



# Verify OpenFx-cli

`openfx-cli` 컴파일까지 완료하였으면, 프레임워크 위 API 단위 응용이 배포되는지를 확인하여야 한다. 이는 다음과 같은 절차로 진행한다. 

## Create folder for CLI testing

```
$ mkdir cli-test
$ cd cli-test
```



## Cloninig `OpenFx-runtime`

```
$ git clone https://github.com/keti-openfx/OpenFx-runtime.git runtime
```



## Create OpenFx function

- 함수를 배포하기 위해 함수의 initialization을 진행(runtime 설정, 함수 이름 설정 및 config.yaml 파일 생성)

  ```bash
  $ openfx-cli function init <FUNCTION NAME> --runtime <RUNTIME NAME> 
  >> 
  Folder: <FUNCTION NAME> created
  Fucntion handler created in folder: <FUNCTION NAME>/src
  Rewrite the function handler code in <FUNCTION NAME>/src folder
  Config file written: config.yaml
  
  $ cd <FUNCTION NAME>
  ```

  > `<FUNCTION NAME>` : 생성하고자할 함수의 이름
  
  > `<RUNTIME NAME>` : 함수 작성 언어 (go, python2, python3, nodejs, ruby, cpp, java, csharp 중 택일)

- 게이트웨이 설정

  함수 초기화 시, 기본적으로 설정되는 게이트웨이의 주소는 `localhost:31113`이다.  **다른 사용자 호스트 OS에 함수를 배포하고자 하는 경우**, `--gateway` 옵션으로 게이트웨이 주소를 변경할 수 있다. 

  ```bash
  $ openfx-cli function init <FUNCTION NAME> --runtime <RUNTIME NAME> --gateway <호스트 OS IP:31113>
  
  >>
  Folder: <FUNCTION NAME> created
  Fucntion handler created in folder: <FUNCTION NAME>/src
  Rewrite the function handler code in <FUNCTION NAME>/src folder
  Config file written: config.yaml
  
  $ cd <FUNCTION NAME>
  ```

- runtime은 `go`, `python2`, `python3`, `nodejs`, `ruby`, `cpp`, `java`, `csharp`을 지원한다.



## Configure `config.yaml`

```
functions:
  <FUNCTION NAME>:
    runtime: <RUNTIME NAME>
    desc: ""
    maintainer: ""
    handler:
      dir: ./src
      file: handler.<RUNTIME>
    docker_registry: <REGISTRY IP>:<PORT>
    image: <REGISTRY IP>:<PORT>/<FUNCTION NAME>
    requests:
      memory: 50Mi
      cpu: 50m
      gpu: ""
openfx:
  gateway: <호스트 OS IP>:31113
```

- `<REGISTRY IP>`, `<PORT>`를 레지스트리에 맞춰 변경한다.
- `gateway`의 <호스트 OS IP>는 `function init` 시 지정한 IP 이다. 
- `requests`는 사용자가 정의할 서비스 별 자원 사용량이며, 각각의 항목은 다음과 같다.
  - memory: 서비스 별 memory 사용량, 최대 200Mi 까지 지정할 수 있으며, 기본 값은 50Mi 이다.
  - cpu: 서비스 별 cpu 사용량, 최대 80m까지 지정할 수 있으며, 기본 값은 50m 이다. 

## Writing Handler

- Handler 코드 작성(함수 init 시 지정한 runtime 선택)

  - Golang

    handler.go

    ```go
    package main
    
    import sdk "github.com/keti-openfx/openfx/executor/go/pb"
    
    func Handler(req sdk.Request) string {
        return string(req.Input)
    }
    ```

  - Python 2.7 / 3.4

    handler.py

    ```python
    def Handler(req):
        return req.input
    ```

    라이브러리 추가시,  `requirements.txt` 에 필요 라이브러리를 명시해야 한다.

    

    다음은 현재 시간을 출력하는 예제이다. 
  
    requirements.txt
  
    ```
    datetime
    ```

     ```python
  import datetime 
    
  def Handler(req):
        return datetime.datetime.now()
     ```
  
  - Node Js
  
    handler.js
  
    ```js
    function Handler(argStr) {
      return argStr;
    }
  
    module.exports = Handler;
  ```
  
  - Ruby
  
    handler.rb
  
    ```ruby
    #!/usr/bin/env ruby
    
    module FxWatcher
    def FxWatcher.Handler(argStr)
        return argStr
    end
    end
  ```
  
  - C++
  
    handler.cc
  
    ```c++
    #include <iostream>
    
    using namespace std;
    
    string Handler(const string req) {
      return req;
    }
  ```
  
- Java
  
  Handler.java
  
    ```java
    package io.grpc.fxwatcher;
    
    import com.google.protobuf.ByteString;
    
    public class Handler {
    
      public static String reply(ByteString input) {
        return input.toStringUtf8() + "test";
      }
    
    }
    ```
  
  - C#
  
    handler.cs
  
    ```c#
    namespace Fx
    {
        class Function
        {
            public byte[] Handler(byte[] Input)
            {
                return Input; 
            }
        }
    }
    ```
  
    라이브러리 추가시, `fxServer.csproj` 에 필요 라이브러리를 명시애햐한다. 
  
    ```xml
    <Project Sdk="Microsoft.NET.Sdk">
    
      <PropertyGroup>
        <OutputType>Exe</OutputType>
        <TargetFrameworks>netcoreapp2.1</TargetFrameworks>
      </PropertyGroup>
      <!-- openfx default installation library. Never modify-->
      <ItemGroup>
        <PackageReference Include="Google.Protobuf" Version="3.7.0" />
        <PackageReference Include="Grpc" Version="1.20.1" />
        <PackageReference Include="Grpc.Tools" Version="1.20.1" />
      </ItemGroup>
    
      <ItemGroup>
        <Protobuf Include="fxwatcher.proto" Link="fxwatcher.proto"/>
      </ItemGroup>
    
      <!-- Input Server Library.-->
    
    </Project>                            
    ```
  
    * Openfx Handler에는 grpc 관련 라이브러리를 설치해야 하므로 위와 같이 명시하였다.  라이브러리 명시시,  `<PackageReference Include="Grpc.Tools" Version="1.20.1" />` 처럼 포매팅이 필요하다.
  
      

## Building Function

- Kubernetes에 생성한 함수를 배포하기 위한 도커 이미지 생성. (도커 이미지는 로컬에 생성됨)

  ```bash
    $ openfx-cli function build
    
    >> 
    Building function (<FUNCTION NAME>) image ...
    Image: <REGISTRY IP>:<PORT>/<FUNCTION NAME> built in local environment.
  ```

- `-v` 옵션으로 이미지가 빌드되는 과정을 로그로 확인할 수 있다.

  ```bash
    $ openfx-cli function build -v
    
    >>
    ...
    Building function (<FUNCTION NAME>) image ...
    Image: <REGISTRY IP>:<PORT>/<FUNCTION NAME> built in local environment.
  ```


## Deploying Funtion

- 생성된 이미지를 통해 Kubernetes에 함수 배포.

  ```bash
    $ openfx-cli function deploy -f config.yaml -v
  
    >> 
    Pushing: echo, Image: <REGISTRY IP>:<PORT>/<FUNCTION NAME> in Registry: <REGISTRY IP>:<PORT>...
    ...
    Deploying: echo ...
    Attempting update... but Function Not Found. Deploying Function...
    http trigger url: http://localhost:31113/function/echo
  ```

- OpenFx는 기본적으로 사용자 함수에 대해 오토스케일링을 제공하고 있다. 사용자 함수의 자원 사용량이 OpenFx API Gateway에서 한계치로 지정한 자원 사용량의 수치를 넘어서게 되면, 사용자 함수의 레플리카 셋을 만든다. 함수 deploy 시, 사용자는 다음과 같이 `--min`, `--max` 옵션을 통해 사용자 함수에 대한 레플리카 셋의 최소값과 최대값을 지정할 수 있다. 

  ```bash
    $ openfx-cli function deploy -f config.yaml --min 2 --max 4 -v
  
    >> 
    Pushing: echo, Image: <REGISTRY IP>:<PORT>/<FUNCTION NAME> in Registry: <REGISTRY IP>:<PORT>...
    ...
    Deploying: echo ...
    Attempting update... but Function Not Found. Deploying Function...
    http trigger url: http://localhost:31113/function/echo
  ```

  - `default min : 1`, `default max : 1`

- 함수 Initialization 시 `--gateway` 옵션으로 게이트웨이를 설정하였다면 함수 배포 시, 마찬가지로 `--gateway` 옵션을 주어야 한다. 

  ```bash
    $ openfx-cli function deploy -f config.yaml --min 2 --max 5 --gateway <호스트 OS IP:31113> -v
  
    >> 
    Pushing: echo, Image: <REGISTRY IP>:<PORT>/<FUNCTION NAME> in Registry: <REGISTRY IP>:<PORT>...
    ...
    Deploying: echo ...
    Attempting update... but Function Not Found. Deploying Function...
    http trigger url: http://<호스트 OS IP:31113>/function/echo
  ```

## Confirm OpenFx function list

- Kubernetes에 배포가 완료된 함수의 목록 확인.

  ```bash
    $ openfx-cli function list
  
    >> 
    Function    Image           Maintainer    Invocations    Replicas    Status    Description
    echo        $(repo)/echo                  0              1           Ready
  ```

- 게이트웨이 옵션

  ```bash
    $ openfx-cli function list --gateway <호스트 OS IP:31113>
  
    >> 
    Function    Image           Maintainer    Invocations    Replicas    Status    Description
    echo        $(repo)/echo                  0              1           Ready
  ```


## Verify deployed function using invoke

- Kubernetes에 배포된 함수를 호출.

  ```bash
    $ echo "Hello" | openfx-cli function call echo
  
    >> 
    Hello
  ```

- 게이트웨이 옵션

  ```bash
    $ echo "Hello" | openfx-cli function call echo --gateway <호스트 OS IP:31113>
  
    >> 
    Hello
  ```

## Rolling update Function

- Handler 코드 수정 후, [Building Function](#Building Function)부터 순차적으로 진행하면 된다. 
