package indexer

import "container/list"

type listIndex struct {
	data      map[*string]*indexData
	orderList *list.List
	toIndex   chan *string
	requests  chan receiver
}

type indexData struct {
	count     int
	orderedEl *list.Element
}

func NewIndex() *listIndex {
	i := &listIndex{
		data:      make(map[*string]*indexData),
		toIndex:   make(chan *string),
		requests:  make(chan receiver),
		orderList: list.New(),
	}
	go i.run()
	return i
}

func (i *listIndex) Add(s *string) {
	i.toIndex <- s
}

func (i *listIndex) Len() int {
	req := &distinctRequest{newReceiver()}
	i.requests <- req
	result := <-req.result

	return result.(int)
}

func (i *listIndex) Top(size int) []TopQuery {
	req := &topQueriesRequest{newReceiver(), size}
	i.requests <- req
	result := <-req.result

	return result.([]TopQuery)
}

func (i *listIndex) run() {
	for {
		// priority to requests
		select {
		case r := <-i.requests:
			switch req := r.(type) {
			case *distinctRequest:
				req.receive(len(i.data))
			case *topQueriesRequest:
				res := make([]TopQuery, 0, req.size)
				for c, el := req.size, i.orderList.Back(); c > 0 && el != nil; c, el = c-1, el.Prev() {
					res = append(res, TopQuery{*(el.Value.(*string)), i.data[el.Value.(*string)].count})
				}
				req.receive(res)
			}
			continue
		default:
		}

		// wait either for a request or an indexation demande
		select {
		case r := <-i.requests:
			switch req := r.(type) {
			case *distinctRequest:
				req.receive(len(i.data))
			case *topQueriesRequest:
				res := make([]TopQuery, 0, req.size)
				for c, el := req.size, i.orderList.Back(); c > 0 && el != nil; c, el = c-1, el.Prev() {
					res = append(res, TopQuery{*(el.Value.(*string)), i.data[el.Value.(*string)].count})
				}
				req.receive(res)
			}
		case s := <-i.toIndex:
			data, exists := i.data[s]
			if !exists {
				i.orderList.PushBack(s)
				i.data[s] = &indexData{1, i.orderList.Back()}
			} else {
				data.count++
				e := data.orderedEl
				for e.Next() != nil && i.data[e.Next().Value.(*string)].count < data.count {
					e = e.Next()
				}
				if data.orderedEl != e {
					i.orderList.MoveAfter(data.orderedEl, e)
				}
			}
		}
	}
}

type receiver interface {
	receive(interface{})
}

type baseReceiver struct {
	result chan interface{}
}

func newReceiver() *baseReceiver {
	return &baseReceiver{
		result: make(chan interface{}, 1),
	}
}

func (r *baseReceiver) receive(v interface{}) {
	r.result <- v
}

type distinctRequest struct {
	*baseReceiver
}

type topQueriesRequest struct {
	*baseReceiver
	size int
}
