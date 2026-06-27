package bplustree

import "cmp"

//--------------Eliminacion (API publica)---------------

// Eliminar: quita la primera entrada con clave exacta.
// Tiempo: O(log_d n) | Espacio: O(h) altura del arbol
func (t *Tree[K, V]) Eliminar(clave K) {
	t.EliminarSi(clave, func(V) bool { return true })
}

// EliminarSi: borra la entrada que cumple clave + predicado (duplicados en secundarios).
// Tiempo: O(log_d n + k) | Espacio: O(h)
func (t *Tree[K, V]) EliminarSi(clave K, coincide func(V) bool) bool {
	hoja := t.irAHoja(clave)
	for actual := hoja; actual != nil; {
		pos := posicionEntrada(actual.entradas, clave)
		for i := pos; i < len(actual.entradas) && compare(actual.entradas[i].clave, clave) == 0; i++ {
			if coincide(actual.entradas[i].valor) {
				padres, posiciones := t.rutaAHoja(actual)
				actual.entradas = append(actual.entradas[:i], actual.entradas[i+1:]...)
				if actual == t.raiz || len(actual.entradas) >= t.minPorNodo {
					return true
				}
				t.corregirUnderflowHoja(actual, padres, posiciones)
				return true
			}
		}
		if len(actual.entradas) == 0 {
			break
		}
		ultima := actual.entradas[len(actual.entradas)-1].clave
		if compare(ultima, clave) != 0 {
			break
		}
		actual = actual.siguienteHoja
	}
	return false
}

// EliminarExacto: en secundarios, borra por clave + ID de cancion.
// Tiempo: O(log_d n) | Espacio: O(h)
func EliminarExacto[K cmp.Ordered](t *Tree[K, Registro], clave K, indice int) bool {
	return t.EliminarSi(clave, coincideIndice(indice))
}

//--------------Underflow: prestamo y fusion (interno)---------------

func (t *Tree[K, V]) ajustarRaiz(n *nodo[K, V]) {
	if n != t.raiz || n.esHoja {
		return
	}
	if len(n.separadores) == 0 && len(n.hijos) == 1 {
		t.raiz = n.hijos[0]
	}
}

func (t *Tree[K, V]) corregirUnderflowHoja(hoja *nodo[K, V], padres []*nodo[K, V], posiciones []int) {
	if len(hoja.entradas) >= t.minPorNodo || len(padres) == 0 {
		return
	}

	padre := padres[len(padres)-1]
	pos := posiciones[len(posiciones)-1]

	if pos+1 < len(padre.hijos) {
		hermanoDer := padre.hijos[pos+1]
		if len(hermanoDer.entradas) > t.minPorNodo {
			t.prestarDerecha(hoja, hermanoDer, padre, pos)
			return
		}
	}

	if pos > 0 {
		hermanoIzq := padre.hijos[pos-1]
		if len(hermanoIzq.entradas) > t.minPorNodo {
			t.prestarIzquierda(hoja, hermanoIzq, padre, pos-1)
			return
		}
	}

	if pos+1 < len(padre.hijos) {
		t.fusionarConDerecha(hoja, padre, pos)
	} else {
		t.fusionarConIzquierda(hoja, padre, pos-1)
	}

	t.ajustarRaiz(padre)
	if padre != t.raiz && len(padre.separadores) < t.minPorNodo {
		t.corregirUnderflowInterno(padre, padres[:len(padres)-1], posiciones[:len(posiciones)-1])
	}
}

func (t *Tree[K, V]) prestarDerecha(hoja, hermano *nodo[K, V], padre *nodo[K, V], pos int) {
	hoja.entradas = append(hoja.entradas, hermano.entradas[0])
	hermano.entradas = hermano.entradas[1:]
	padre.separadores[pos] = hermano.entradas[0].clave
}

func (t *Tree[K, V]) prestarIzquierda(hoja, hermano *nodo[K, V], padre *nodo[K, V], pos int) {
	hoja.entradas = append([]entrada[K, V]{hermano.entradas[len(hermano.entradas)-1]}, hoja.entradas...)
	hermano.entradas = hermano.entradas[:len(hermano.entradas)-1]
	padre.separadores[pos] = hoja.entradas[0].clave
}

func (t *Tree[K, V]) fusionarConDerecha(hoja *nodo[K, V], padre *nodo[K, V], pos int) {
	hermano := padre.hijos[pos+1]
	hoja.entradas = append(hoja.entradas, hermano.entradas...)
	hoja.siguienteHoja = hermano.siguienteHoja
	padre.separadores = append(padre.separadores[:pos], padre.separadores[pos+1:]...)
	padre.hijos = append(padre.hijos[:pos+1], padre.hijos[pos+2:]...)
}

func (t *Tree[K, V]) fusionarConIzquierda(hoja *nodo[K, V], padre *nodo[K, V], pos int) {
	hermano := padre.hijos[pos]
	hermano.entradas = append(hermano.entradas, hoja.entradas...)
	hermano.siguienteHoja = hoja.siguienteHoja
	padre.separadores = append(padre.separadores[:pos], padre.separadores[pos+1:]...)
	padre.hijos = append(padre.hijos[:pos+1], padre.hijos[pos+2:]...)
}

func (t *Tree[K, V]) corregirUnderflowInterno(n *nodo[K, V], padres []*nodo[K, V], posiciones []int) {
	if len(padres) == 0 {
		return
	}

	padre := padres[len(padres)-1]
	pos := posiciones[len(posiciones)-1]

	if pos+1 < len(padre.hijos) {
		hermanoDer := padre.hijos[pos+1]
		if len(hermanoDer.separadores) > t.minPorNodo {
			t.prestarSeparadorDerecha(n, hermanoDer, padre, pos)
			return
		}
	}

	if pos > 0 {
		hermanoIzq := padre.hijos[pos-1]
		if len(hermanoIzq.separadores) > t.minPorNodo {
			t.prestarSeparadorIzquierda(n, hermanoIzq, padre, pos-1)
			return
		}
	}

	if pos+1 < len(padre.hijos) {
		t.fusionarInternoConDerecha(n, padre, pos)
	} else {
		t.fusionarInternoConIzquierda(n, padre, pos-1)
	}

	t.ajustarRaiz(padre)
	if padre != t.raiz && len(padre.separadores) < t.minPorNodo {
		t.corregirUnderflowInterno(padre, padres[:len(padres)-1], posiciones[:len(posiciones)-1])
	}
}

func (t *Tree[K, V]) prestarSeparadorDerecha(nd, hermano *nodo[K, V], padre *nodo[K, V], pos int) {
	nd.separadores = append(nd.separadores, padre.separadores[pos])
	nd.hijos = append(nd.hijos, hermano.hijos[0])
	padre.separadores[pos] = hermano.separadores[0]
	hermano.separadores = hermano.separadores[1:]
	hermano.hijos = hermano.hijos[1:]
}

func (t *Tree[K, V]) prestarSeparadorIzquierda(nd, hermano *nodo[K, V], padre *nodo[K, V], pos int) {
	nd.separadores = append([]K{padre.separadores[pos]}, nd.separadores...)
	nd.hijos = append([]*nodo[K, V]{hermano.hijos[len(hermano.hijos)-1]}, nd.hijos...)
	padre.separadores[pos] = hermano.separadores[len(hermano.separadores)-1]
	hermano.separadores = hermano.separadores[:len(hermano.separadores)-1]
	hermano.hijos = hermano.hijos[:len(hermano.hijos)-1]
}

func (t *Tree[K, V]) fusionarInternoConDerecha(n *nodo[K, V], padre *nodo[K, V], pos int) {
	hermano := padre.hijos[pos+1]
	n.separadores = append(n.separadores, padre.separadores[pos])
	n.separadores = append(n.separadores, hermano.separadores...)
	n.hijos = append(n.hijos, hermano.hijos...)
	padre.separadores = append(padre.separadores[:pos], padre.separadores[pos+1:]...)
	padre.hijos = append(padre.hijos[:pos+1], padre.hijos[pos+2:]...)
}

func (t *Tree[K, V]) fusionarInternoConIzquierda(n *nodo[K, V], padre *nodo[K, V], pos int) {
	hermano := padre.hijos[pos]
	hermano.separadores = append(hermano.separadores, padre.separadores[pos])
	hermano.separadores = append(hermano.separadores, n.separadores...)
	hermano.hijos = append(hermano.hijos, n.hijos...)
	padre.separadores = append(padre.separadores[:pos], padre.separadores[pos+1:]...)
	padre.hijos = append(padre.hijos[:pos+1], padre.hijos[pos+2:]...)
}
