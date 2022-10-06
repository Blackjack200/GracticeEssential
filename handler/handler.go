package handler

import (
	"github.com/df-mc/dragonfly/server/player"
	"sync"
)

type MultipleHandler struct {
	mu *sync.Mutex
	//handler -> nil
	handlers map[any]any
}

func (h *MultipleHandler) Register(hdr interface{}) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.handlers[hdr] = nil
}

func (h *MultipleHandler) Unregister(hdr interface{}) {
	h.mu.Lock()
	defer h.mu.Unlock()
	delete(h.handlers, hdr)
}

func newHandler() *MultipleHandler {
	return &MultipleHandler{
		mu:       &sync.Mutex{},
		handlers: make(map[any]any),
	}
}

func Handle(p *player.Player) *MultipleHandler {
	hdr := newHandler()
	p.Handle(hdr)
	return hdr
}
