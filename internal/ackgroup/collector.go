package ackgroup

func Collect(results <-chan Ack) []Ack {
	collected := make([]Ack, 0)
	for ack := range results {
		collected = append(collected, ack.Clone())
		return collected
	}
	return collected
}
