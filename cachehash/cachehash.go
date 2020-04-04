package cachehash

type Store interface {
	Open() error
	Cloas() error
	Flush() error
	Hash(string) (string, error)
	Lock(string) (func() error, error)
	LoadHashData(hash string) (*HashData, error)
	DeleteHashData(hash string) error
	SaveHashData(hash string, firstExpired int64, lastExpired int64, data *HashData) error
}
type HashData struct {
	Key     string
	Expired int64
	Data    []byte
}
type HashDataStatus struct {
	FirstExpired int64
	LastExpired  int64
	Size         int
	Changed      bool
}

func (s *HashDataStatus) calc(data *HashData, current int64) bool {
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
func NewHashDataStatus() *HashDataStatus {
	return &HashDataStatus{
		FirstExpired: 0,
		LastExpired:  0,
		Size:         0,
		Changed:      false,
	}
}

type Hash []*HashData

func (h *Hash) set(data *HashData, current int64) *HashDataStatus {
	result := make(Hash, 0, len(*h)+1)
	status := NewHashDataStatus()
	status.Changed = true
	for k := range *h {
		if (*h)[k].Key != data.Key || status.calc((*h)[k], current) {
			result = append(result, (*h)[k])
		}
	}
	if status.calc(data, current) {
		*h = append(*h, data)
	}
	h = &result
	return status
}
func (h *Hash) delete(key string, current int64) *HashDataStatus {
	result := make(Hash, 0, len(*h))
	status := NewHashDataStatus()
	for k := range *h {
		if (*h)[k].Key != key || status.calc((*h)[k], current) {
			result = append(result, (*h)[k])
		}
	}
	h = &result
	return status
}
func (h *Hash) scan(current int64) *HashDataStatus {
	result := make(Hash, 0, len(*h))
	status := NewHashDataStatus()
	for k := range *h {
		if status.calc((*h)[k], current) {
			result = append(result, (*h)[k])
		}
	}
	h = &result
	return status
}

func (h *Hash) get(key string) *HashData {
	for k := range *h {
		if (*h)[k].Key == key {
			return (*h)[k]
		}
	}
	return nil
}
