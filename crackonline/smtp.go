package crackonline

import (
	"crypto/tls"
	"github.com/rock-go/rock/logger"
	smtp "net/smtp"
	"strconv"
	"strings"
)

func (o *online) smtpdial(ipstopf *bool, userstop *bool, s *int, ip string, port int, user string, pass string) {
	str := ip + ":" + strconv.Itoa(port)
	c, err := smtp.Dial(str)
	if err != nil {
		//println("dial",err.Error())
		logger.Errorf("dial %s err : %s", ip, err.Error())
		*userstop = true
		*ipstopf = true
		*s += 1
		return
	}
	auth := smtp.PlainAuth("", user, pass, ip)

	if ok, _ := c.Extension("STARTTLS"); ok {
		config := &tls.Config{ServerName: ip, InsecureSkipVerify: true}
		if err = c.StartTLS(config); err != nil {
			//fmt.Println("call start tls")
			//return err
		}
	}

	if auth != nil {
		if ok, _ := c.Extension("AUTH"); ok {
			if err = c.Auth(auth); err != nil {
				if strings.Contains(err.Error(), "504 Unrecognized authentication type") {
					logger.Errorf("smtp crack error: 线程数量过多！ %s", err.Error())
				}
				//密码错误
			} else {
				//println(user,pass)
				o.ev(ip, user, pass, port, "smtp hit")
				*userstop = true
			}
		}
	}
	*s += 1
}
