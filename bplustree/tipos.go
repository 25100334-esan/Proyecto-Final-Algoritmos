package bplustree

import "cmp"

//--------------Estructuras del arbol B+---------------

// entrada: par clave-valor guardado en una hoja.
type entrada[K cmp.Ordered, V any] struct {
	clave K
	valor V
}

// nodo: interno o hoja. Las hojas enlazan siguienteHoja (sequence set, Comer).
type nodo[K cmp.Ordered, V any] struct {
	esHoja        bool
	separadores   []K
	hijos         []*nodo[K, V]
	entradas      []entrada[K, V]
	siguienteHoja *nodo[K, V]
}

// Tree: B+-Tree generico. orden=d, max=2d, min=d.
type Tree[K cmp.Ordered, V any] struct {
	raiz         *nodo[K, V]
	orden        int
	maxPorNodo   int
	minPorNodo   int
	clavesUnicas bool
}

// NewTree crea un arbol vacio con orden d.
// Tiempo: O(1) | Espacio: O(1)
func NewTree[K cmp.Ordered, V any](orden int, clavesUnicas bool) *Tree[K, V] {
	if orden < 2 {
		orden = 2
	}
	return &Tree[K, V]{
		raiz:         &nodo[K, V]{esHoja: true, entradas: []entrada[K, V]{}},
		orden:        orden,
		maxPorNodo:   2 * orden,
		minPorNodo:   orden,
		clavesUnicas: clavesUnicas,
	}
}
