package bplustree

import "testing"

func TestInsertarYBuscar(t *testing.T) {
	arbol := New(4)

	for i := 1; i <= 20; i++ {
		arbol.Insertar(Registro{Indice: i, TrackName: "Cancion", Artists: "Artista"})
	}

	for i := 1; i <= 20; i++ {
		reg, ok := arbol.Buscar(i)
		if !ok || reg.Indice != i {
			t.Fatalf("no se encontró índice %d", i)
		}
	}

	if _, ok := arbol.Buscar(99); ok {
		t.Fatal("índice inexistente no debe encontrarse")
	}
}

func TestBuscarRango(t *testing.T) {
	arbol := New(4)

	for i := 1; i <= 15; i++ {
		arbol.Insertar(Registro{Indice: i * 10, TrackName: "Cancion"})
	}

	rango := arbol.BuscarRango(30, 80)
	if len(rango) != 6 {
		t.Fatalf("esperaba 6 resultados, obtuvo %d", len(rango))
	}
	if rango[0].Indice != 30 || rango[len(rango)-1].Indice != 80 {
		t.Fatalf("rango incorrecto: %+v", rango)
	}
}

func TestBusquedaBajaHastaHoja(t *testing.T) {
	arbol := New(2)

	for i := 1; i <= 5; i++ {
		arbol.Insertar(Registro{Indice: i, TrackName: "C"})
	}

	reg, ok := arbol.Buscar(4)
	if !ok || reg.Indice != 4 {
		t.Fatal("búsqueda exacta debe leer el dato desde la hoja, no del índice")
	}

	hoja := arbol.irAHoja(4)
	if hoja.esHoja != true {
		t.Fatal("irAHoja(4) debe terminar en una hoja aunque 4 sea separador interno")
	}
}

func TestSequenceSet(t *testing.T) {
	arbol := New(3)
	total := 20

	for i := 1; i <= total; i++ {
		arbol.Insertar(Registro{Indice: i, TrackName: "C"})
	}

	recorrido := recorrerSequenceSet(arbol)
	if len(recorrido) != total {
		t.Fatalf("sequence set tiene %d registros, esperaba %d", len(recorrido), total)
	}

	for i, indice := range recorrido {
		if indice != i+1 {
			t.Fatalf("orden incorrecto en posición %d: obtuvo %d", i, indice)
		}
	}

	rango := arbol.BuscarRango(1, total)
	if len(rango) != total {
		t.Fatalf("BuscarRango devolvió %d registros, esperaba %d", len(rango), total)
	}
	for i, reg := range rango {
		if reg.Indice != recorrido[i] {
			t.Fatalf("BuscarRango y sequence set difieren en posición %d", i)
		}
	}
}
