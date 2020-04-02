package cachehash

type HashStore interface {
	Hash(string) (string, error)
	Lock(string) (func() error, error)
	LoadHashData(hash string) (*HashData, error)
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
	Changed      bool
}

func NewHashDataStatus() *HashDataStatus {
	return &HashDataStatus{
		FirstExpired: -1,
		LastExpired:  -1,
		Changed:      false,
	}
}

type Hash []*HashData

func (h *Hash) set(data *HashData) {
	for k := range *h {
		if (*h)[k].Key == data.Key {
			(*h)[k] = data
		}
	}
	*h = append(*h, data)
}

func (h *Hash) get(key string) *HashData {
	for k := range *h {
		if (*h)[k].Key == key {
			return (*h)[k]
		}
	}
	return nil
}
