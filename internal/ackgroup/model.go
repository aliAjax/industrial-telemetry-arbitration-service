package ackgroup

type Ack struct {
	GatewayID string
	Sequence  uint64
	Metadata  map[string]string
}

func (a Ack) Clone() Ack {
	clone := a
	clone.Metadata = make(map[string]string, len(a.Metadata))
	for key, value := range a.Metadata {
		clone.Metadata[key] = value
	}
	return clone
}
