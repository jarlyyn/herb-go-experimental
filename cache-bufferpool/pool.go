package bufferpool

func New(defaultBufferSize int, maxBufferSize int) *Pool {
	p := &Pool{
		Buffers:           make(chan *Buffer, 40),
		DefaultBufferSize: defaultBufferSize,
		MaxBufferSize:     maxBufferSize,
	}
	for i := 0; i < 40; i++ {
		p.Buffers <- &Buffer{
			Pool:  p,
			Bytes: make([]byte, defaultBufferSize),
		}
	}
	return p
}

type Pool struct {
	Buffers           chan *Buffer
	DefaultBufferSize int
	MaxBufferSize     int
}

func (p *Pool) GetBuffer() *Buffer {
	return <-p.Buffers
}
