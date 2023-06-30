package clients

import (
	//"errors"
	"container/list"
	//"log"
	"sync"
	//"time"
)

type Client struct {
	sync.Mutex
	id   string
	name string
}

type SourceList struct {
	sync.Mutex
	SList list.List
}

type Source struct {
	sync.Mutex
	Id         string `json:"id"`
	Name       string `json:"username"`
	Localtion  string `json:"localtion"`
	Categorize string `json:"categorize"`
}

// func NewSource() *Source {
// 	log.SetFlags(log.Ldate | log.Lshortfile)
// 	return &Source{}
// }

// func (s *Source) GetId() string {
// 	return s.id
// }
// func (s *Source) SetId(string id) {
// 	s.id = id
// }

// func (s *Source) GetName() string {
// 	return s.name
// }
// func (s *Source) SetName(string name) {
// 	s.name = name
// }
