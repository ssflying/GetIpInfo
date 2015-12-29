# GetIpInfo

调用ip.taobao.com返回ip的地理和运营商信息

Usage:

```bash
cat iplist | ./GetIpInfo -p <qps>
./GetIpInfo -p <qps> < iplist
```

qps: 默认10,taobao的限速
