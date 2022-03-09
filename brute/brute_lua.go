package brute

import (
	"github.com/rock-go/rock/cidr"
	"github.com/rock-go/rock/lua"
)

func (b *brute) sshL(L *lua.LState) int {
	b.append(newBruteSsh(L))
	return 0
}

func (b *brute) ftpL(L *lua.LState) int {
	b.append(newBruteFtp(L))
	return 0
}

func (b *brute) elasticL(L *lua.LState) int {
	b.append(newBruteElastic(L))
	return 0
}

func (b *brute) useL(L *lua.LState) int {
	//读取堆栈
	//写入
	//b.config.dict.use
	return 0
}

func (b *brute) useFL(L *lua.LState) int {
	//读取路径
	//赋值
	//条件: 如果 user.txt 文件大小 小于 dict.buffer  将全部缓存 user []string
	return 0
}

func (b *brute) passL(L *lua.LState) int {
	//读取堆栈
	//写入
	//b.config.dict.pass
	return 0
}

func (b *brute) passFL(L *lua.LState) int {
	//读取堆栈
	//写入
	//b.config.dict.pass
	//条件: 如果 pass.txt 文件大小 小于 dict.buffer  将全部缓存 pass []string
	return 0
}

func (b *brute) cidrL(L *lua.LState) int {
	b.cfg.cidr = cidr.Check(L)
	return 0
}

func (b *brute) hostL(L *lua.LState) int {
	return 0
}

func (b *brute) Index(L *lua.LState, key string) lua.LValue {
	switch key {
	case "use":
		return lua.NewFunction(b.useL)

	case "pass":
		return lua.NewFunction(b.passL)

	case "Use":
		return lua.NewFunction(b.useFL)

	case "Pass":
		return lua.NewFunction(b.passFL)

	case "cidr":
		return lua.NewFunction(b.cidrL)

	case "host":
		return lua.NewFunction(b.hostL)

	case "ssh":
		return lua.NewFunction(b.sshL)

	case "elastic":
		return lua.NewFunction(b.elasticL)

	case "ftp":
		return lua.NewFunction(b.ftpL)

	default:
		//其他协议

	}

	return nil
}
