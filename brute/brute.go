package brute

import (
	"github.com/rock-go/rock/cidr"
	"github.com/rock-go/rock/lua"
	"github.com/rock-go/rock/pipe"
	"gopkg.in/tomb.v2"
	"net"
)

type brute struct {
	lua.Super
	cfg   *config
	queue chan Tx
	tom   *tomb.Tomb
}

func newBrute(cfg *config) *brute {
	b := &brute{ cfg:cfg }
	b.V(lua.INIT , typeof)
	return b
}

func (b *brute) append(s service) {
	b.cfg.service = append(b.cfg.service , s)
}

func (b *brute) succeed(ev *event) {
	pipe.Do(b.cfg.succeed , ev, b.cfg.co , func(err error) {
		xEnv.Errorf("%s call succeed pipe fail %v" , b.Name() , err)
	})
}

func (b *brute) verbose(ev *event) {
	pipe.Do(b.cfg.verbose , ev , b.cfg.co , func(err error) {
		xEnv.Errorf("%s call verbose pipe fail %v" , b.Name() , err)
	})
}

func (b *brute) help(s service) func(net.IP) {
	fn := func(ip net.IP) {

		//路由是否可达
		if !s.Ping(ip) {
			b.verbose(&event{
				ip: ip.String(),
				port: s.port,
				stat: Unreachable,
			})
			return
		}

		//开始遍历字典
		iter := b.cfg.dict.Iterator()
		defer iter.Close()
		for info := iter.Next(); !info.over ; info = iter.Next() {
			select {

			case <-b.tom.Dying():
				return

			default:
				b.queue <- Tx{ ip:ip , info: info , iter:iter , service: s }
			}
		}
	}

	return fn
}

func (b *brute) async() {
	n := len(b.cfg.service)
	if n == 0 {
		return
	}

	for i := 0 ; i < n ; i++ {
		go func(s service){
			cidr.Visit(b.tom , b.cfg.cidr , b.help(s)) //这里会阻塞
		}(b.cfg.service[i])
	}
}

func (b *brute) Start() error {
	return nil
}

func (b *brute) Close() error {
	return nil
}

func (b *brute) thread(idx int) {
	for tx := range b.queue {
		ev := &event{
			ip: tx.ip.String(),
			user: tx.info.name,
			pass: tx.info.pass,
			service: tx.service.auth.Name(),
		}

		tx.service.Do(ev)

		switch ev.stat {

		case Succeed:
			b.succeed(ev)
			if tx.service.skip {
				tx.iter.Skip()
				goto done
			}

			//跳过当前用户名
			tx.iter.SkipU()

		case Denied:
			//用户被锁定 跳过
			tx.iter.SkipU()

		case Fail:
			//todo

		case Unreachable:
			//todo
		}

		done:
		b.verbose(ev)
	}

}