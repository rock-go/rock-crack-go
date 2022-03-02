package crackonline

import (
	"fmt"
	"gopkg.in/olivere/elastic.v3"
)

func (o *online) elasticdial(ipstopf *bool, userstopf *bool, s *int, ip string, port int, user string, pass string) {
	_, err := elastic.NewClient(elastic.SetURL(fmt.Sprintf("http://%v:%v", ip, port)),
		elastic.SetMaxRetries(3),
		elastic.SetSniff(false),
		elastic.SetBasicAuth(user, pass),
	)
	if err == nil {
		//println(pass)
		o.ev(ip, user, pass, port, "elastic hit")
		*userstopf = true
		//_, _, err = client.Ping(fmt.Sprintf("http://%v:%v", ip, port)).Do()
		//if err == nil {
		//	println(pass)
		//	o.ev(ip,user,pass,port,"elastic hit")
		//	*userstopf = true
		//}
	} else {
		//println(user,pass,err.Error())
	}
	*s += 1
}
