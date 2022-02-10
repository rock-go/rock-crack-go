package john

import (
	"github.com/rock-go/rock/audit"
	"github.com/rock-go/rock/lua"
	"reflect"
	"time"
)

const (
	MD5    uint8 = iota + 1
	SHA256
	SHA512
	SHADOW
)

var fileTypeOf = reflect.TypeOf((*john)(nil)).String()

type john struct {
	lua.Super
	cfg    *config
}

func newJohn( cfg *config) *john {
	obj := &john{ cfg: cfg }
	obj.V(lua.MODE , time.Now())
	return obj
}

func (j *john) Start() error {
	return nil
}

func (j *john) Close() error {
	return nil
}

func (j *john) compareVM(co1 *lua.LState , co2 *lua.LState) bool {
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

func (j *john) Index(L *lua.LState , key string) lua.LValue {
	switch key {


	case "pipe":
		return L.NewFunction(j.pipe)

	case "dict":
		return L.NewFunction(j.dict)

	case "shadow":
		return L.NewFunction(func(co *lua.LState) int {
			j.attack(SHADOW, co.IsString(1))
			return 0
		})

	case "md5":
		return L.NewFunction(func(co *lua.LState) int {
			j.attack(MD5 , co.IsString(1))
			return 0
		})

	case "sha256":
		return L.NewFunction(func(co *lua.LState) int {
			j.attack(SHA256, co.IsString(1))
			return 0
		})

	case "sha512":
		return L.NewFunction(func(co *lua.LState) int {
			j.attack(SHA512, co.IsString(1))
			return 0
		})
	}

	return nil
}

func (j *john) call(ev *audit.Event) {
	cp := xEnv.P(j.cfg.pipe)
	co := xEnv.Clone(j.cfg.co)
	defer xEnv.Free(co)
	err := co.CallByParam(cp , ev)
	xEnv.Errorf("%v" , err)

}

func (j *john) shadow(raw string) {
	//1. 首先解析shadow raw 字符串
	//2. 开始爆破
	//3. 命中后运行pipe中的逻辑
	/*4.
		ev := audit.NewEvent("john").User(u).Msg("hash:%s pass:%s" , hash , pass)
		j.call(ev)
	*/
}

func (j *john) md5(raw string){
	//shadow 类似
}

func (j *john) sha256(raw string) {
	//shadow 类似
}

func (j *john) sha512(raw string) {
	//shadow 类似
}

func (j *john) pipe(L *lua.LState) int {
	j.cfg.pipe = L.IsFunc(1)
	return 0
}

func (j *john) dict(L *lua.LState) int {
	//1. 判断是ext 后缀是否为 txt dict 等文件路径
	//2. 如果是文件 运行时打开io
	//3. 如果是文本 运行是 strings.NewReader("xxxxx")
	return 0
}

func (j *john) attack(method uint8 , raw string) {
	if raw == "" {
		return
	}

	switch method {
	case MD5:    j.md5(raw)
	case SHA256: j.sha256(raw)
	case SHA512: j.sha512(raw)
	case SHADOW: j.shadow(raw)
	}
}