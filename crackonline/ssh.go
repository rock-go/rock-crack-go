package crackonline

import (
	"fmt"
	"golang.org/x/crypto/ssh"
	"net"
	"time"
)

//ssh     root@127.0.0.1  -p 123456
func (o *online) sshdial(ipstopf *bool, userstopf *bool, s *int, ip string, port int, user string, pass string) {
	config := &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{
			ssh.Password(pass),
		},
		Timeout: time.Duration(o.cfg.timeout) * time.Second,
		//HostKeyCallback: ssh.FixedHostKey(sshpublickey),
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			return nil
		},
	}
	client, err := ssh.Dial("tcp", fmt.Sprintf("%v:%v", ip, port), config)
	if err == nil {
		defer client.Close()
		o.ev(ip, user, pass, port, "ssh hit")
		*userstopf = true
	} else { //密码错误
		//println(pass , " : ",err.Error())
	}
	*s += 1
}
