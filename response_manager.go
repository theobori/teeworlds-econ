package teeworldsecon

import "sync"

type EconResponseManager struct {
	id         int
	payloadsCh map[int]chan string
	mu         sync.Mutex
}

func NewEconResponseManager() *EconResponseManager {
	return &EconResponseManager{
		id:         0,
		payloadsCh: make(map[int]chan string),
	}
}

func (erm *EconResponseManager) Add(payloadCh chan string) int {
	erm.mu.Lock()
	defer erm.mu.Unlock()

	id := erm.id
	erm.payloadsCh[erm.id] = payloadCh

	erm.id++

	return id
}

func (erm *EconResponseManager) Delete(id int) {
	erm.mu.Lock()
	defer erm.mu.Unlock()

	delete(erm.payloadsCh, id)
}

func (erm *EconResponseManager) Send(payload string) {
	erm.mu.Lock()
	defer erm.mu.Unlock()

	for id, ch := range erm.payloadsCh {
		select {
		case ch <- payload:
			Debug("Sent payload to channel with id %d", id)
		default:
			Debug("Channel with id %d is blocked or closed", id)
		}
	}
}
