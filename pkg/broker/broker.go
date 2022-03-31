package broker

import "sync"

type Broker struct {
	C  chan bool
	wg *sync.WaitGroup
}

func New() *Broker {
	b := &Broker{
		C:  make(chan bool),
		wg: &sync.WaitGroup{},
	}

	b.wg.Add(1)
	return b
}

func (b *Broker) Done() {
	b.wg.Done()
}

func (b *Broker) Cancel() {
	b.C <- true
	b.wg.Wait()
}
