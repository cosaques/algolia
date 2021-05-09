package indexer

type Index interface {
	Add(string)
	Len() int
	Top(int) []TopQuery
}

type TopQuery struct {
	Query string
	Count int
}
