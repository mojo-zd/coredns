### 下载
- download
```
> wget -O coredns.zip http://files.git.oschina.net/group1/M00/04/96/PaAvDFtxDoiAT5_rAM_rm0oYy4I638.zip?token=551d67b379111f684bbe9a5df219b79f&ts=1534135954&attname=coredns_1.2.0.1_linux_amd64.zip
```

- unzip
```
> unzip coredns.zip
```
> 如果zip工具未安装 执行`yum install -y unzip zip`

- move coredns file
```
> mkdir /etc/coredns

> cp Corefile z_health.org /etc/coredns

> mv coredns /usr/bin/

> chmod +x /usr/bin/coredns
```

- 配置systemd

```
> sudo vi /etc/systemd/system/coredns.service

[Unit]
Description=coredns: The dns server
Documentation=https://github.com/coredns/coredns

[Service]
ExecStart=/usr/bin/coredns
Restart=always
StartLimitInterval=0
RestartSec=10

[Install]
WantedBy=multi-user.target
```

- start service
```
> chmod +x /etc/systemd/system/coredns.service

> systemctl start coredns

> systemctl status coredns
```

### Coredns plugin z_health_check setting
- Corefile定义
```
.:53 {
    loadbalance round_robin
    z_health_check /etc/coredns/z_health.org
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

`www.registry.com`  域名

`ips` backend services的host ip

`api` server任意可调用的get请求的api

`protocol` 不填写默认为 "http://"

### Test
- before test
```
vi /etc/resolv.conf

# Generated by NetworkManager
# nameserver 10.0.0.1
nameserver 10.0.0.131  //当前dns service ip地址
```

- dig `dig @localhost www.registry.com`
- ping `ping www.registry.com`  反复`ping`将会返回不同的service ip
