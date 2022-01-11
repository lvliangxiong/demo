## Requisite

A rabbitmq docker container should be running on your machine with an address of `127.0.0.1:5672`. Instead, specify your rabbitmq instance addr in the `conf.local.yaml` file. 

## Publisher

每隔一段时间就发送一条日志消息。

## Subscriber

能够接收到自**自身运行**以来 publisher 发送的所有消息。通过不同的启动方式，可以将日志记录在不同的地方。

1. 记录到文件

```bash
go run subcriber/subscriber.go 2>logs
```

2. 打印在 console 中

```bash
go run subcriber/subscriber.go
```

