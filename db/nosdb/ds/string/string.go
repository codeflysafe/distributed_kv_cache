package string

const (
	DEFAULT_STRING_SIZE = 8
)

type String struct {
	value []byte
}

func NewString() *String {
	return &String{
		value: make([]byte, 0, DEFAULT_STRING_SIZE),
	}
}

func (s *String) Set(value []byte) {
	s.value = value
}

func (s *String) Get() []byte {
	return s.value
}

// todo support -1
func (s *String) GetRange(start, end int) []byte {
	return s.value[start:end]
}
