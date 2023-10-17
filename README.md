# go-chord
Implementation of [Chord (DHash)](https://pdos.csail.mit.edu/papers/chord:sigcomm01/chord_sigcomm.pdf) in GoLang.

Currently a work in progress. 

This implementation stays true to the designs of the original paper, and 
presents the single operation, `lookup`, to get the address of the node that 
is responsible for handling that specific key. Unlike the original paper that 
presents Chord in the context of a distributed hash table (dht), this specific 
project deals with the barest usage of Chord: a mapping between keys and nodes. 

All of the relevant APIs are in `chord/chord.go` which includes Chord configurations, 
initialization, lookups, and planned exits. Please try out the `cli.go` file to see 
a live version of a Chord ring in action. 