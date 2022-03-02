package crackonline

import (
	"database/sql"
	"fmt"
	_ "github.com/netxfly/mysql"
	"github.com/rock-go/rock/logger"
	"strings"
)

func (o *online) mysqldial(ipstopf *bool, userstopf *bool, s *int, ip string, port int, user string, pass string) {
	dataSourceName := fmt.Sprintf("%v:%v@tcp(%v:%v)/%v?charset=utf8", user,
		pass, ip, port, "mysql")
	db, err := sql.Open("mysql", dataSourceName)
	if err == nil {
		defer db.Close()
		err = db.Ping()
		if err == nil {
			//println(pass)
			o.ev(ip, user, pass, port, "mysql hit")
			*userstopf = true
			goto retn
		} else {
			if !strings.Contains(err.Error(), "Error 1045") { //非密码错误问题
				logger.Errorf("mysql crack err , ip is baned : &s", err.Error())
				*userstopf = true
				*ipstopf = true
			}
			goto retn
		}
	}
	//ip有问题
	*userstopf = true
	*ipstopf = true

retn:
	*s += 1
}
