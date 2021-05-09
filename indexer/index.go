package indexer

type Index interface {
	Add(string) error
	Len() int
	Top(int) []TopQuery
}

type TopQuery struct {
	Query string
	Count int
}
