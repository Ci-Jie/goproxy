# GoProxy
GoProxy is a module proxy server for golang.
It provide a datastore to store your private go packages from gitlab.

## Getting Started

First, you need to copy the configuration file from `etc/goproxy.example.yaml` and name it as `etc/goproxy.yaml`

```
$ cp etc/goproxy.example.yaml etc/goproxy.yaml
```

Then, you could modify it by yourself according to your environment.
It maybe like below:

```
$ vim etc/goproxy.yaml

port: 8078                             # which port you want to use for goproxy
gitlabs:
- domain: "test.proxy.org"             # your domain name
  endpoint: http://127.0.0.1:30000     # your gitlab endpoint
  token: <Token>                       # your gitlab token which have permission to access gitlab api
storage:
  provider: local                      # what kind of backend storage you want to use. supported: local, s3(in the future)
  local:
    path: tmp/                         # what dirctory you want to store your go packages
```

Finally, start goproxy service. 

```
$ goproxy -c etc/goproxy.yaml start
```

> If you don't provide `-c` flag. it will read `/etc/goproxy/goproxy.yaml` path by default.

## Installation

### Binary

Prerequisites

* Golang `>=16`

You just need to use `make` command to build go proxy binary

```
$ make
```

If you want to build binary and container image at one time.
You can use `make all` command to build them.

```
$ make all
```

### Docker

You can use `make build-image` commnad to build container image.
It will build binary based on `golang:1.17.3` image and copy it on the other container which is baes on `alpine:latest`
And running go proxy service on it.

```
$ make build-image
```

More detail information for can click [here](https://github.com/Ci-Jie/goproxy/blob/master/Dockerfile)