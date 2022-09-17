package main

import "sync"

type Cache struct {
	m    sync.RWMutex
	data map[int]string
}

func New() *Cache {
	c := &Cache{}
	c.data = make(map[int]string)
	return c
}

func main() {

}
