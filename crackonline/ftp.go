package crackonline

import (
	"fmt"
	"github.com/jlaffaye/ftp"
	"strings"
	"time"
)

func (o *online) ftpdial(ipstopf *bool, userstopf *bool, s *int, ip string, port int, user string, pass string) {
	conn, err := ftp.DialTimeout(fmt.Sprintf("%v:%v", ip, port), time.Duration(o.cfg.timeout)*time.Second)
	if err == nil {
		err = conn.Login(user, pass)
		if err == nil {
			defer conn.Logout()
			//println(pass)
			o.ev(ip, user, pass, port, "ftp hit")
			*userstopf = true
		} else {
			//println(user,":",pass,err.Error())
			if strings.Contains(err.Error(), "Permission denied") {
				*userstopf = true
			}
		}
	} else {
		//该ip被封了
		*ipstopf = true
		//println(err.Error())
	}

	*s += 1
}
