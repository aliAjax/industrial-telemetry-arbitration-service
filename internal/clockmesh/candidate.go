package clockmesh

type Candidate struct {
	ID     string
	Offset int64
	Labels map[string]string
}

func (c Candidate) Clone() Candidate {
	clone := c
	clone.Labels = make(map[string]string, len(c.Labels))
	for key, value := range c.Labels {
		clone.Labels[key] = value
	}
	return clone
}
