package chord

type successors struct {
	data []string
}

func createSuccessorQueue() *successors {
	return &successors{data: make([]string, 0)}
}

func (q *successors) push(address string) {
	q.data = append(q.data, address)
}

func (q *successors) pop() string {
	d := q.data[0]
	q.data = q.data[1:]
	return d
}

func (q *successors) get(i int) string {
	return q.data[i]
}

func (q *successors) length() int {
	return len(q.data)
}

func (q *successors) empty() bool {
	return len(q.data) == 0
}
