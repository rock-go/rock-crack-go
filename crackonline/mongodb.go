package crackonline

import (
	"fmt"
	"github.com/rock-go/rock/logger"
	"gopkg.in/mgo.v2"
	"time"
)

func (o *online) mongodbdial(ipstopf *bool, userstop *bool, s *int, ip string, port int, user string, pass string) {
	url := fmt.Sprintf("mongodb://%v:%v@%v:%v/%v", user, pass, ip, port, "admin")
	session, err := mgo.DialWithTimeout(url, time.Duration(o.cfg.timeout)*time.Second)

	if err == nil {
		defer session.Close()
		err = session.Ping()
		if err == nil {
			//println(pass)
			o.ev(ip, user, pass, port, "mongodb hit")
			*userstop = true
		} else {
			//println("xxx:",user,pass,err.Error())
			logger.Error("ping mongodb err : %s", err.Error())
		}
	} else {
		//println(user,pass,err.Error())
		//logger.Errorf("connect mongodb err : %s",err.Error())
		//*userstop = true
		//*ipstopf = true
	}
	*s += 1
}
