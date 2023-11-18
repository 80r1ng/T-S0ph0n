## T-S0ph0n

​	用途:特殊Windows环境下进行用户/组添加操作



### 安装:

```shell
go mod tidy
go build -o tsophon.exe
.\tsophon.exe -h(获取帮助)
.\tsophon.exe add --t 1 -u xxx -p xxx (使用方式一进行用户添加)
.\tsophon.exe add --t 1 -u xxx -p xxx -g xxx (使用方式一进行用户添加,并将用户添加至指定组)
```





### Usage:

```yaml
Usage:
  tsophon [command]

Available Commands:
  add         绕过反病毒软件添加指定用户

Flags:
  --t uint8        指定绕过方式
  -u, --user string    指定用户名
  -p, --pwd string     指定密码
  -g, --group string   指定组
```



绕过方式:

​	方式一:调用WindowsAPI进行绕过(--t 1)

​	方式二:复制net1.exe文件至其它目录进行执行,这里默认是D盘(--t 2)





参考视频:

​	https://www.bilibili.com/video/BV1HR4y1q7aL/?spm_id_from=333.880.my_history.page.click&vd_source=7daf77b55473abbb4599de8e3ea76a1b