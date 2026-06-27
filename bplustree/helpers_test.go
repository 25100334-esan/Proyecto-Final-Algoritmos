package bplustree

func alturaArbol(raiz *nodo[int, Registro]) int {
	if raiz == nil {
		return 0
	}
	if raiz.esHoja {
		return 1
	}
	maxH := 0
	for _, hijo := range raiz.hijos {
		if h := alturaArbol(hijo); h > maxH {
			maxH = h
		}
	}
	return 1 + maxH
}

func hojaMasIzquierda(raiz *nodo[int, Registro]) *nodo[int, Registro] {
	nodo := raiz
	for !nodo.esHoja {
		nodo = nodo.hijos[0]
	}
	return nodo
}

func contarHojas(raiz *nodo[int, Registro]) int {
	if raiz == nil {
		return 0
	}
	if raiz.esHoja {
		return 1
	}
	total := 0
	for _, hijo := range raiz.hijos {
		total += contarHojas(hijo)
	}
	return total
}

func indicesDeHoja(hoja *nodo[int, Registro]) []int {
	indices := make([]int, len(hoja.entradas))
	for i, e := range hoja.entradas {
		indices[i] = e.clave
	}
	return indices
}

func recorrerSequenceSet(arbol *BPlusTree) []int {
	indices := []int{}
	hoja := hojaMasIzquierda(arbol.raiz)

	for hoja != nil {
		for _, e := range hoja.entradas {
			indices = append(indices, e.clave)
		}
		hoja = hoja.siguienteHoja
	}

	return indices
}
