package clockmesh

type Candidate struct {
	ID     string
	Offset int64
	Labels map[string]string
}

func (c Candidate) Clone() Candidate {
	cp := c
	if c.Labels != nil {
		cp.Labels = make(map[string]string, len(c.Labels))
		for k, v := range c.Labels {
			cp.Labels[k] = v
		}
	}
	return cp
}
