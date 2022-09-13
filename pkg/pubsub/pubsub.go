package pubsub

import (
	"sync"
	"time"
)

type (
	Subscriber chan interface{}
	Topic      func(v interface{}) bool
)

type Publisher struct {
	subscribers map[Subscriber]Topic
	buffer      int
	timeout     time.Duration
	m           sync.RWMutex
}

func NewPublisher(publishTimeout time.Duration, buffer int) *Publisher {
	return &Publisher{
		buffer:      buffer,
		timeout:     publishTimeout,
		subscribers: make(map[Subscriber]Topic),
	}
}

func (p *Publisher) Subscribe() Subscriber {
	return p.SubscribeTopic(nil)
}

func (p *Publisher) SubscribeTopic(topic Topic) Subscriber {
	p.m.Lock()
	defer p.m.Unlock()

	ch := make(Subscriber, p.buffer)
	p.subscribers[ch] = topic

	return ch
}

func (p *Publisher) Cancel(sub Subscriber) {
	p.m.Lock()
	defer p.m.Unlock()

	delete(p.subscribers, sub)
	close(sub)
}

func (p *Publisher) Close() {
	p.m.Lock()
	defer p.m.Unlock()

	for sub := range p.subscribers {
		delete(p.subscribers, sub)
		close(sub)
	}
}

func (p *Publisher) Publish(v interface{}) {
	p.m.Lock()
	defer p.m.Unlock()

	var wg sync.WaitGroup
	for sub, topic := range p.subscribers {
		wg.Add(1)
		go p.send(sub, topic, v, &wg)
	}

	wg.Wait()
}

func (p *Publisher) send(sub Subscriber, topic Topic, v interface{}, wg *sync.WaitGroup) {
	defer wg.Done()

	if topic != nil && !topic(v) {
		return
	}

	select {
	case sub <- v:
	case <-time.After(p.timeout):
	}
}
