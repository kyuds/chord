package chord

type successors struct {
	ip   string
	data []string
}

func createSuccessorQueue(ownIP string) *successors {
	return &successors{
		ip:   ownIP,
		data: make([]string, 0),
	}
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

func (q *successors) getSuccessor() string {
	if q.empty() {
		return q.ip
	}
	return q.get(0)
}
