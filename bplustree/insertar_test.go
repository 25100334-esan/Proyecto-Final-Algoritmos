package bplustree

import (
	"fmt"
	"testing"
)

func TestOrdenComer(t *testing.T) {
	arbol := New(2)
	for i := 1; i <= 4; i++ {
		arbol.Insertar(Registro{Indice: i, TrackName: "C"})
	}
	if !arbol.raiz.esHoja || len(arbol.raiz.entradas) != 4 {
		t.Fatal("con 4 claves (2d) no debe haber split")
	}
	arbol.Insertar(Registro{Indice: 5, TrackName: "C"})
	if arbol.raiz.esHoja {
		t.Fatal("la 5.ª inserción debe provocar split (overflow > 2d)")
	}

	arbol4 := New(4)
	for i := 1; i <= 8; i++ {
		arbol4.Insertar(Registro{Indice: i, TrackName: "C"})
	}
	if !arbol4.raiz.esHoja || len(arbol4.raiz.entradas) != 8 {
		t.Fatalf("d=4: esperaba 8 claves sin split, obtuvo %d", len(arbol4.raiz.entradas))
	}
}

func TestInsertarDuplicado(t *testing.T) {
	arbol := New(4)
	arbol.Insertar(Registro{Indice: 5, TrackName: "A"})
	arbol.Insertar(Registro{Indice: 5, TrackName: "B"})

	reg, ok := arbol.Buscar(5)
	if !ok || reg.TrackName != "A" {
		t.Fatal("insertar duplicado no debe sobrescribir")
	}
}

func TestOverflowYSplit(t *testing.T) {
	arbol := New(2)

	for i := 1; i <= 5; i++ {
		arbol.Insertar(Registro{Indice: i, TrackName: fmt.Sprintf("C%d", i)})
	}

	if arbol.raiz.esHoja {
		t.Fatal("después del split la raíz debe ser un nodo interno")
	}
	if len(arbol.raiz.separadores) != 1 || arbol.raiz.separadores[0] != 4 {
		t.Fatalf("separador promovido incorrecto: %v", arbol.raiz.separadores)
	}
	if len(arbol.raiz.hijos) != 2 {
		t.Fatalf("raíz debe tener 2 hijos, tiene %d", len(arbol.raiz.hijos))
	}

	izq := arbol.raiz.hijos[0]
	der := arbol.raiz.hijos[1]
	if !izq.esHoja || !der.esHoja {
		t.Fatal("los hijos de la raíz deben ser hojas")
	}
	if len(izq.entradas) != 3 || izq.entradas[0].clave != 1 || izq.entradas[2].clave != 3 {
		t.Fatalf("hoja izquierda incorrecta: %+v", indicesDeHoja(izq))
	}
	if len(der.entradas) != 2 || der.entradas[0].clave != 4 || der.entradas[1].clave != 5 {
		t.Fatalf("hoja derecha incorrecta: %+v", indicesDeHoja(der))
	}

	if _, ok := arbol.Buscar(4); !ok {
		t.Fatal("el índice promovido debe seguir en la hoja (copia, no movimiento)")
	}
	if izq.siguienteHoja != der {
		t.Fatal("las hojas deben estar enlazadas horizontalmente tras el split")
	}

	for i := 6; i <= 9; i++ {
		arbol.Insertar(Registro{Indice: i, TrackName: fmt.Sprintf("C%d", i)})
	}
	if contarHojas(arbol.raiz) < 2 {
		t.Fatal("con 9 registros debe haber al menos 2 hojas")
	}
	for i := 1; i <= 9; i++ {
		if _, ok := arbol.Buscar(i); !ok {
			t.Fatalf("índice %d perdido tras múltiples splits", i)
		}
	}
}
