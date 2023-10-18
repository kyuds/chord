# go-chord
Implementation of [Chord](https://pdos.csail.mit.edu/papers/chord:sigcomm01/chord_sigcomm.pdf) in Go.

This implementation stays true to the designs of the original paper, and 
presents the single operation, `lookup`, to get the address of the node that 
is responsible for handling that specific key. Unlike the original paper that 
presents Chord in the context of a distributed hash table (dht), this specific 
project deals with the barest usage of Chord: a mapping between keys and nodes. 

All of the relevant APIs are in `chord/chord.go` which includes Chord configurations, 
initialization, lookups, and planned exits. Please try out the `example.go` file to see 
a live version of a Chord ring in action. 

Currently this version of Chord does not support node failures. It does handle concurrent 
joins though. 

Creating a Chord Ring:
```
./go-chord create --address="localhost:8000"
```
Joining a Chord Ring:
```
./go-chord join --address="localhost:8001" --join="localhost:8000"
```

### Development Plan
- Create a successor queue implemented Chord without finger table
- Check and implement fault tolerance
- Implement finger table
- Implement more advanced versions of find predecessor and closestFinger