package utils

import "testing"

func TestFloat64RoundString(t *testing.T) {
	t.Log(Float64RoundString(0.123456789))
}

func TestFloat64RoundString2(t *testing.T) {
	t.Log(Float64RoundString2(0.123456789))
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
