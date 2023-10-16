package chord

type node struct {
	id string
	ip string
	cnf *Config
}

func newNode(conf *Config) (*node, error) {
	return nil, nil
}

func (nd *node) joinNode(address string) error {
	return nil
}

