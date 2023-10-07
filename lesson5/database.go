package main

var vclock = make(map[string]uint64)

var snap = "{}"

type Transaction struct {
	Source  string
	Id      uint64
	Payload string
}

var journal = make([]Transaction, 0)

var mySource = "robozmey"

var peers []string
