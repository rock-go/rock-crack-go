# 离线密码爆破
基于rock-go的密码爆破模块

## john
- hash 爆破模块
```lua
	local function handle(ev)
	  ev.Put(false , true)
	  rock.ERR("%v" , ev)
	end
        --local john = crack.john("shadow").dict("share/dict/pass.dict").pipe(function(ev) ev.Put(true , true) end)
	local johnn = crack.john{
		name = "shadow",
		dict = "C:\\easypass.txt",
		pipe = handle,
		salt = "123"
	}

	johnn.shadow("root:$6$X7Z9HGT8$.810fPiPxQtCKFH2ecvG/xxtMdzE0pJG.amPTz5W/21/kJQ0O3Wl0:18896:0:99999:7:::")
	johnn.md5("ttttt")
	johnn.sha256("ttttt")
	johnn.sha512("840ee0a9f4deb2ca30714b1b518aee33954a2a468281e9049ef3e3fa23112f5cc0298c396878a2b92a5145094eca605afd195b58c771c5c19a6ff6ec5b738948")
```
### 字段解析

#### `dict` 适用于所有方法
明文的密码文件路径即弱口令文件
#### `salt` 适用于除shadow之外的方法
盐值

### 方法解析

#### `shadow`
破解shadow文件原始字符串，会自动根据不同算法进行爆破
#### `md5`,`sha256`,`sha512`
破解`md5`,`sha256`,`sha512`加密字符串，如果含有salt值的话加密方式为将$string$salt进行组合加密并比较
