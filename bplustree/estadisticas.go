package bplustree

import "fmt"

//--------------Estadisticas del arbol (UI sin canvas)---------------

// EstadisticasArbol resume un índice sin serializar todos los nodos.
// Se usa cuando hay más de 500 canciones y el canvas se desactiva.
type EstadisticasArbol struct {
	TipoIndice     string `json:"tipoIndice"`
	EtiquetaIndice string `json:"etiquetaIndice"`
	TipoClave      string `json:"tipoClave"`
	TotalCanciones int    `json:"totalCanciones"`
	EntradasHoja   int    `json:"entradasHoja"`
	NodosInternos  int    `json:"nodosInternos"`
	NodosHoja      int    `json:"nodosHoja"`
	NodosTotales   int    `json:"nodosTotales"`
	Orden          int    `json:"orden"`
	Altura         int    `json:"altura"`
}

//--------------Conteo interno---------------

// contarNodos recorre el árbol y devuelve nodos internos, hojas y entradas totales.
func (t *Tree[K, V]) contarNodos() (internos, hojas, entradas int) {
	if t == nil || t.raiz == nil {
		return 0, 0, 0
	}
	var recorrer func(*nodo[K, V])
	recorrer = func(n *nodo[K, V]) {
		if n.esHoja {
			hojas++
			entradas += len(n.entradas)
			return
		}
		internos++
		for _, hijo := range n.hijos {
			recorrer(hijo)
		}
	}
	recorrer(t.raiz)
	return internos, hojas, entradas
}

// alturaArbol calcula la altura del árbol (raíz → hojas).
func (t *Tree[K, V]) alturaArbol() int {
	if t == nil || t.raiz == nil {
		return 0
	}
	var prof func(*nodo[K, V]) int
	prof = func(n *nodo[K, V]) int {
		if n.esHoja {
			return 1
		}
		maxH := 0
		for _, hijo := range n.hijos {
			if h := prof(hijo); h > maxH {
				maxH = h
			}
		}
		return maxH + 1
	}
	return prof(t.raiz)
}

// totalCanciones devuelve el total de entradas en el índice primario.
func (c *CatalogoIndices) totalCanciones() int {
	if c == nil || c.PorIndice == nil {
		return 0
	}
	_, _, ent := c.PorIndice.contarNodos()
	return ent
}

// armarStats construye el struct EstadisticasArbol para respuesta JSON.
func armarStats(tipo, etiqueta, tipoClave string, orden, canciones, in, ho, ent, alt int) EstadisticasArbol {
	return EstadisticasArbol{
		TipoIndice:     tipo,
		EtiquetaIndice: etiqueta,
		TipoClave:      tipoClave,
		TotalCanciones: canciones,
		EntradasHoja:   ent,
		NodosInternos:  in,
		NodosHoja:      ho,
		NodosTotales:   in + ho,
		Orden:          orden,
		Altura:         alt,
	}
}

//--------------API publica---------------

// EstadisticasPorTipo: conteos sin serializar todo el arbol.
// Tiempo: O(n) recorrido | Espacio: O(1) auxiliar
func (c *CatalogoIndices) EstadisticasPorTipo(tipo string) (EstadisticasArbol, error) {
	if c == nil {
		return EstadisticasArbol{}, fmt.Errorf("catálogo no inicializado")
	}
	canciones := c.totalCanciones()

	switch tipo {
	case "", "indice", "id":
		in, ho, ent := c.PorIndice.contarNodos()
		return armarStats("indice", "Indice (id)", "int", c.Orden, canciones, in, ho, ent, c.PorIndice.alturaArbol()), nil
	case "nombre", "trackname", "name":
		in, ho, ent := c.PorNombre.contarNodos()
		return armarStats("nombre", "TrackName (string)", "string", c.Orden, canciones, in, ho, ent, c.PorNombre.alturaArbol()), nil
	case "popularidad", "popularity":
		in, ho, ent := c.PorPopularidad.contarNodos()
		return armarStats("popularidad", "Popularity (int)", "int", c.Orden, canciones, in, ho, ent, c.PorPopularidad.alturaArbol()), nil
	case "tempo":
		in, ho, ent := c.PorTempo.contarNodos()
		return armarStats("tempo", "Tempo (float64)", "float64", c.Orden, canciones, in, ho, ent, c.PorTempo.alturaArbol()), nil
	case "danceability", "bailabilidad":
		in, ho, ent := c.PorDanceability.contarNodos()
		return armarStats("danceability", "Danceability (float64)", "float64", c.Orden, canciones, in, ho, ent, c.PorDanceability.alturaArbol()), nil
	default:
		return EstadisticasArbol{}, fmt.Errorf("tipo de índice desconocido: %s", tipo)
	}
}
