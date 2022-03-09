package brute

import (
	"github.com/rock-go/rock/lua"
	"gopkg.in/olivere/elastic.v3"
)

type Elastic struct {
	scheme string
}

func newBruteElastic(L *lua.LState) service {
	e := &Elastic{scheme: "http"}
	return newService(L , e)
}

func (ec *Elastic) Login(ev *event) {
	_ , err := elastic.NewClient(elastic.SetURL(ec.scheme + "://" + ev.Server()),
		elastic.SetSniff(false),
		elastic.SetBasicAuth(ev.user, ev.pass),
	)

	if err == nil {
		ev.stat = Succeed
	} else {
		ev.stat = Fail
		ev.banner = err.Error()
	}
}

func (ec *Elastic) Name() string {
	return "elastic"
}

func (ec *Elastic) Index(L *lua.LState , key string) lua.LValue {
	return lua.LNil
}