package top

import "sync"

type MultiChanObserverInterface interface {
	AddChannel(ch chan *ProgramPool, f func(*ProgramPool))
	RemoveChannel(ch chan *ProgramPool)
}
type MultiChanObserver struct {
	MultiChanObserverInterface
	events chan *Program
	chMap  map[string]chan *Program
	mutex  sync.Mutex
}

func NewMultiChanObserver() *MultiChanObserver {
	obj := MultiChanObserver{chMap: make(map[string]chan *Program),
		events: make(chan *Program, 10),
	}
	return &obj
}

func (o *MultiChanObserver) onPoolChanged() {
	for program := range o.events {
		o.mutex.Lock()
		for _, ch := range o.chMap {
			select {
			case ch <- program:
			default:
				// Если канал заполнен, игнорируем значение и переходим к следующему каналу
				continue
			}
		}
		o.mutex.Unlock()
	}
}

func (o *MultiChanObserver) Notify(pp *Program) {
	o.mutex.Lock()
	defer o.mutex.Unlock()
	o.events <- pp
}

func (o *MultiChanObserver) AddChannel(name string) *chan *Program {
	o.mutex.Lock()
	defer o.mutex.Unlock()
	if len(o.chMap) == 0 {
		go o.onPoolChanged()
	}
	ch := make(chan *Program, 5)
	o.chMap[name] = ch

	return &ch
}

func (o *MultiChanObserver) RemoveChannel(name string) {
	o.mutex.Lock()
	defer o.mutex.Unlock()
	delete(o.chMap, name)
}
