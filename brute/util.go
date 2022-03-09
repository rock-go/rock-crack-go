package brute

import "github.com/rock-go/rock/lua"


func toPort(L *lua.LState) int {
	p := L.IsInt(1)
	if p > 0 && p < 65535 {
		return p
	}

	L.RaiseError("invalid port %d" , p)
	return p
}

func checkServicePorts(L *lua.LState) []int {
	var ports []int

	n := L.GetTop()
	if n == 0 {
		return ports
	}

	tmp := make(map[int]struct{})

	add := func(p int) {
		if p <= 0 || p > 65535 {
			return
		}

		if _ , ok := tmp[p]; ok {
			return
		}

		tmp[p] = struct{}{}
		ports = append(ports , p)
	}

	for i := 1 ; i <= n ; i++ {
		val := L.Get(i)

		switch val.Type() {

		case lua.LTInt:
			add(int(val.(lua.LInt)))

		case lua.LTNumber:
			add(int(val.(lua.LNumber)))

		case lua.LTString:
			//解析 10-100 ,1000-2000 的端口范围
			//todo
		}
	}

	return ports
}