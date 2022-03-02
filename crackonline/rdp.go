package crackonline

import (
	"github.com/22ke/gordp/glog"
	gor "github.com/22ke/gordp/login"
	"strconv"
	"strings"
)

func (o *online) rdpdial(ipstopf *bool, userstopf *bool, s *int, ip string, port int, user string, pass string) {
	target := ip + ":" + strconv.Itoa(port)
	var err error
	g := gor.NewClient(target, glog.NONE)
	//SSL协议登录测试
	err = g.LoginForSSL("", user, pass)
	if err == nil {
		*userstopf = true
		*ipstopf = true
		*s += 1
		return
	}
	//println(err.Error())
	if strings.Contains(err.Error(), "success") {
		//println("login success , ", target, " : ", user, " : ", pass)
		o.ev(ip, user, pass, port, "rdp hit")
		*userstopf = true
	}
	*s += 1
}
