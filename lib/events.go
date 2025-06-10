package main

import "sync"

type Event struct {
	Event  string
	Method func(wg *sync.WaitGroup, data *any)
}
