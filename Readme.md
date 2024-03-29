2024.2.26 开始毕业论文设计

	暂定:区块链商城项目
	前端端口:xxx:8088
	后端端口:xxx:8080
	

# 代码的实现

## 1.安装Go环境

项目使用 Golang 1.22版本

## 2.下载并配置环境

``` 
下载地址：https://golang.google.cn/dl/
```	

### 2.确定使用Gin框架

```golang
// 下载GIn框架包
go get -u github.com/gin-gonic/gin
```


## 3.确定项目目录


## 4.代码编写（太多已省略）

## 5.linux服务器部署（docker）

### 5.1. 下载并配置docker


### 5.2. 在容器中启动 Mysql

```docker
docker run -d \
  --name mysql-container \
  -e MYSQL_ROOT_PASSWORD=<你的密码> \
  -e MYSQL_DATABASE=<数据库名> \
  -e MYSQL_USER=<用户名> \
  -e MYSQL_PASSWORD=<用户密码> \
  -p 3306:3306 \
  mysql:latest
```

### 5.3. 在容器中启动 Redis

```docker
docker run -d \
  --name redis-container \
  -p 6379:6379 \
  redis:latest
```


## 6.用户商家信息

1. xxxx00 手机品牌专卖店
2. xxxx01 苹果专卖店
3. xxxx02 小米专卖店
4. xxxx03 华为专卖店
5. xxxx04 男装专卖店
6.  xxxx05 女装专卖店
7.  xxxx06 特色专卖店
