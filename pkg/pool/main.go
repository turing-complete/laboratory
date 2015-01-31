package pool

type Pool struct {
	collection chan interface{}
	create     func() interface{}
}

func New(capacity int, create func() interface{}) *Pool {
	return &Pool{
		collection: make(chan interface{}, capacity),
		create:     create,
	}
}

func (p *Pool) Get() interface{} {
	select {
	case item := <-p.collection:
		return item
	default:
		return p.create()
	}
}

func (p *Pool) Put(item interface{}) {
	select {
	case p.collection <- item:
	default:
	}
}
