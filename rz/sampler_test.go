package rz

import (
	"testing"
	"time"
)

var samplers = []struct {
	name    string
	sampler func() LogSampler
	total   int
	wantMin int
	wantMax int
}{
	{
		"SamplerBasic_1",
		func() LogSampler {
			return &SamplerBasic{N: 1}
		},
		100, 100, 100,
	},
	{
		"SamplerBasic_5",
		func() LogSampler {
			return &SamplerBasic{N: 5}
		},
		100, 20, 20,
	},
	{
		"SamplerRandom",
		func() LogSampler {
			return SamplerRandom(5)
		},
		100, 10, 30,
	},
	{
		"SamplerBurst",
		func() LogSampler {
			return &SamplerBurst{Burst: 20, Period: time.Second}
		},
		100, 20, 20,
	},
	{
		"SamplerBurstNext",
		func() LogSampler {
			return &SamplerBurst{Burst: 20, Period: time.Second, NextSampler: &SamplerBasic{N: 5}}
		},
		120, 40, 40,
	},
}

func TestSamplers(t *testing.T) {
	for i := range samplers {
		s := samplers[i]
		t.Run(s.name, func(t *testing.T) {
			sampler := s.sampler()
			got := 0
			for t := s.total; t > 0; t-- {
				if sampler.Sample(0) {
					got++
				}
			}
			if got < s.wantMin || got > s.wantMax {
				t.Errorf("%s.Sample(0) == true %d on %d, want [%d, %d]", s.name, got, s.total, s.wantMin, s.wantMax)
			}
		})
	}
}

func BenchmarkSamplers(b *testing.B) {
	for i := range samplers {
		s := samplers[i]
		b.Run(s.name, func(b *testing.B) {
			sampler := s.sampler()
			b.RunParallel(func(pb *testing.PB) {
				for pb.Next() {
					sampler.Sample(0)
				}
			})
		})
	}
}
