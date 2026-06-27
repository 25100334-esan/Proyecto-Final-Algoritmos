package bplustree

import (
	"fmt"
	"testing"
)

var tamanosBenchmark = []int{1_000, 10_000, 100_000, 1_000_000}

func BenchmarkInsertar(b *testing.B) {
	for _, n := range tamanosBenchmark {
		b.Run(fmt.Sprintf("N=%d", n), func(b *testing.B) {
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				arbol := New(128)
				for j := 1; j <= n; j++ {
					arbol.Insertar(Registro{Indice: j, TrackName: "Cancion"})
				}
			}
		})
	}
}

func BenchmarkBuscar(b *testing.B) {
	for _, n := range tamanosBenchmark {
		arbol := New(128)
		for j := 1; j <= n; j++ {
			arbol.Insertar(Registro{Indice: j, TrackName: "Cancion"})
		}

		b.Run(fmt.Sprintf("N=%d", n), func(b *testing.B) {
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				arbol.Buscar((i % n) + 1)
			}
		})
	}
}

func BenchmarkBuscarRango(b *testing.B) {
	const tamRango = 50

	for _, n := range tamanosBenchmark {
		if n <= tamRango {
			continue
		}
		arbol := New(128)
		for j := 1; j <= n; j++ {
			arbol.Insertar(Registro{Indice: j, TrackName: "Cancion"})
		}

		b.Run(fmt.Sprintf("N=%d", n), func(b *testing.B) {
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				desde := (i % (n - tamRango)) + 1
				arbol.BuscarRango(desde, desde+tamRango-1)
			}
		})
	}
}
