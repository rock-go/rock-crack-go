package brute

import (
	"github.com/rock-go/rock/logger"
	"github.com/rock-go/rock/lua"
	"github.com/spf13/cast"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type web struct {
	url    string
	proxy  *url.URL
	method string
	//userfiled string
	//pass      string
	query       string
	checkstatus int
	checkstr    string
	contenttype string
	timeout     time.Duration
}

func (w *web) Name() string {
	return "web"
}

func checkconfig(sv *web) {

	m := strings.ToUpper(sv.method)
	if m != "GET" && m != "POST" {
		logger.Errorf("web method error: %s", sv.method)
		sv.method = ""
	}
	if !(strings.Contains(sv.query, "{user}") && strings.Contains(sv.query, "{pass}")) {
		logger.Errorf("web user or pass error: %s", sv.query)
	}
	if sv.checkstatus > 506 || sv.checkstatus < 100 {
		logger.Errorf("web status error: %s", sv.checkstatus)
	}
	if sv.checkstr == "nil" {
		sv.checkstr = ""
	}

	return
}

//func delheader(opt *lua.LTable) http.Header {
//	for k,v := range opt.Array(){
//
//	}
//}

func newBruteWeb(L *lua.LState) service {
	opt := L.CheckTable(1)
	port := 80

	p := opt.RawGetString("proxy").String()
	u, e := url.Parse(p)
	if e != nil {
		logger.Errorf("web proxy error: %s ; %s", p, e.Error())
		u = nil
	}

	status, err := strconv.Atoi(opt.RawGetString("checkstatus").String())
	if err != nil {
		logger.Errorf("web status error: %s ; %s", status, err.Error())
		status = 200
	}

	sv := &web{
		url:         opt.RawGetString("url").String(),
		proxy:       u,                                   //"http://172.31.61.82:8080"
		method:      opt.RawGetString("method").String(), //"GET"
		timeout:     time.Duration(cast.ToInt(opt.RawGetString("timeout").String())),
		query:       opt.RawGetString("query").String(), //""
		contenttype: opt.RawGetString("contenttype").String(),
		checkstatus: status,
		checkstr:    opt.RawGetString("checkstr").String(),
	}

	checkconfig(sv)

	//proxy := opt.RawGetString("proxy").String()//"http://172.31.61.82:8080"

	s := newService(L, sv, port)
	s.ping = false
	return s
}

func (w *web) Login(ev *event) {
	//println("WEB LOGIN")
	//新建请求客户端
	c := http.Client{
		Timeout: 5 * time.Second,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	//设置代理
	if w.proxy != nil {
		c.Transport = &http.Transport{
			Proxy: http.ProxyURL(w.proxy),
		}
	}

	b := w.query
	b = strings.ReplaceAll(b, "{user}", ev.user)
	b = strings.ReplaceAll(b, "{pass}", ev.pass)

	if strings.ToUpper(w.method) == "POST" {
		resp, err := c.Post(
			w.url,
			w.contenttype,
			strings.NewReader(b),
		)
		if err != nil {
			//logger.Errorf("post body fial %s", err.Error())
			ev.stat = Unreachable
			return
		}
		status := resp.StatusCode
		r, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			logger.Errorf("")
		}
		//println(string(r))
		if status == w.checkstatus && strings.Contains(string(r), w.checkstr) {
			println("wright! :", ev.user, " : ", ev.pass)
			ev.stat = Succeed
			return
		}
		//println("wrong! :", ev.user, " : ",ev.pass)
	}

}
