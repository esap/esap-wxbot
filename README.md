# esap-wxbot
Folk from ManiacMike/go-wxbot

根目录下创建一个config.ini，就可以跑了，mac/linux/windows都支持

config.ini的内容如下:

```
[esap]
# 远程esap服务器API
remote = http://192.168.99.10:9090/robot/

# 本地ESAP回调地址，用于扫码登陆等，ESAP服务器要能访问到
local = 192.168.99.10

# 本地服务端口
port = 19090
```
