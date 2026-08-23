package waveform

type Assembler struct {
	samples []Sample
	limit   int
}

func NewAssembler(limit int) *Assembler { return &Assembler{limit: limit} }

func (a *Assembler) Append(next []Sample) []Sample {
	combined := make([]Sample, 0, len(a.samples)+len(next))
	combined = append(combined, a.samples...)
	combined = append(combined, next...)
	if a.limit > 0 && len(combined) > a.limit {
		combined = Clone(combined[len(combined)-a.limit:])
	}
	a.samples = Clone(combined)
	return Clone(a.samples)
}

func (a *Assembler) Samples() []Sample { return Clone(a.samples) }
