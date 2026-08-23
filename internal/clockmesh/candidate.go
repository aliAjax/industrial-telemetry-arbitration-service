package clockmesh

type Candidate struct {
	ID     string
	Offset int64
	Labels map[string]string
}

func (c Candidate) Clone() Candidate {
	return c
}
