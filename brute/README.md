## 线上爆破模块
### v1.0
##### 实现并已成功测试的组件包括：
`elastic`
`ftp`
`mongodb`
`mysql`
`rdp`
`redis`
`smtp`
`ssh`
`公司使用vpn`
##### 已经完成但是还没有经过测试的组件包括：
`mssql`
`oracle`
`postgres`
`smb`
`snmp`

## 使用手册

### `example`
```lua
local function handle(ev)
    ev.Put(true , true)
end

local online = crack.online{
    name    = "crackonline",
    iplist  = "127.0.0.1",
    dict    = "share/third/A.txt",
    user    = "root",
    proxy   = "http://127.0.0.1:10809",
    threads = 3,
    timeout = 3,

    pipe = handle
}
```
### 字段含义
#### `iplist`: 不适用于`VPN`爆破模块
ip地址，可填ip文件存放路径或者单独ip地址
#### `dict`: 使用于所有组件
爆破密码，可填字典文件存放路径或者单独密码
#### `user`: 不适用于`redis`组件
用户名，可填用户名文件存放路径或者单独用户
#### `proxy`: 仅适用于`vpn`爆破组件(内网无法访问`vpn`登录页面)
代理地址，用于`vpn`爆破
#### `threads`:适用于所有组件
单个ip爆破时使用的线程数，目的时为了防止ip被封禁
>ftp爆破的时候一个ip连接不能超过3个
>smtp爆破的时候一个线程数量不能过高建议10个
#### `timeout`:适用于所有组件
tcp链接的超时时间

### 使用代码

```lua
online.ssh(22)

--online.mysql(3306)

--online.mongodb(27017)

--online.elastic(9200)

--online.rdp(3389)

--不需要用户名
--online.redis(6379)

--ftp爆破的时候一个ip连接不能超过3个
--online.ftp(21)

--由于外网访问限制，所以可能需要代理地址
--online.vpn(1)

--smtp爆破的时候一个线程数量不能过高建议10个，否则有问题
--email.eastmoney.com的ip地址：[222.73.55.214]
--online.smtp(25)

--online.smb(445) --未测试
--online.snmp(161) --未测试
```
### 使用介绍
根据调用的组件和传入的端口进行线上爆破
爆破成功会触发相应的事件日志，后续增加告警等措施
