package crackonline

import (
	"database/sql"
	"fmt"
)

func (o *online) postgresdial(ipstopf *bool, userstop *bool, s *int, ip string, port int, user string, pass string) {
	dataSourceName := fmt.Sprintf("postgres://%v:%v@%v:%v/%v?sslmode=%v", user,
		pass, ip, port, "postgres", "disable")
	db, err := sql.Open("postgres", dataSourceName)

	if err == nil {
		defer db.Close()
		err = db.Ping()
		if err == nil {
			o.ev(ip, user, pass, port, "postgres hit")
		}
	}
}
