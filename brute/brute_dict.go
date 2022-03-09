package brute

type dictEntry struct {
	name string
	pass string
    over bool
}

type Iterator interface {
	SkipU() //stop user
	Skip() //stop
	Next() dictEntry
	Close() error
}

type Dict interface {
	Iterator() Iterator
}