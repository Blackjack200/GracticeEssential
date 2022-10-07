package mhandler

import (
	"fmt"
	"reflect"
	"sync"
)

type MultipleHandler struct {
	mu *sync.Mutex
	//mhandler -> nil
	handlers map[any]any
}

func (h *MultipleHandler) Register(hdr interface{}) (unregister func()) {
	r := reflect.ValueOf(hdr)
	if r.NumMethod() < 1 {
		panic(fmt.Errorf("invalid mhandler: %v", hdr))
	}
	h.mu.Lock()
	defer h.mu.Unlock()
	h.handlers[hdr] = nil
	unregister = func() { h.Unregister(hdr) }
	return
}

func (h *MultipleHandler) Unregister(hdr interface{}) {
	h.mu.Lock()
	defer h.mu.Unlock()
	delete(h.handlers, hdr)
}

func New() *MultipleHandler {
	return &MultipleHandler{
		mu:       &sync.Mutex{},
		handlers: make(map[any]any),
	}
}
