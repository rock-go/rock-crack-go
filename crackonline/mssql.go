package crackonline

import (
	"database/sql"
	"fmt"
)

func (o *online) mssqldial(ipstopf *bool, userstop *bool, s *int, ip string, port int, user string, pass string) {
	dataSourceName := fmt.Sprintf("server=%v;port=%v;user id=%v;password=%v;database=%v", ip, port, user, pass, "master")

	db, err := sql.Open("mssql", dataSourceName)
	if err == nil {
		defer db.Close()
		err = db.Ping()
		if err == nil {
			o.ev(ip, user, pass, port, "mssql hit")
			*userstop = true
		}
	} else {
		*ipstopf = true
	}
	*s += 1
}
