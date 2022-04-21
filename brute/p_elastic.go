package brute

import (
	"github.com/rock-go/rock/lua"
	"github.com/spf13/cast"
	"gopkg.in/olivere/elastic.v3"
)

type Elastic struct {
	scheme string
}

func (ec *Elastic) Name() string {
	return "elastic"
}

func newBruteElastic(L *lua.LState) service {
	val := L.CheckTable(1)
	port := cast.ToInt(val.RawGetString("port").String())

	e := &Elastic{scheme: "http"}
	return newService(L, e, port)
}

func (ec *Elastic) Login(ev *event) {
	_, err := elastic.NewClient(elastic.SetURL(ec.scheme+"://"+ev.Server()),
		elastic.SetSniff(false),
		elastic.SetBasicAuth(ev.user, ev.pass),
	)

	if err == nil {
		ev.stat = Succeed
		println("success: ", ev.user, ev.pass)
	} else {
		ev.stat = Fail
		//println("fail: ",ev.user,ev.pass)
		ev.banner = err.Error()
	}
}
