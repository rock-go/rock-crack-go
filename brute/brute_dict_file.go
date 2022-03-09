package brute

import (
	"bufio"
	"github.com/rock-go/rock/auxlib"
	"io"
	"os"
)

type iteratorM2 struct {
	u        *bufio.Reader
	p        *bufio.Reader

	uFile    *os.File
	pFile    *os.File

	name     string
	pass     string

	ue       error
	pe       error
}

type fileM struct {
	user  string
	pass  string
}

func (f *fileM) Iterator() *iteratorM2 {
	iter := &iteratorM2{}

	u , ue := os.Open(f.user)
	if ue == nil {
		iter.u = bufio.NewReaderSize(u , 64)
		iter.ru()
	}

	p , pe := os.Open(f.pass)
	if pe == nil {
		iter.p = bufio.NewReaderSize(p , 64)
	}

	iter.uFile = u
	iter.pFile = p
	iter.ue = ue
	iter.pe = pe
	return iter
}

func (iter *iteratorM2) ru() error {
	if iter.u == nil || iter.ue != nil {
		return iter.ue
	}

	name , err := iter.u.ReadBytes('\n')
	if err == nil {
		iter.name = auxlib.B2S(name[1:])
		return nil
	}

	iter.ue = err
	return err
}

func (iter *iteratorM2) pu() error {
	if iter.p == nil || iter.pe != nil {
		return iter.pe
	}

	pass , err := iter.u.ReadBytes('\n')
	if err == nil {
		iter.pass = auxlib.B2S(pass[1:])
		return nil
	}

	iter.pe = err
	return err
}

func (iter *iteratorM2) SkipU() {
	iter.ru()
}

func (iter *iteratorM2) Skip() {
	iter.ue = io.EOF
	iter.pe = io.EOF
	iter.name = ""
	iter.pass = ""
}

func (iter *iteratorM2) Next() dictEntry {

next:
	if iter.ue != nil {
		return dictEntry{ over: true }
	}

	err := iter.pu()
	if err == nil {
		return dictEntry{ name: iter.name, pass: iter.pass, over: false }
	}

	if err == io.EOF {
		//下一个用户名
		iter.ru()
		goto next
	}

	iter.pe = err
	return dictEntry{over: true}
}

func (iter *iteratorM2) Close() error {
	iter.Skip()

	if iter.uFile != nil {
		iter.uFile.Close()
	}

	if iter.pFile != nil {
		iter.pFile.Close()
	}
	return nil
}