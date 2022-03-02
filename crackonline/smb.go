package crackonline

import "github.com/stacktitan/smb/smb"

func (o *online) smbdial(ipstopf *bool, userstopf *bool, s *int, ip string, port int, user string, pass string) {
	options := smb.Options{
		Host:        ip,
		Port:        port,
		User:        user,
		Password:    pass,
		Domain:      "",
		Workstation: "",
	}

	session, err := smb.NewSession(options, false)
	if err == nil {
		session.Close()
		if session.IsAuthenticated {
			println(pass)
			o.ev(ip, user, pass, port, "smb hit")
			*userstopf = true
		}
	} else {
		//println(pass,err.Error())
	}
	*s += 1
}
