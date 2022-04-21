package brute

type dictEntry struct {
	name string
	pass string
	over bool
}

type Iterator interface {
	SkipU() //stop user
	Skip()  //stop
	Next() dictEntry
	Close() error
	Updatem(*memory)
	Updatef(*fileM)
}

type Dict interface {
	Iterator() Iterator
}

type Super struct{}

func (su *Super) Updatem(*memory) {}
func (su *Super) Updatef(*fileM)  {}
