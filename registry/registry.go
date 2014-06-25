package registry

import (
	"bytes"
	"errors"
	"sync"
)

var (
	StructNotFound = errors.New("Struct not found")
)

type Hashable interface {
	Hash() []byte
}

type Registry struct {
	nameToStruct map[string]*registeredStruct
	registry     []*registeredStruct
}

type registeredStruct struct {
	hashableStruct *Hashable
	m              *sync.Mutex
	name           string
	previousHash   []byte
}

func NewRegistry() *Registry {
	return &Registry{
		make(map[string]*registeredStruct),
		make([]*registeredStruct, 0),
	}
}

func (r *Registry) RegisterStruct(name string, h *Hashable) {
	rs := &registeredStruct{
		hashableStruct: h,
		m:              &sync.Mutex{},
		name:           name,
		previousHash:   (*h).Hash(),
	}
	r.registry = append(r.registry, rs)
	r.nameToStruct[name] = rs
}

func (r *Registry) ChangedStructs() []*Hashable {
	changed := make([]*Hashable, 0, len(r.registry))
	for _, rs := range r.registry {
		rs.m.Lock()
		hash := (*rs.hashableStruct).Hash()
		rs.m.Unlock()
		if !bytes.Equal(hash, rs.previousHash) {
			changed = append(changed, rs.hashableStruct)
			rs.previousHash = hash
		}
	}
	return changed
}

func (r *Registry) Checkout(name string) (*Hashable, *sync.Mutex, error) {
	rs, ok := r.nameToStruct[name]
	if !ok {
		return nil, nil, StructNotFound
	}
	return rs.hashableStruct, rs.m, nil
}
