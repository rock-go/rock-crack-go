package john

import (
	"github.com/rock-go/rock/auxlib"
	"github.com/rock-go/rock/lua"
	"github.com/rock-go/rock/pipe"
)

type config struct {
	speed int
	name  string
	dict  string
	salt  string

	co   *lua.LState
	pipe []pipe.Pipe
}

func (c config) verify() interface{} {
	return nil
}

func (c *config) Index(L *lua.LState, key string, val lua.LValue) {
	switch key {
	case "name":
		c.name = auxlib.CheckProcName(val, L)

	case "speed":
		c.speed = lua.IsInt(val)

	case "dict":
		c.dict = lua.IsString(val)

	case "pipe":
		if pv := pipe.LValue(val); val == nil {
			L.RaiseError("invalid pipe type")
		} else {
			c.pipe = []pipe.Pipe{pv}
		}

	case "salt":
		c.salt = lua.IsString(val)

	default:
		L.RaiseError("invalid %s field", key)
		return
	}
}

func newConfig(L *lua.LState) *config {
	val := L.Get(1)
	cfg := &config{
		co: xEnv.Clone(L),
	}

	switch val.Type() {
	case lua.LTString:
		cfg.name = auxlib.CheckProcName(val, L)

	case lua.LTTable:
		val.(*lua.LTable).Range(func(key string, val lua.LValue) { cfg.Index(L, key, val) })

	default:
		L.RaiseError("invalid config type must string or table , got %s", val.Type().String())
		return nil
	}

	if e := cfg.verify(); e != nil {
		L.RaiseError("%v", e)
		return nil
	}

	return cfg

}
