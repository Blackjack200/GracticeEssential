package mhandler

import (
	"fmt"
	"golang.org/x/exp/slices"
	"reflect"
	"sync"
)

type MultipleHandler struct {
	mu       *sync.Mutex
	handlers []any
}

func (h *MultipleHandler) Register(hdr any) (unregister func()) {
	r := reflect.ValueOf(hdr)
	if r.NumMethod() < 1 {
		panic(fmt.Errorf("invalid mhandler: %v", hdr))
	}
	h.mu.Lock()
	defer h.mu.Unlock()
	h.handlers = append(h.handlers, hdr)
	unregister = func() { h.Unregister(hdr) }
	return
}

func (h *MultipleHandler) Clear() {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.handlers = nil
}

func (h *MultipleHandler) Unregister(hdr any) {
	h.mu.Lock()
	defer h.mu.Unlock()
	if idx := slices.IndexFunc(h.handlers, func(v any) bool {
		return v == hdr
	}); idx != -1 {
		h.handlers = slices.Delete(h.handlers, idx, idx+1)
	}
}

func New() *MultipleHandler {
	return &MultipleHandler{
		mu:       &sync.Mutex{},
		handlers: nil,
	}
}
