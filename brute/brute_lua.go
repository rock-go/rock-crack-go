package brute

import (
	"fmt"
	"github.com/rock-go/rock/cidr"
	"github.com/rock-go/rock/lua"
	"os"
)

func (b *brute) Index(L *lua.LState, key string) lua.LValue {
	println("lll:", L.GetTop())
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

	case "start":
		return lua.NewFunction(b.startL)

	case "ssh":
		return lua.NewFunction(b.sshL)

	case "elastic":
		return lua.NewFunction(b.elasticL)

	case "ftp":
		return lua.NewFunction(b.ftpL)

	case "mongodb":
		return lua.NewFunction(b.mongodbL)

	case "mssql":
		return lua.NewFunction(b.mssqlL)

	case "mysql":
		return lua.NewFunction(b.mysqlL)

	case "oracle":
		return lua.NewFunction(b.oracleL)

	case "postgres":
		return lua.NewFunction(b.postgresL)

	case "rdp":
		return lua.NewFunction(b.rdpL)

	case "redis":
		return lua.NewFunction(b.redisL)

	case "smb":
		return lua.NewFunction(b.smbL)

	case "smtp":
		return lua.NewFunction(b.smtpL)

	case "snmp":
		return lua.NewFunction(b.snmpL)

	case "web":
		return lua.NewFunction(b.webL)

	default:
		//其他协议

	}

	return nil
}

func (b *brute) ret(L *lua.LState) int {
	L.Push(L.NewProcData(b))
	return 1
}

//func (s *service) ret(L *lua.LState) int {
//	L.Push(L.NewLightUserData(s))
//	return 1
//}

func (b *brute) checkVM(L *lua.LState) bool {
	cu, nu := b.cfg.co.CodeVM(), L.CodeVM()
	if cu != nu {
		L.RaiseError("%s proc start must be %s , but %s", b.Name(), cu, nu)
		return false
	}
	return true
}

func (b *brute) useL(L *lua.LState) int {
	n := L.GetTop()
	m := &memory{}
	for i := 1; i <= n; i++ {
		m.use = append(m.use, L.Get(i).String())
	}
	if b.cfg.dict == nil {
		b.cfg.dict = m.Iterator()
	} else {
		//m.up(b.cfg.dict.Iterator())
		b.cfg.dict.Iterator().Updatem(m)
	}

	//a.Update(m)

	return b.ret(L)
	//读取堆栈
	//写入
	//b.cfg.dict.Iterator().Next()
}

func (b *brute) useFL(L *lua.LState) int {
	userp := L.Get(1).String()
	if !b.checkfile(userp) {
		fmt.Printf("error! open file :%v  error.", userp)
		return b.ret(L)
	}
	m := &fileM{
		userf: "",
		passf: "",
	}
	m.userf = userp
	if b.cfg.dict == nil {
		b.cfg.dict = m.Iterator()
	} else {
		b.cfg.dict.Iterator().Updatef(m)
	}

	//读取路径
	//赋值
	//条件: 如果 user.txt 文件大小 小于 dict.buffer  将全部缓存 user []string
	return b.ret(L)
}

func (b *brute) passL(L *lua.LState) int {
	n := L.GetTop()

	m := &memory{}

	for i := 1; i <= n; i++ {
		m.pass = append(m.pass, L.Get(i).String())
	}
	if b.cfg.dict == nil {
		b.cfg.dict = m.Iterator()
	} else {
		b.cfg.dict.Iterator().Updatem(m)
	}
	//b.cfg.dict = m.Iterator()
	return b.ret(L)
}

func (b *brute) passFL(L *lua.LState) int {
	passp := L.Get(1).String()
	if !b.checkfile(passp) {
		fmt.Printf("error! open file :%v  error.", passp)
		return b.ret(L)
	}
	m := &fileM{
		userf: "",
		passf: "",
	}
	m.passf = passp
	if b.cfg.dict == nil {
		b.cfg.dict = m.Iterator()
	} else {
		b.cfg.dict.Iterator().Updatef(m)
	}

	//读取堆栈
	//写入
	//b.config.dict.pass
	//条件: 如果 pass.txt 文件大小 小于 dict.buffer  将全部缓存 pass []string
	return b.ret(L)
}

func (b *brute) checkfile(userf string) bool {
	f, e := os.Open(userf)
	defer f.Close()
	if e != nil {
		return false
	} else {
		return true
	}
}

func (b *brute) cidrL(L *lua.LState) int {
	b.cfg.cidr = cidr.Check(L)
	return b.ret(L)
}

func (b *brute) hostL(L *lua.LState) int {
	return 0
}

func (b *brute) startL(L *lua.LState) int {
	b.Start()
	return b.ret(L)
	//if b.checkVM(L) {
	//	xEnv.Start(b., func(err error) {
	//		L.RaiseError("%v", err)
	//	})
	//}
	//return b.ret(L)
}

//----------------协议处理--------------------
func (b *brute) mongodbL(L *lua.LState) int {
	b.append(newBruteMongodb(L))
	return b.ret(L)
}

func (b *brute) mssqlL(L *lua.LState) int {
	b.append(newBruteMssql(L))
	return b.ret(L)
}

func (b *brute) mysqlL(L *lua.LState) int {
	b.append(newBruteMysql(L))
	return b.ret(L)
}

func (b *brute) oracleL(L *lua.LState) int {
	b.append(newBruteOracle(L))
	return b.ret(L)
}

func (b *brute) postgresL(L *lua.LState) int {
	b.append(newBrutePostgres(L))
	return b.ret(L)
}

func (b *brute) rdpL(L *lua.LState) int {
	b.append(newBruteRdp(L))
	return b.ret(L)
}

func (b *brute) redisL(L *lua.LState) int {
	b.append(newBruteRedis(L))
	return b.ret(L)
}

func (b *brute) smbL(L *lua.LState) int {
	b.append(newBruteSmb(L))
	return b.ret(L)
}

func (b *brute) smtpL(L *lua.LState) int {
	b.append(newBruteSmtp(L))
	return b.ret(L)
}

func (b *brute) snmpL(L *lua.LState) int {
	b.append(newBruteSnmp(L))
	return b.ret(L)
}

func (b *brute) webL(L *lua.LState) int {
	b.append(newBruteWeb(L))
	return b.ret(L)
}

func (b *brute) sshL(L *lua.LState) int {
	s := newBruteSsh(L)
	b.append(s)
	return b.ret(L)
}

func (b *brute) ftpL(L *lua.LState) int {
	b.append(newBruteFtp(L))
	return b.ret(L)
}

func (b *brute) elasticL(L *lua.LState) int {
	b.append(newBruteElastic(L))
	return b.ret(L)
}
