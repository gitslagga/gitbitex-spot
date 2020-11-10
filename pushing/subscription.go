package pushing

import (
	"sync"
)

type subscription struct {
	subscribers map[string]map[int64]*Client
	mu          sync.RWMutex
}

func newSubscription() *subscription {
	return &subscription{subscribers: map[string]map[int64]*Client{}}
}

func (s *subscription) subscribe(channel string, client *Client) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	_, found := s.subscribers[channel]
	if !found {
		s.subscribers[channel] = map[int64]*Client{}
	}

	_, found = s.subscribers[channel][client.id]
	if found {
		return false
	}
	s.subscribers[channel][client.id] = client
	return true
}

func (s *subscription) unsubscribe(channel string, client *Client) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	_, found := s.subscribers[channel]
	if !found {
		return false
	}

	_, found = s.subscribers[channel][client.id]
	if !found {
		return false
	}
	delete(s.subscribers[channel], client.id)
	return true
}

func (s *subscription) publish(channel string, msg interface{}) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	_, found := s.subscribers[channel]
	if !found {
		return
	}

	for _, c := range s.subscribers[channel] {
		c.writeCh <- msg
	}
}
