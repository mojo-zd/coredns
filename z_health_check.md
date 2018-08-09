### coredns z_health_check使用
- Corefile定义
```
.:53 {
    loadbalance round_robin
    z_health_check z_health.org
    log
}
```

1. `loadbalance`为lb插件 参考 https://coredns.io/plugins/loadbalance

2. `z_health_check`为server健康检查插件
`z_health.org`为指定health check文件名称. 文件内容示例如下:
```
{
  "www.registry.com": {
    "ips": ["139.159.228.163", "139.159.225.130", "118.31.50.65"],
    "api": "i18n/lang/en-us-lang.json",
    "protocol":"http://"
  }
}
```

说明:

`www.registry.com`为设定的域名

`ips` 为后端运行的host ip

`api` server任意可调用的get请求的api

`protocol` 不填写默认为 "http://"

