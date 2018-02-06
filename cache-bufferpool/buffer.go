package bufferpool

type Buffer struct {
	Pool  *Pool
	Bytes []byte
}

func (b *Buffer) Close() {
	if len(b.Bytes) > b.Pool.MaxBufferSize {
		b.Bytes = make([]byte, b.Pool.DefaultBufferSize)
	}
	b.Pool.Buffers <- b
}
