package broker

import "sync"

// Broker Broker
type Broker struct {
	C  chan bool
	wg *sync.WaitGroup
}

// New new
func New() *Broker {
	b := &Broker{
		C:  make(chan bool),
		wg: &sync.WaitGroup{},
	}

	b.wg.Add(1)
	return b
}

// Done Done
func (b *Broker) Done() {
	b.wg.Done()
}

// Cancel Cancel
func (b *Broker) Cancel() {
	b.C <- true
	b.wg.Wait()
}
