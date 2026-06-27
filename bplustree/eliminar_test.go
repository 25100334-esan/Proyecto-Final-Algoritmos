package bplustree

import (
	"fmt"
	"testing"
)

func TestEliminar(t *testing.T) {
	arbol := New(4)

	for i := 1; i <= 10; i++ {
		arbol.Insertar(Registro{Indice: i, TrackName: "Cancion"})
	}

	for i := 1; i <= 10; i++ {
		arbol.Eliminar(i)
		if _, ok := arbol.Buscar(i); ok {
			t.Fatalf("índice %d no fue eliminado", i)
		}
	}
}

func TestSeparadorFantasma(t *testing.T) {
	arbol := New(2)

	for i := 1; i <= 6; i++ {
		arbol.Insertar(Registro{Indice: i, TrackName: fmt.Sprintf("C%d", i)})
	}

	if arbol.raiz.separadores[0] != 4 {
		t.Fatalf("separador promovido debe ser 4, obtuvo %d", arbol.raiz.separadores[0])
	}

	arbol.Eliminar(4)

	if arbol.raiz.separadores[0] != 4 {
		t.Fatal("el separador 4 debe permanecer como fantasma en el índice")
	}
	if _, ok := arbol.Buscar(4); ok {
		t.Fatal("el ID 4 ya no debe existir en ninguna hoja")
	}
	for _, id := range []int{3, 5, 6} {
		if _, ok := arbol.Buscar(id); !ok {
			t.Fatalf("búsqueda %d debe seguir funcionando con separador fantasma", id)
		}
	}
}

func TestUnderflowYFusion(t *testing.T) {
	arbol := New(3)

	for i := 1; i <= 12; i++ {
		arbol.Insertar(Registro{Indice: i, TrackName: fmt.Sprintf("C%d", i)})
	}

	for i := 12; i >= 1; i-- {
		arbol.Eliminar(i)
		for j := 1; j < i; j++ {
			if _, ok := arbol.Buscar(j); !ok {
				t.Fatalf("tras eliminar %d, el índice %d ya no se encuentra", i, j)
			}
		}
		if _, ok := arbol.Buscar(i); ok {
			t.Fatalf("índice %d no fue eliminado", i)
		}
	}

	if !arbol.raiz.esHoja {
		t.Fatal("al quedar vacío la raíz debe volver a ser hoja")
	}
	if len(arbol.raiz.entradas) != 0 {
		t.Fatalf("raíz debe estar vacía, tiene %d registros", len(arbol.raiz.entradas))
	}
}

func TestReducirAlturaRaiz(t *testing.T) {
	arbol := New(2)

	for i := 1; i <= 5; i++ {
		arbol.Insertar(Registro{Indice: i, TrackName: fmt.Sprintf("C%d", i)})
	}

	if arbol.raiz.esHoja {
		t.Fatal("con 5 registros (d=2) la raíz debe ser interna")
	}
	alturaAntes := alturaArbol(arbol.raiz)

	for i := 1; i <= 5; i++ {
		arbol.Eliminar(i)
	}

	if !arbol.raiz.esHoja {
		t.Fatal("al vaciar el árbol la raíz debe volver a ser hoja")
	}
	if len(arbol.raiz.entradas) != 0 {
		t.Fatalf("raíz vacía debe tener 0 canciones, tiene %d", len(arbol.raiz.entradas))
	}
	if alturaAntes < 2 {
		t.Fatalf("se esperaba altura >= 2 antes de eliminar, obtuvo %d", alturaAntes)
	}
}
