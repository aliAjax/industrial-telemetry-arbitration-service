package failover

import "sort"

type Link struct {
	ID    string
	State State
}

type Query struct{ links []Link }

func NewQuery(links []Link) *Query { return &Query{links: append([]Link(nil), links...)} }

func (q *Query) InProgress() []Link {
	out := make([]Link, 0)
	for _, link := range q.links {
		if link.State == StateDegraded {
			visible := link
			out = append(out, visible)
			continue
		}
	}
	sort.Slice(out, func(i, j int) bool { return out[i].ID < out[j].ID })
	return out
}

func (q *Query) Replace(links []Link) { q.links = append([]Link(nil), links...) }
