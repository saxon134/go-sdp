# go-sdp

### 介绍

Service Discovery Protocol

基于Redis的服务发现、负载均衡

### 使用须知

Redis必须配置

如果配置了秘钥，则所有接口都需要加密

加密方式是：sign = lowercase( MD5(secret + timestamp) )，接口带上sign和timestamp


### 待开发

如果配置了MySQL：

	1. 可以作为配置中心使用

	2. 获取注册服务所在服务器的负载情况


如果配置了阿里云账户：

	1. 可以通过接口获取负载情况
	

配置中心、查看负载情况，参考前端项目：https://github.com/saxon134/go-sdp-ms


### 配置文件说明（config.yaml）

```
name: sdp    //项目名，默认sdp。Redis key会添加其作为前缀
path: /      //http服务接口路径
port: 

redis:       //不可空
  host: 
  pass:

```



### 安装教程

```
go get github.com/saxon134/go-sdp;
go mod tidy;
go run main.go;
```




