package bplustree

import "cmp"

//--------------Insercion (API publica)---------------

// Insertar: agrega clave-valor; si la hoja rebalsa hace split.
// Tiempo: O(log_d n) | Espacio: O(log_d n) por la ruta de padres
func (t *Tree[K, V]) Insertar(clave K, valor V) {
	if t.clavesUnicas {
		if _, ok := t.Buscar(clave); ok {
			return
		}
	}

	hoja, padres, posiciones := t.bajarGuardandoRuta(clave)
	pos := posicionEntrada(hoja.entradas, clave)
	hoja.entradas = append(hoja.entradas, entrada[K, V]{})
	copy(hoja.entradas[pos+1:], hoja.entradas[pos:])
	hoja.entradas[pos] = entrada[K, V]{clave: clave, valor: valor}

	if len(hoja.entradas) <= t.maxPorNodo {
		return
	}

	claveCopiada, nuevaHoja := t.partirHoja(hoja)
	t.subirDivisionHoja(padres, posiciones, claveCopiada, nuevaHoja)
}

//--------------Split y overflow (interno)---------------

func (t *Tree[K, V]) agregarEntradaOrdenada(hoja *nodo[K, V], clave K, valor V) {
	pos := posicionEntrada(hoja.entradas, clave)
	hoja.entradas = append(hoja.entradas, entrada[K, V]{})
	copy(hoja.entradas[pos+1:], hoja.entradas[pos:])
	hoja.entradas[pos] = entrada[K, V]{clave: clave, valor: valor}
}

// partirHoja: divide hoja llena; promocion por copia al padre.
func (t *Tree[K, V]) partirHoja(hoja *nodo[K, V]) (K, *nodo[K, V]) {
	medio := (len(hoja.entradas) + 1) / 2

	derecha := &nodo[K, V]{
		esHoja:        true,
		entradas:      append([]entrada[K, V]{}, hoja.entradas[medio:]...),
		siguienteHoja: hoja.siguienteHoja,
	}
	hoja.entradas = hoja.entradas[:medio]
	hoja.siguienteHoja = derecha

	return derecha.entradas[0].clave, derecha
}

func insertarEnNodoInterno[K cmp.Ordered, V any](nodo *nodo[K, V], posHijo int, clave K, hijo *nodo[K, V]) {
	nodo.separadores = append(nodo.separadores, *new(K))
	copy(nodo.separadores[posHijo+1:], nodo.separadores[posHijo:])
	nodo.separadores[posHijo] = clave

	nodo.hijos = append(nodo.hijos, nil)
	copy(nodo.hijos[posHijo+2:], nodo.hijos[posHijo+1:])
	nodo.hijos[posHijo+1] = hijo
}

func (t *Tree[K, V]) subirDivisionHoja(padres []*nodo[K, V], posiciones []int, clave K, nuevaHoja *nodo[K, V]) {
	if len(padres) == 0 {
		t.raiz = &nodo[K, V]{
			esHoja:      false,
			separadores: []K{clave},
			hijos:       []*nodo[K, V]{t.raiz, nuevaHoja},
		}
		return
	}

	padre := padres[len(padres)-1]
	posHijo := posiciones[len(posiciones)-1]
	insertarEnNodoInterno(padre, posHijo, clave, nuevaHoja)

	if len(padre.separadores) <= t.maxPorNodo {
		return
	}

	claveSubida, nuevoInterno := t.partirNodoInterno(padre)
	t.subirDivisionInterna(padres[:len(padres)-1], posiciones[:len(posiciones)-1], claveSubida, nuevoInterno)
}

func (t *Tree[K, V]) partirNodoInterno(n *nodo[K, V]) (K, *nodo[K, V]) {
	medio := len(n.separadores) / 2
	claveSubida := n.separadores[medio]

	derecha := &nodo[K, V]{
		esHoja:      false,
		separadores: append([]K{}, n.separadores[medio+1:]...),
		hijos:       append([]*nodo[K, V]{}, n.hijos[medio+1:]...),
	}

	n.separadores = n.separadores[:medio]
	n.hijos = n.hijos[:medio+1]

	return claveSubida, derecha
}

func (t *Tree[K, V]) subirDivisionInterna(padres []*nodo[K, V], posiciones []int, clave K, nuevoHijo *nodo[K, V]) {
	if len(padres) == 0 {
		t.raiz = &nodo[K, V]{
			esHoja:      false,
			separadores: []K{clave},
			hijos:       []*nodo[K, V]{t.raiz, nuevoHijo},
		}
		return
	}

	padre := padres[len(padres)-1]
	posHijo := posiciones[len(posiciones)-1]
	insertarEnNodoInterno(padre, posHijo, clave, nuevoHijo)

	if len(padre.separadores) <= t.maxPorNodo {
		return
	}

	claveSubida, nuevoInterno := t.partirNodoInterno(padre)
	t.subirDivisionInterna(padres[:len(padres)-1], posiciones[:len(posiciones)-1], claveSubida, nuevoInterno)
}
