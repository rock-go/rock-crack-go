# 密码爆破
基于rock-go的密码爆破模块

## john
- hash 爆破模块
```lua
	local function handle(ev)
	  ev.Put(false , true)
	  rock.ERR("%v" , ev)
	end

	local johnn = crack.john{
		name = "shadow",
		dict = "C:\\easypass.txt",
		pipe = handle,
		salt = "123"
	}

	--johnn.shadow("root:$6$X7Z9HGT8$.810fZP6mWm19PKSboWRLqCjGFyrH5doETlIqfPiPxQtCKFH2ecvG/xxtMdzE0pJG.amPTz5W/21/kJQ0O3Wl0:18896:0:99999:7:::")
	johnn.sha512("840ee0a9f4deb2ca30714b1b518aee33954a2a468281e9049ef3e3fa23112f5cc0298c396878a2b92a5145094eca605afd195b58c771c5c19a6ff6ec5b738948")
    
<!--     local john = crack.john("shadow").dict("share/dict/pass.dict").pipe(function(ev) ev.Put(true , true) end)
    
    john.shadow("$2$aaaaaaa")
    john.md5('xxx')
    john.sha256('xxx')
    john.sha512('xxx222') -->
```
