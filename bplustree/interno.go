package bplustree

import "cmp"

//--------------Recorrido interno del arbol---------------
// Helpers privados: bajan desde la raiz hasta una hoja.
// Los usan buscar.go, insertar.go y eliminar.go.

func compare[K cmp.Ordered](a, b K) int {
	return cmp.Compare(a, b)
}

func noMenor[K cmp.Ordered](a, b K) bool {
	return compare(a, b) >= 0
}

func menor[K cmp.Ordered](a, b K) bool {
	return compare(a, b) < 0
}

func mayor[K cmp.Ordered](a, b K) bool {
	return compare(a, b) > 0
}

func posicionEntrada[K cmp.Ordered, V any](entradas []entrada[K, V], clave K) int {
	for i, e := range entradas {
		if !menor(e.clave, clave) {
			return i
		}
	}
	return len(entradas)
}

func (t *Tree[K, V]) indiceHijo(n *nodo[K, V], clave K) int {
	i := 0
	if t.clavesUnicas {
		for i < len(n.separadores) && noMenor(clave, n.separadores[i]) {
			i++
		}
	} else {
		for i < len(n.separadores) && mayor(clave, n.separadores[i]) {
			i++
		}
	}
	return i
}

// irAHoja: baja a la hoja candidata. Tiempo O(log_d n), espacio O(1).
func (t *Tree[K, V]) irAHoja(clave K) *nodo[K, V] {
	n := t.raiz
	for !n.esHoja {
		n = n.hijos[t.indiceHijo(n, clave)]
	}
	return n
}

func (t *Tree[K, V]) bajarGuardandoRuta(clave K) (*nodo[K, V], []*nodo[K, V], []int) {
	padres := []*nodo[K, V]{}
	posiciones := []int{}
	n := t.raiz

	for !n.esHoja {
		i := t.indiceHijo(n, clave)
		padres = append(padres, n)
		posiciones = append(posiciones, i)
		n = n.hijos[i]
	}
	return n, padres, posiciones
}

func (t *Tree[K, V]) rutaAHoja(target *nodo[K, V]) ([]*nodo[K, V], []int) {
	padres := []*nodo[K, V]{}
	posiciones := []int{}
	var dfs func(*nodo[K, V]) bool
	dfs = func(n *nodo[K, V]) bool {
		if n == target {
			return true
		}
		if n.esHoja {
			return false
		}
		for i, h := range n.hijos {
			if dfs(h) {
				padres = append(padres, n)
				posiciones = append(posiciones, i)
				return true
			}
		}
		return false
	}
	if !dfs(t.raiz) {
		return nil, nil
	}
	for i, j := 0, len(padres)-1; i < j; i, j = i+1, j-1 {
		padres[i], padres[j] = padres[j], padres[i]
		posiciones[i], posiciones[j] = posiciones[j], posiciones[i]
	}
	return padres, posiciones
}

func (t *Tree[K, V]) primeraHoja() *nodo[K, V] {
	n := t.raiz
	for n != nil && !n.esHoja {
		n = n.hijos[0]
	}
	return n
}
