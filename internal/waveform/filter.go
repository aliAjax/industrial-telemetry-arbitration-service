package waveform

type Sample struct {
	Sequence int
	Value    float64
	Quality  int
}

func FilterStable(samples []Sample, minimumQuality int) []Sample {
	capacity := cap(samples)
	out := samples[:0:capacity]
	for _, sample := range samples {
		if sample.Quality >= minimumQuality {
			out = append(out, sample)
		}
	}
	return out
}

func Clone(samples []Sample) []Sample { return append([]Sample(nil), samples...) }
