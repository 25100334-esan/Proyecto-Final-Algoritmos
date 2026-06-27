package bplustree

//--------------Busqueda (API publica)---------------

// Buscar: baja hasta la hoja y devuelve el primer valor con clave exacta.
// Tiempo: O(log_d n) | Espacio: O(1) auxiliar
func (t *Tree[K, V]) Buscar(clave K) (V, bool) {
	hoja := t.irAHoja(clave)
	pos := posicionEntrada(hoja.entradas, clave)
	if pos < len(hoja.entradas) && compare(hoja.entradas[pos].clave, clave) == 0 {
		return hoja.entradas[pos].valor, true
	}
	var zero V
	return zero, false
}

// BuscarTodos: todos los valores con la misma clave (indices secundarios).
// Tiempo: O(log_d n + k) | Espacio: O(k)
func (t *Tree[K, V]) BuscarTodos(clave K) []V {
	return t.BuscarExactos(clave)
}

// BuscarExactos: igual que BuscarRango(k, k) recorriendo el sequence set.
// Tiempo: O(log_d n + k) | Espacio: O(k)
func (t *Tree[K, V]) BuscarExactos(clave K) []V {
	return t.BuscarRango(clave, clave)
}

// BuscarRango: range scan horizontal entre dos claves.
// Tiempo: O(log_d n + k) | Espacio: O(k) por el slice resultado
func (t *Tree[K, V]) BuscarRango(desde, hasta K) []V {
	if compare(desde, hasta) > 0 {
		return nil
	}

	hoja := t.irAHoja(desde)
	resultado := []V{}

	for hoja != nil {
		for _, e := range hoja.entradas {
			if menor(e.clave, desde) {
				continue
			}
			if mayor(e.clave, hasta) {
				return resultado
			}
			resultado = append(resultado, e.valor)
		}
		if len(hoja.entradas) == 0 {
			break
		}
		ultima := hoja.entradas[len(hoja.entradas)-1].clave
		if mayor(ultima, hasta) {
			break
		}
		hoja = hoja.siguienteHoja
	}
	return resultado
}
