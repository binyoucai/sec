# 腾讯云函数下 c2 隐蔽
## 使用方法
### 1、修改c2服务端配置文件
   本项目下tecent_cloud_func.profile
    ```
    ./teamserver ip passowrd win_tecent_cloud_func.profile
    ```
### 2、修改配置文件c2地址
    修改config.yaml中配置文件 C2Srv下 Address 为c2的回源地址
    
### 3、编译上传至腾讯云函数
```
GOOS=linux GOARCH=amd64 go build -o main main.go
zip main.zip main config.yaml
```    




