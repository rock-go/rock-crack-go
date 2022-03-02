package crackonline

import (
	"bufio"
	"github.com/rock-go/rock/audit"
	"github.com/rock-go/rock/logger"
	"github.com/rock-go/rock/lua"
	"net"
	"os"
	"reflect"
	"strconv"
	"time"
)

var fileTypeOf = reflect.TypeOf((*online)(nil)).String()

const (
	SSH uint8 = iota + 1
	ELASTIC
	FTP
	MONGODB
	MSSQL
	MYSQL
	ORACLE
	POSTGRES
	REDIS
	SMB
	SNMP
	SMTP
	RDP
	VPN
)

type online struct {
	lua.Super
	cfg *config
}

func newJohn(cfg *config) *online {
	obj := &online{cfg: cfg}
	obj.V(lua.MODE, time.Now())
	return obj
}

func (o *online) Start() error {
	return nil
}

func (o *online) Close() error {
	return nil
}

func (o *online) ret(L *lua.LState) int {
	L.Push(L.NewLightUserData(o))
	return 1
}

func (o *online) compareVM(co1 *lua.LState, co2 *lua.LState) bool {
	if co1 == nil || co2 == nil {
		return false
	}

	vm1 := co1.CodeVM()
	vm2 := co2.CodeVM()

	if vm1 == "" || vm2 == "" {
		return false
	}

	return vm1 == vm2
}

//判断是否为文件，以能否打开作为判断依据
func (o *online) isfile(name string) *os.File {
	f, err := os.Open(name)
	if err != nil {
		logger.Infof("'%s' is not the file , so as the string to solve", name)
		return nil
	}
	return f
}

func (o *online) Index(L *lua.LState, key string) lua.LValue {
	o.cfg.method = key
	switch key {
	case "pipe":
		return L.NewFunction(o.pipe)

	case "dict":
		return L.NewFunction(o.dict)

	case "ssh":
		return L.NewFunction(func(co *lua.LState) int {
			o.crackonline(SSH, co.IsInt(1))
			return 0
		})
	case "mysql":
		return L.NewFunction(func(co *lua.LState) int {
			o.crackonline(MYSQL, co.IsInt(1))
			return 0
		})
	case "smb":
		return L.NewFunction(func(co *lua.LState) int {
			o.crackonline(SMB, co.IsInt(1))
			return 0
		})
	case "elastic":
		return L.NewFunction(func(co *lua.LState) int {
			o.crackonline(ELASTIC, co.IsInt(1))
			return 0
		})
	case "ftp":
		return L.NewFunction(func(co *lua.LState) int {
			o.crackonline(FTP, co.IsInt(1))
			return 0
		})
	case "mongodb":
		return L.NewFunction(func(co *lua.LState) int {
			o.crackonline(MONGODB, co.IsInt(1))
			return 0
		})
	case "mssql":
		return L.NewFunction(func(co *lua.LState) int {
			o.crackonline(MSSQL, co.IsInt(1))
			return 0
		})
	case "oracle":
		return L.NewFunction(func(co *lua.LState) int {
			o.crackonline(ORACLE, co.IsInt(1))
			return 0
		})
	case "postgres":
		return L.NewFunction(func(co *lua.LState) int {
			o.crackonline(POSTGRES, co.IsInt(1))
			return 0
		})
	case "redis":
		return L.NewFunction(func(co *lua.LState) int {
			o.crackonline(REDIS, co.IsInt(1))
			return 0
		})
	case "snmp":
		return L.NewFunction(func(co *lua.LState) int {
			o.crackonline(SNMP, co.IsInt(1))
			return 0
		})
	case "smtp":
		return L.NewFunction(func(co *lua.LState) int {
			o.crackonline(SMTP, co.IsInt(1))
			return 0
		})
	case "rdp":
		return L.NewFunction(func(co *lua.LState) int {
			o.crackonline(RDP, co.IsInt(1))
			return 0
		})
	case "vpn":
		return L.NewFunction(func(co *lua.LState) int {
			o.crackonline(VPN, co.IsInt(1))
			return 0
		})
	}

	return nil
}

func (o *online) call(ev *audit.Event) {
	if o.cfg.pipe == nil {
		return
	}
	cp := xEnv.P(o.cfg.pipe)
	co := xEnv.Clone(o.cfg.co)
	defer xEnv.Free(co)
	err := co.CallByParam(cp, ev)
	xEnv.Errorf("%v", err)
}

func (o *online) pipe(L *lua.LState) int {
	o.cfg.pipe = L.IsFunc(1)
	return o.ret(L)
}

func (o *online) dict(L *lua.LState) int {
	//1. 判断是ext 后缀是否为 txt dict 等文件路径
	//2. 如果是文件 运行时打开io
	//3. 如果是文本 运行是 strings.NewReader("xxxxx")
	return o.ret(L)
}

func (o *online) checkip(ip string, port int) bool {
	host := ip + ":" + strconv.Itoa(port)
	var d net.Dialer
	d.Timeout = time.Duration(o.cfg.timeout) * time.Second
	_, err := d.Dial("tcp", host)
	if err != nil {
		return false
	}
	return true
}

func (o *online) dealuserfile(speed int, f func(*bool, *bool, *int, string, int, string, string), ip string, port int) { //获取用户名
	s := speed

	//判断端口是否通
	if !o.checkip(ip, port) && o.cfg.method != "vpn" { //ip不可达并且不是vpn返回
		//println("不通：",ip,":",port)
		logger.Errorf("ip: %s , port: %d dialed fail", ip, port)
		return
	}
	ipstopf := false

	if o.cfg.method == "redis" {
		go o.dealpassfile(&ipstopf, &s, f, ip, port, o.cfg.user)
		return
	}
	//判断用户名是否为文件列表
	file := o.isfile(o.cfg.user)
	defer file.Close()
	if file != nil {
		userlist := bufio.NewScanner(file)
		for userlist.Scan() {
			if ipstopf {
				//println("该ip有问题:",ip)
				logger.Errorf("the ip : %s has problem", ip)
				return
			}
			go o.dealpassfile(&ipstopf, &s, f, ip, port, userlist.Text())
		}
	} else {
		if ipstopf {
			logger.Errorf("the ip : %s has problem", ip)
			//println("该ip有问题:",ip)
			return
		}
		go o.dealpassfile(&ipstopf, &s, f, ip, port, o.cfg.user)
	}
	return
}

func (o *online) dealpassfile(ipstopf *bool, s *int, f func(*bool, *bool, *int, string, int, string, string), ip string, port int, user string) { //获取密码
	userstopf := false
	//处理密码文件
	file := o.isfile(o.cfg.dict)
	defer file.Close()
	if file != nil {
		passlist := bufio.NewScanner(file)
		for passlist.Scan() {
			pass := passlist.Text()
			if userstopf || *ipstopf {
				if !*ipstopf {
					logger.Infof("the user %s cracked successfully or has error", user)
				}
				//logger.Infof("the user %s has")
				//println("该user不用爆破了:",user)
				return
			}
			//println(pass)
			o.crk(f, ipstopf, &userstopf, s, ip, port, user, pass)
		}
		logger.Infof("the user : %s cracked failed", user)
		//println(user," ：扫描完毕！")
	} else {
		if userstopf || *ipstopf {
			if !*ipstopf {
				logger.Infof("the user %s cracked successfully or has error", user)
			}
			return
		}
		o.crk(f, ipstopf, &userstopf, s, ip, port, user, o.cfg.dict)
		logger.Infof("the user : %s cracked failed", user)
		//println(user," ：扫描完毕！")
	}
}

func (o *online) crk(f func(*bool, *bool, *int, string, int, string, string), ipstopf *bool, userstopf *bool, s *int, ip string, port int, user string, pass string) {
	for true {
		if *s > 0 {
			//println(*s)
			*s -= 1
			go f(ipstopf, userstopf, s, ip, port, user, pass)
			break
		}
		time.Sleep(1)
	}
}

func (o *online) crack(f func(*bool, *bool, *int, string, int, string, string), port int) { //获取ip//一个线程处理一个ip
	file := o.isfile(o.cfg.iplist)
	defer file.Close()
	if file != nil {
		iplist := bufio.NewScanner(file)
		for iplist.Scan() {
			go o.dealuserfile(o.cfg.threads, f, iplist.Text(), port)
		}
	} else {
		go o.dealuserfile(o.cfg.threads, f, o.cfg.iplist, port)
	}
}

func (o *online) ev(ip string, user string, pass string, port int, subject string) {
	ev := audit.NewEvent("crackonline").User(user).Msg("ip:%s user:%s pass:%s port:%d", ip, user, pass, port)
	ev.Subject(subject).From(o.cfg.co.CodeVM())
	o.call(ev)
}

func (o *online) crackonline(method uint8, port int) {
	var f func(*bool, *bool, *int, string, int, string, string)
	switch method {
	case SSH:
		f = o.sshdial
	case ELASTIC:
		f = o.elasticdial
	case FTP:
		f = o.ftpdial
	case MONGODB:
		f = o.mongodbdial
	case MSSQL:
		f = o.mssqldial
	case MYSQL:
		f = o.mysqldial
	case ORACLE:
		f = o.oracledial
	case POSTGRES:
		f = o.postgresdial
	case REDIS:
		f = o.redisdial
	case SMB:
		f = o.smbdial
	case SNMP:
		f = o.snmpdial
	case SMTP:
		f = o.smtpdial
	case RDP:
		f = o.rdpdial
	case VPN:
		f = o.vpndial
	default:
		panic("no such method")
	}
	o.crack(f, port)
}
