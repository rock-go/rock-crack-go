package crackonline

import (
	"fmt"
	"github.com/go-redis/redis"
	"time"
)

func (o *online) redisdial(ipstopf *bool, userstop *bool, s *int, ip string, port int, user string, pass string) {
	opt := redis.Options{Addr: fmt.Sprintf("%v:%v", ip, port),
		Password: pass, DB: 0, DialTimeout: time.Duration(o.cfg.timeout) * time.Second}
	client := redis.NewClient(&opt)
	defer client.Close()
	_, err := client.Ping().Result()
	if err == nil {
		//println(pass)
		o.ev(ip, user, pass, port, "redis hit")
		*ipstopf = true
		*userstop = true
	} else {
		//println(pass,err.Error())
	}

	*s += 1
}
