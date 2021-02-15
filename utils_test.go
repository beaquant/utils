package utils

import (
	"strings"
	"testing"
)

func TestFloat64RoundString(t *testing.T) {
	t.Log(Float64RoundString(0.123456789))

	t.Log(Float64RoundString(0.12))
}

func TestFloat64RoundString2(t *testing.T) {
	t.Log(Float64RoundString2(0.123456789))
	t.Log(Float64RoundString2(0.12))
}

func TestFloat64RoundString3(t *testing.T) {
	t.Log(Float64RoundString3(0.123456789))
	t.Log(Float64RoundString3(0.12))
}

func TestFloatToString(t *testing.T) {
	t.Log(FloatToString(0.123456789, 4))
	t.Log(FloatToString(0.12, 4))
}

func TestFloat64ToString(t *testing.T) {
	t.Log(Float64ToString(0.123456789))
	t.Log(Float64ToString(0.12))
}

func TestFloat64Round(t *testing.T) {
	t.Log(Float64Round(0.123456789, 8))
}

func TestFloat64Round2(t *testing.T) {
	t.Log(Float64Round2(0.123456789))
}

func TestFloat64Round3(t *testing.T) {
	t.Log(Float64Round3(0.123456789))
}

func BenchmarkFloat64RoundString(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Float64RoundString(0.123456789)
	}
}

func BenchmarkFloat64RoundString2(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Float64RoundString2(0.123456789)
	}
}

func BenchmarkFloat64ToString(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Float64ToString(0.123456789)
	}
}

func BenchmarkFloat64ToString2(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Float64ToString2(0.123456789)
	}
}

func BenchmarkFloat64Round(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Float64Round(0.123456789)
	}
}

func BenchmarkFloat64Round2(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Float64Round2(0.123456789)
	}
}

func BenchmarkFloat64Round3(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Float64Round3(0.123456789, 4)
	}
}

// 117 ns/op
func BenchmarkGuuid(b *testing.B) {
	b.ResetTimer()
	s := "0.123456789"
	n := 2
	for i := 0; i < b.N; i++ {

		s1 := strings.Split(s, ".")
		amt := s1[0] + "."
		if len(s1) == 2 && len(s1[1]) > n {
			amt += s1[1][:n]
		}
	}
}

// 285 ns/op
func BenchmarkGuuid2(b *testing.B) {
	b.ResetTimer()
	s := "0.123456789"
	n := 2
	for i := 0; i < b.N; i++ {
		f := StringToFloat64(s)
		Float64RoundString(f, n)
	}
}
