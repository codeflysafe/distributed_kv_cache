package main

import "testing"

func BenchmarkLock1(b *testing.B) {
	for i := 0; i < b.N; i++ {
		lock1()
	}
}

func BenchmarkChan1(b *testing.B) {
	for i := 0; i < b.N; i++ {
		chan1()
	}
}
