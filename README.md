## Socks5 to HTTP

> 这是一个超简单的 Socks5 代理转换成 HTTP 代理的小工具。

### 如何安装？

#### Golang 用户

```sh
# Required Go 1.17+
go install github.com/mritd/s2h@master
```

#### Docker 用户

```sh
docker pull mritd/s2h
```

### 如何使用？

#### 二进制安装用户

```sh
# -l 本地 HTTP 监听地址
# -s 远程 Socks5 服务器地址
s2h -l 127.0.0.1:8081 -s 127.0.0.1:1080
```

#### Docker 用户

```sh
docker run --rm -it --network=host mritd/s2h -l 127.0.0.1:8081 -s 127.0.0.1:1080
```