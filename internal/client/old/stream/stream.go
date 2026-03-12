package stream

type Stream struct {
	ID     uint8
	Domain string
}

func New(domain string, id uint8) *Stream {
	return &Stream{
		ID:     id,
		Domain: domain,
	}
}
