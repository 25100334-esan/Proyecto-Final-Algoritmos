package bplustree

import "strings"

//-----------------Prefix B+-Tree Adicional ------------------
// Busqueda por TrackName: rango lexicografico y autocompletado por prefijo.
// Usa el mismo sequence set del arbol de strings (Comer).

// BuscarRangoString: range scan sobre nombres (case-insensitive).
// Si desde/hasta son una letra (M-O), incluye todos los nombres de esas letras.
// Tiempo: O(log_d n + k) | Espacio: O(k)
func BuscarRangoString(t *Tree[string, Registro], desde, hasta string) []Registro {
	desdeNorm := strings.ToLower(strings.TrimSpace(desde))
	hastaNorm := strings.ToLower(strings.TrimSpace(hasta))
	if desdeNorm == "" || hastaNorm == "" || t == nil {
		return nil
	}
	if strings.Compare(desdeNorm, hastaNorm) > 0 {
		return nil
	}

	if len(desdeNorm) == 1 && len(hastaNorm) == 1 {
		return buscarRangoStringSemi(t, desdeNorm, limiteSuperiorPrefijo(hastaNorm))
	}

	return buscarRangoStringInclusivo(t, desdeNorm, hastaNorm)
}

// BuscarPorPrefijo: [prefijo, limiteSuperior(prefijo)).
// Tiempo: O(log_d n + k) | Espacio: O(k)
func BuscarPorPrefijo(t *Tree[string, Registro], prefijo string) []Registro {
	prefijo = strings.TrimSpace(prefijo)
	if prefijo == "" || t == nil {
		return nil
	}
	prefijoNorm := strings.ToLower(prefijo)
	return buscarRangoStringSemi(t, prefijoNorm, limiteSuperiorPrefijo(prefijoNorm))
}

//--------------Helpers de rango string---------------

// limiteSuperiorPrefijo: "comed"->"comee", "o"->"p"
func limiteSuperiorPrefijo(prefijoNorm string) string {
	b := []byte(prefijoNorm)
	for i := len(b) - 1; i >= 0; i-- {
		if b[i] < 0xff {
			b[i]++
			return string(b[:i+1])
		}
	}
	return prefijoNorm + "\U0010FFFF"
}

func buscarRangoStringInclusivo(t *Tree[string, Registro], desdeNorm, hastaNorm string) []Registro {
	hoja := irAHojaPrefijo(t, desdeNorm)
	resultado := make([]Registro, 0)

	for hoja != nil {
		for _, e := range hoja.entradas {
			claveNorm := strings.ToLower(e.clave)
			if claveNorm < desdeNorm {
				continue
			}
			if claveNorm > hastaNorm {
				return resultado
			}
			resultado = append(resultado, e.valor)
		}
		if len(hoja.entradas) == 0 {
			break
		}
		ultimaNorm := strings.ToLower(hoja.entradas[len(hoja.entradas)-1].clave)
		if ultimaNorm < desdeNorm {
			hoja = hoja.siguienteHoja
			continue
		}
		if ultimaNorm > hastaNorm {
			break
		}
		hoja = hoja.siguienteHoja
	}
	return resultado
}

func buscarRangoStringSemi(t *Tree[string, Registro], desdeNorm, hastaExcl string) []Registro {
	hoja := irAHojaPrefijo(t, desdeNorm)
	resultado := make([]Registro, 0)

	for hoja != nil {
		for _, e := range hoja.entradas {
			claveNorm := strings.ToLower(e.clave)
			if claveNorm < desdeNorm {
				continue
			}
			if claveNorm >= hastaExcl {
				return resultado
			}
			resultado = append(resultado, e.valor)
		}
		if len(hoja.entradas) == 0 {
			break
		}
		ultimaNorm := strings.ToLower(hoja.entradas[len(hoja.entradas)-1].clave)
		if ultimaNorm < desdeNorm {
			hoja = hoja.siguienteHoja
			continue
		}
		if ultimaNorm >= hastaExcl {
			break
		}
		hoja = hoja.siguienteHoja
	}
	return resultado
}

func irAHojaPrefijo(t *Tree[string, Registro], prefijoNorm string) *nodo[string, Registro] {
	n := t.raiz
	for !n.esHoja {
		i := 0
		for i < len(n.separadores) && strings.ToLower(n.separadores[i]) < prefijoNorm {
			i++
		}
		n = n.hijos[i]
	}
	return n
}
