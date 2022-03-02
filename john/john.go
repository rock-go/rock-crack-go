package john

import (
	"bufio"
	"crypto/md5"
	"crypto/sha256"
	"crypto/sha512"
	"fmt"
	"github.com/rock-go/rock/audit"
	"github.com/rock-go/rock/lua"
	"github.com/rock-go/rock/pipe"
	"github.com/tredoe/osutil/user/crypt"
	"github.com/tredoe/osutil/user/crypt/md5_crypt"
	"github.com/tredoe/osutil/user/crypt/sha256_crypt"
	"github.com/tredoe/osutil/user/crypt/sha512_crypt"
	"hash"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"
)

const (
	MD5 uint8 = iota + 1
	SHA256
	SHA512
	SHADOW
)

var fileTypeOf = reflect.TypeOf((*john)(nil)).String()

type john struct {
	lua.Super
	cfg *config
}

func newJohn(cfg *config) *john {
	obj := &john{cfg: cfg}
	obj.V(lua.MODE, time.Now())
	return obj
}

func (j *john) Start() error {
	return nil
}

func (j *john) Close() error {
	return nil
}

func (j *john) ret(L *lua.LState) int {
	L.Push(L.NewLightUserData(j))
	return 1
}

func (j *john) compareVM(co1 *lua.LState, co2 *lua.LState) bool {
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

func (j *john) Index(L *lua.LState, key string) lua.LValue {
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
			j.attack(MD5, co.IsString(1))
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
	if j.cfg.pipe == nil {
		return
	}
	pipe.Do(j.cfg.pipe, ev, j.cfg.co, func(err error) {
		xEnv.Errorf("%s pipe call fail %v", j.Name(), err)
	})
}

func (j *john) ifdictpath() bool {
	passf, err := os.Open(j.cfg.dict)
	if err != nil {
		return false
	}
	defer passf.Close()
	return true
}

func (j *john) shadow(raw string) {
	//1. 首先解析shadow raw 字符串
	//2. 开始爆破
	//3. 命中后运行pipe中的逻辑
	/*4.
	ev := audit.NewEvent("john").User(u).Msg("hash:%s pass:%s" , hash , pass)
	j.call(ev)
	*/
	//root:$6$X7Z9HGT8$.810fZP6mWm19PKSboWRLqCjGFyrH5doETlIqfPiPxQtCKFH2ecvG/xxtMdzE0pJG.amPTz5W/21/kJQ0O3Wl0:18896:0:99999:7:::

	//获取加密方式
	passtype := strings.Split(raw, "$")
	if len(passtype) < 4 {
		return
	}
	salt := "$" + passtype[1] + "$" + passtype[2] + "$"
	tp := passtype[1]

	t, err := strconv.Atoi(tp)
	if err != nil {
		panic(err)
	}
	var cryp crypt.Crypter
	switch t {
	case 1:
		cryp = md5_crypt.New()
	case 5:
		cryp = sha256_crypt.New()
	case 6:
		cryp = sha512_crypt.New()
	default:
		panic("nil cryp")
	}
	//获取加密shadow
	passhash := strings.Split(raw, ":")
	if len(passhash) < 4 {
		panic("length shadow err")
	}
	user := passhash[0]
	hashedpass := passhash[1]

	//将密码字典进行加密并比较
	var ismatch bool
	var plain string
	err, ismatch, plain = j.checkshadow(cryp, hashedpass, salt)
	if err != nil {
		panic(err)
	}
	if ismatch {
		ev := audit.NewEvent("john").User(user).Msg("hash:%s pass:%s", hashedpass, plain)
		ev.Subject("shadow hit").From(j.cfg.co.CodeVM())
		j.call(ev)
	}
}

func (j *john) checkshadow(crypt crypt.Crypter, hashedpass string, salt string) (error, bool, string) {
	var ismatch bool
	var plain string

	passf, err := os.Open(j.cfg.dict)
	defer passf.Close()
	if err != nil {
		return checkshadowstr(crypt, hashedpass, salt, j.cfg.dict)
	}
	passfsc := bufio.NewScanner(passf)
	for passfsc.Scan() {
		err, ismatch, plain = checkshadowstr(crypt, hashedpass, salt, passfsc.Text())
		if ismatch {
			return nil, true, plain
		}
	}
	return nil, false, ""
}

func checkshadowstr(crypt crypt.Crypter, hashed string, salt string, plainpass string) (error, bool, string) {
	newhash, err := crypt.Generate([]byte(plainpass), []byte(salt))
	if err != nil {
		return err, false, ""
	}
	if newhash == hashed {
		return nil, true, plainpass
	}
	return nil, false, ""
}

func (j *john) checkcrypt(h hash.Hash, raw string) (bool, string) {
	var ismatch bool
	var plain string

	passf, err := os.Open(j.cfg.dict)
	defer passf.Close()
	if err != nil {
		return j.checkcryptstr(h, j.cfg.dict, raw)
	}
	passfsc := bufio.NewScanner(passf)
	for passfsc.Scan() {
		ismatch, plain = j.checkcryptstr(h, passfsc.Text(), raw)
		if ismatch == true {
			return true, plain
		}
	}
	return false, ""
}

func (j *john) checkcryptstr(h hash.Hash, src string, raw string) (bool, string) {
	salt := j.cfg.salt
	h.Write([]byte(src))
	if len(salt) != 0 {
		h.Write([]byte(salt))
	}
	if fmt.Sprintf("%x", h.Sum(nil)) == raw {
		h.Reset()
		return true, src
	}
	h.Reset()
	return false, ""
}

func (j *john) md5(raw string) { //raw : eeda50edb56d...
	var ismatch bool
	var plain string
	h := md5.New()
	ismatch, plain = j.checkcrypt(h, raw)
	if !ismatch && plain != "" {
		panic(plain)
	}
	if ismatch {
		ev := audit.NewEvent("john").User("").Msg("md5:%s pass:%s", raw, plain)
		ev.Subject("md5 hit").From(j.cfg.co.CodeVM())
		j.call(ev)
	}
}

func (j *john) sha256(raw string) {
	var ismatch bool
	var plain string
	h := sha256.New()

	ismatch, plain = j.checkcrypt(h, raw)
	if !ismatch && plain != "" {
		panic(plain)
	}
	if ismatch {
		ev := audit.NewEvent("john").User("").Msg("sha256:%s pass:%s", raw, plain)
		ev.Subject("sha256 hit").From(j.cfg.co.CodeVM())
		j.call(ev)
	}
}

func (j *john) sha512(raw string) {
	var ismatch bool
	var plain string
	h := sha512.New()

	ismatch, plain = j.checkcrypt(h, raw)
	if !ismatch && plain != "" {
		panic(plain)
	}
	if ismatch {
		ev := audit.NewEvent("john").User("").Msg("sha512:%s pass:%s", raw, plain)
		ev.Subject("sha512 hit").From(j.cfg.co.CodeVM())
		j.call(ev)
	}
}

func (j *john) pipe(L *lua.LState) int {
	pv := pipe.LValue(L.Get(1))
	if pv == nil {
		return 0
	}

	j.cfg.pipe = append(j.cfg.pipe, pv)
	return 0
}

func (j *john) dict(L *lua.LState) int {
	//1. 判断是ext 后缀是否为 txt dict 等文件路径
	//2. 如果是文件 运行时打开io
	//3. 如果是文本 运行是 strings.NewReader("xxxxx")
	return j.ret(L)
}

func (j *john) attack(method uint8, raw string) {
	if raw == "" {
		return
	}

	//hash方式  $pass$salt
	switch method {
	case MD5:
		j.md5(raw)
	case SHA256:
		j.sha256(raw)
	case SHA512:
		j.sha512(raw)
	case SHADOW:
		j.shadow(raw)
	}
}
