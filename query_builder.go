package gosqlbuilder

import "fmt"

type QueryBuilder interface {
	Reset()
	NextArg() string
}

type IntQueryBuilder struct {
	ArgCounter uint32
}

func (b *IntQueryBuilder) NextArg() string {
	s := fmt.Sprintf("$%d", b.ArgCounter)
	b.ArgCounter += 1
	return s
}

func (b *IntQueryBuilder) Reset() {
	b.ArgCounter = 1
}

type QuestionMarkQueryBuilder struct{}

func (b *QuestionMarkQueryBuilder) Reset() {}

func (b *QuestionMarkQueryBuilder) NextArg() string {
	return "?"
}
