package cachehash

type Store interface {
	Open() error
	Close() error
	Flush() error
	Hash(string) (string, error)
	Lock(string) (func(), error)
	Load(hash string) (*Hash, error)
	Delete(hash string) error
	Save(hash string, status *Status, data *Hash) error
}
type Data struct {
	Key     string
	Expired int64
	Data    []byte
}

func NewData(key string, expired int64, data []byte) *Data {
	return &Data{
		Key:     key,
		Expired: expired,
		Data:    data,
	}
}

type Status struct {
	FirstExpired int64
	LastExpired  int64
	Size         int
	Delta        int
	Changed      bool
}

func (s *Status) calc(data *Data, current int64) bool {
	if data.Expired < current {
		s.Changed = true
		return false
	}
	if s.LastExpired <= 0 || s.FirstExpired > data.Expired {
		s.FirstExpired = data.Expired
	}
	if s.LastExpired < data.Expired {
		s.LastExpired = data.Expired
	}
	s.Size = s.Size + len(data.Data)
	return true
}
func NewStatus() *Status {
	return &Status{
		FirstExpired: 0,
		LastExpired:  0,
		Size:         0,
		Changed:      false,
	}
}

type Hash []*Data

func (h *Hash) isEmpty() bool {
	return len(*h) == 0
}
func (h *Hash) set(data *Data, current int64) *Status {
	result := make(Hash, 0, len(*h)+1)
	status := NewStatus()
	status.Changed = true
	for k := range *h {
		if (*h)[k].Key != data.Key || status.calc((*h)[k], current) {
			result = append(result, (*h)[k])
			status.Size = status.Size + len((*h)[k].Data)
		} else {
			status.Delta = status.Delta - len((*h)[k].Data)
		}
	}
	if status.calc(data, current) {
		*h = append(*h, data)
		status.Size = status.Size + len(data.Data)
		status.Delta = status.Delta + len(data.Data)
	}
	h = &result
	return status
}
func (h *Hash) update(data *Data, current int64) *Status {
	result := make(Hash, 0, len(*h)+1)
	status := NewStatus()
	status.Changed = true

	for k := range *h {
		var delta int
		d := (*h)[k]
		if d.Key == data.Key {
			d.Expired = data.Expired
			d.Data = data.Data
			delta = len(data.Data) - len(d.Data)
		}
		if status.calc((*h)[k], current) {
			result = append(result, (*h)[k])
			status.Size = status.Size + len(data.Data)
			status.Delta = status.Delta + delta
		} else {
			status.Delta = status.Delta - len(d.Data)
		}
	}
	h = &result
	return status
}
func (h *Hash) expired(key string, expired int64, current int64) *Status {
	result := make(Hash, 0, len(*h)+1)
	status := NewStatus()
	status.Changed = true
	for k := range *h {
		d := (*h)[k]
		if d.Key == key {
			d.Expired = expired
		}
		if status.calc((*h)[k], current) {
			result = append(result, (*h)[k])
			status.Size = status.Size + len((*h)[k].Data)
		} else {
			status.Delta = status.Delta - len((*h)[k].Data)
		}
	}
	h = &result
	return status
}
func (h *Hash) delete(key string, current int64) *Status {
	result := make(Hash, 0, len(*h))
	status := NewStatus()
	for k := range *h {
		if (*h)[k].Key != key || status.calc((*h)[k], current) {
			result = append(result, (*h)[k])
			status.Size = status.Size + len((*h)[k].Data)
		} else {
			status.Delta = status.Delta - len((*h)[k].Data)
		}
	}
	h = &result
	return status
}
func (h *Hash) scan(current int64) *Status {
	result := make(Hash, 0, len(*h))
	status := NewStatus()
	for k := range *h {
		if status.calc((*h)[k], current) {
			result = append(result, (*h)[k])
			status.Size = status.Size + len((*h)[k].Data)
		} else {
			status.Delta = status.Delta - len((*h)[k].Data)
		}
	}
	h = &result
	return status
}

func (h *Hash) get(key string, current int64) *Data {
	for k := range *h {
		if (*h)[k].Key == key {
			if (*h)[k].Expired >= current {
				return (*h)[k]
			}
			return (*h)[k]
		}
	}
	return nil
}
