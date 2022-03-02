package crackonline

import (
	"github.com/gosnmp/gosnmp"
	"time"
)

func (o *online) snmpdial(ipstopf *bool, userstop *bool, s *int, ip string, port int, user string, pass string) {
	gosnmp.Default.Target = ip
	gosnmp.Default.Port = uint16(port)
	gosnmp.Default.Community = pass
	gosnmp.Default.Timeout = time.Duration(o.cfg.timeout) * time.Second

	err := gosnmp.Default.Connect()
	if err == nil {
		oids := []string{"1.3.6.1.2.1.1.4.0", "1.3.6.1.2.1.1.7.0"}
		_, err := gosnmp.Default.Get(oids)
		if err == nil {
			println(pass)
			o.ev(ip, user, pass, port, "redis hit")
		} else {
			println(pass, err.Error())
		}
	} else {
		println(pass, err.Error())
	}
}
