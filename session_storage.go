package tbot

import "sync"

type SessionStorage interface {
	Set(int64, string)
	Get(int64) string
	Reset(int64)
}

type InMemoryStorage struct {
	sync.Mutex
	sessions map[int64]string
}

func NewSessionStorage() SessionStorage {
	return &InMemoryStorage{sessions: make(map[int64]string)}
}

func (ims *InMemoryStorage) Get(id int64) string {
	ims.Lock()
	path := ims.sessions[id]
	ims.Unlock()
	return path
}

func (ims *InMemoryStorage) Set(id int64, path string) {
	ims.Lock()
	ims.sessions[id] = path
	ims.Unlock()
}

func (ims *InMemoryStorage) Reset(id int64) {
	ims.Lock()
	delete(ims.sessions, id)
	ims.Unlock()
}
