package brute

import (
	"fmt"
	"github.com/rock-go/rock/cidr"
	"github.com/rock-go/rock/logger"
	"github.com/rock-go/rock/lua"
	"github.com/rock-go/rock/pipe"
	"gopkg.in/tomb.v2"
	"net"
	"sync"
)

type brute struct {
	lua.LFace
	cfg      *config
	queue    chan Tx
	tom      *tomb.Tomb
	ipskip   map[string]bool
	userskip map[string]bool
}

func newBrute(cfg *config) *brute {
	b := &brute{cfg: cfg}
	//b.V(lua.INIT , typeof)
	return b
}

func (b *brute) append(s service) {
	//println(s.port)
	b.cfg.service = append(b.cfg.service, s)
}

func (b *brute) succeed(ev *event) {
	pipe.Do(b.cfg.succeed, ev, b.cfg.co, func(err error) {
		xEnv.Errorf("%s call succeed pipe fail %v", b.Name(), err)
	})
}

func (b *brute) verbose(ev *event) {
	//e := audit.NewEvent("crackonline").User(ev.user).Msg("ip:%s user:%s pass:%s port:%d", ev.ip, ev.user, ev.pass, ev.port)
	//e.Subject(ev.service).From(b.tom)
	//if o.cfg.pipe == nil {
	//	return
	//}
	//cp := xEnv.P(o.cfg.pipe)
	//co := xEnv.Clone(o.cfg.co)
	//defer xEnv.Free(co)
	//err := co.CallByParam(cp, ev)
	//xEnv.Errorf("%v", err)

	pipe.Do(b.cfg.verbose, ev, b.cfg.co, func(err error) {
		xEnv.Errorf("%s call verbose pipe fail %v", b.Name(), err)
	})
}

func (b *brute) help(s service) func(net.IP) {
	fn := func(ip net.IP) {
		//println("help fn")
		//路由是否可达
		if !s.Ping(ip) {
			logger.Errorf("IP %v , port: %v can not connect! ", ip, s.port)
			b.verbose(&event{
				ip:   ip.String(),
				port: s.port,
				stat: Unreachable,
			})
			return
		}

		//开始遍历字典
		i := 1
		iter := b.cfg.dict.Iterator()
		defer iter.Close()
		for info := iter.Next(); !info.over; info = iter.Next() {
			//fmt.Printf("%v : %#v\n",i, info)
			if i%1000 == 0 {
				fmt.Printf("%v : %#v\n", i, info)
			}
			i++
			select {

			case <-b.tom.Dying():
				println("dying")
				return

			default:
				if b.ipskip[ip.String()] == true {
					iter.Skip()
					println("ip stop!")
					break
				}
				//println("in queue")
				b.queue <- Tx{ip: ip, info: info, iter: iter, service: s}
			}
		}
	}

	return fn
}

func (b *brute) async() {
	//wg := sync.WaitGroup{}

	n := len(b.cfg.service)
	//print(b.cfg.service[1].port)
	//println("n：",n)
	if n == 0 {
		return
	}

	//wg.Add(n)
	for i := 0; i < n; i++ {
		//println(i)
		go func(s service) {
			cidr.Visit(b.tom, b.cfg.cidr, b.help(s)) //这里会阻塞
			//wg.Done()
		}(b.cfg.service[i])
	}
	//wg.Wait()
}

func (b *brute) Start() error {
	b.tom = new(tomb.Tomb)
	b.queue = make(chan Tx, 1024)
	go b.async()
	wg := &sync.WaitGroup{}
	for i := 0; i < b.cfg.thread; i++ {
		wg.Add(1)
		go b.thread(i, wg)
	}
	//wg.Wait()
	return nil
}

func (b *brute) Close() error {
	b.tom.Kill(fmt.Errorf("close"))
	close(b.queue)
	return nil
}

func (b *brute) thread(idx int, wg *sync.WaitGroup) {
	xEnv.Errorf("b thread %d start", idx)
	defer func() {
		xEnv.Errorf("b thread %d close", idx)
	}()
	//println(len(b.queue))
	for tx := range b.queue {
		//if checkskip(tx){
		//	return
		//}
		ev := &event{
			ip:      tx.ip.String(),
			user:    tx.info.name,
			pass:    tx.info.pass,
			port:    tx.service.port,
			service: tx.service.auth.Name(),
		}
		//fmt.Printf("%v/n",ev )

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
	wg.Done()

}
