# chord

Implementation of [Chord](https://pdos.csail.mit.edu/papers/chord:sigcomm01/chord_sigcomm.pdf) in Go.

This implementation stays true to the designs of the original paper, and 
presents the single operation, `lookup`, to get the address of the node that 
is responsible for handling that specific key. Unlike the original paper that 
presents Chord in the context of a distributed hash table (dht), this specific 
project deals with the barest usage of Chord: a mapping between keys and nodes. 

All of the relevant APIs are in `chord/chord.go` which includes Chord configurations, 
initialization, lookups, and planned exits. Please try out the `example.go` file to see 
a live version of a Chord ring in action. Fault tolerance is implemented and prevents 
the whole ring from crashing from a node crash. Key locations are also automatically 
reassigned after a short period of time.

Configurations:
```go
type Config struct {
	Address       string
	Hash          func() hash.Hash
	Joining       bool
	JoinAddress   string
	MaxRepli      int
	ServerOptions []grpc.ServerOption
	DialOptions   []grpc.DialOption
	Timeout       time.Duration
	MaxIdle       time.Duration
	Stabilization time.Duration
	FingerFix     time.Duration
}
```

Creating a Chord Ring with Default Settings:
```
./go-chord create --address="localhost:8000"
```
Joining a Chord Ring with Default Settings:
```
./go-chord join --address="localhost:8001" --join="localhost:8000"
```

### Version History
- [X] V3: Multiplexing for gRPC servers
- [X] V2: Fault tolerance and exit support
- [X] V1: Basic Chord joining
