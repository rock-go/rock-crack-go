# 密码爆破
基于rock-go的密码爆破模块

## john
- hash 爆破模块
```lua
    local john = crack.john("shadow").dict("share/dict/pass.dict").pipe(function(ev) ev.Put(true , true) end)
    
    john.shadow("$2$aaaaaaa")
    john.md5('xxx')
    john.sha256('xxx')
    john.sha512('xxx222')
```