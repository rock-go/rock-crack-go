package brute

import (
	"sync/atomic"
)

type iteratorM struct {
	object   *memory

	uSize    uint32
	pSize    uint32

	uSeek    uint32
	pSeek    uint32
}

type memory struct {
	use     []string
	pass    []string
}

func (m *memory) Iterator() *iteratorM {
	return &iteratorM{
		object: m,
		uSeek: 0,
		pSeek: 0,
		uSize: uint32(len(m.use)),
		pSize: uint32(len(m.pass)),
	}
}

func (iter *iteratorM) SkipU() {
	atomic.AddUint32(&iter.uSeek , 1)
}

func (iter *iteratorM) Skip() {
	atomic.AddUint32(&iter.uSeek , iter.uSize + 1)
	atomic.AddUint32(&iter.pSize , iter.pSize + 1)
}

func (iter *iteratorM) Next() dictEntry {

	u := atomic.LoadUint32(&iter.uSeek)

next:
	if u >= iter.uSize {
		return dictEntry{over: true}
	}

	p := atomic.AddUint32(&iter.pSeek , 1)
	if p < iter.pSize {
		return dictEntry{
			name: iter.object.use[u],
			pass: iter.object.pass[p],
			over: false,
		}
	}

	u = atomic.AddUint32(&iter.uSeek , 1)
	goto next
}

func (iter *iteratorM) Close() error {
	iter.Skip()
	return nil
}