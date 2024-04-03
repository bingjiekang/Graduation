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

### 5.4 在容器中启动vue

```docker
# 使用 Node 作为基础镜像
FROM node:alpine as build-stage

# 设置工作目录
WORKDIR /app

# 将 package.json 和 package-lock.json 复制到工作目录
COPY package*.json ./

# 安装项目依赖
RUN npm install

# 将项目文件复制到工作目录
COPY . .

# 构建 Vue 项目
RUN npm run build

# 使用 Nginx 作为另一个基础镜像
FROM nginx:alpine

# 复制构建好的 Vue 项目到 Nginx 的默认站点目录
COPY --from=build-stage /app/dist /usr/share/nginx/html

# 暴露 Nginx 默认端口
EXPOSE 80

# 启动 Nginx
CMD ["nginx", "-g", "daemon off;"]

```

```docker
docker build -t vue-app .
docker run -d -p 8080:80 vue-app
```


## 6.用户商家信息

1. xxxx00 手机品牌专卖店
2. xxxx01 苹果专卖店
3. xxxx02 小米专卖店
4. xxxx03 华为专卖店
5. xxxx04 男装专卖店
6.  xxxx05 女装专卖店
7.  xxxx06 特色专卖店
