# reverse proxy

reverse proxy 是一个命令行程序。

它含有三个基本参数，分别是：

- listen/l：该参数指定软件启动后监听的IP与端口，默认值是 "127.0.0.1:8081"
- proxy/p: 该参数为流量转发时配置的代理服务。
- target/t: 该参数指定软件启动后转发的目标，该参数必须提供，示例："https://example.com/"

示例1：
command:

```shell
rp -l "127.0.0.1:3001" -p "http://127.0.0.1:3080" -t "https://example.com/"
```

预期行为：

软件将监听 127.0.0.1:3001 , 当用户发送请求 `http://127.0.0.1:3001/path/to/server?pa=1` 后，
系统将使用http代理`http://127.0.0.1:3080`请求 `https://example.com/path/to/server?pa=1`并返回
