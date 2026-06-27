package bplustree

import (
	"fmt"
	"strconv"
)

//--------------Tipos de exportacion (JSON para Vue)---------------

// NodoExport representa un nodo del árbol listo para serializar a JSON.
type NodoExport struct {
	ID          int      `json:"id"`
	EsHoja      bool     `json:"esHoja"`
	Separadores []string `json:"separadores,omitempty"`
	Claves      []string `json:"claves"`
	Indices     []int    `json:"indices,omitempty"`
	IdsRegistro []int    `json:"idsRegistro,omitempty"`
	Hijos       []int    `json:"hijos,omitempty"`
}

// EstructuraExport contiene el árbol completo para la visualización en Vue.js.
type EstructuraExport struct {
	Orden          int          `json:"orden"`
	MinClaves      int          `json:"minClaves"`
	MaxClaves      int          `json:"maxClaves"`
	MaxHijos       int          `json:"maxHijos"`
	Raiz           int          `json:"raiz"`
	Nodos          []NodoExport `json:"nodos"`
	CadenaHojas    []int        `json:"cadenaHojas"`
	TotalCanciones int          `json:"totalCanciones"`
	TipoIndice     string       `json:"tipoIndice"`
	TipoClave      string       `json:"tipoClave"`
	EtiquetaIndice string       `json:"etiquetaIndice"`
	ClavesUnicas   bool         `json:"clavesUnicas"`
	IndicesActivos []string     `json:"indicesActivos,omitempty"`
}

//--------------Exportacion por tipo de clave (interno)---------------

// exportarTreeInt serializa un árbol con claves int (ID, Popularidad).
func exportarTreeInt(t *Tree[int, Registro], tipo, etiqueta string) EstructuraExport {
	if t == nil || t.raiz == nil {
		return EstructuraExport{TipoIndice: tipo, TipoClave: "int", EtiquetaIndice: etiqueta, ClavesUnicas: true}
	}

	contador := 0
	nodos := []NodoExport{}
	mapNodos := map[*nodo[int, Registro]]int{}

	var recorrer func(*nodo[int, Registro]) int
	recorrer = func(n *nodo[int, Registro]) int {
		id := contador
		contador++
		mapNodos[n] = id

		exp := NodoExport{ID: id, EsHoja: n.esHoja}
		if n.esHoja {
			exp.Claves = make([]string, len(n.entradas))
			exp.Indices = make([]int, len(n.entradas))
			exp.IdsRegistro = make([]int, len(n.entradas))
			for i, e := range n.entradas {
				exp.Claves[i] = strconv.Itoa(e.clave)
				exp.Indices[i] = e.clave
				exp.IdsRegistro[i] = e.valor.Indice
			}
		} else {
			for _, s := range n.separadores {
				exp.Separadores = append(exp.Separadores, strconv.Itoa(s))
			}
			exp.Hijos = make([]int, len(n.hijos))
			for i, hijo := range n.hijos {
				exp.Hijos[i] = recorrer(hijo)
			}
		}
		nodos = append(nodos, exp)
		return id
	}

	raizID := recorrer(t.raiz)
	return armarExport(t.orden, t.minPorNodo, t.maxPorNodo, raizID, nodos, mapNodos, t.primeraHoja(), tipo, "int", etiqueta, t.clavesUnicas)
}

// exportarTreeString serializa un árbol con claves string (TrackName).
func exportarTreeString(t *Tree[string, Registro], tipo, etiqueta string) EstructuraExport {
	if t == nil || t.raiz == nil {
		return EstructuraExport{TipoIndice: tipo, TipoClave: "string", EtiquetaIndice: etiqueta, ClavesUnicas: false}
	}

	contador := 0
	nodos := []NodoExport{}
	mapNodos := map[*nodo[string, Registro]]int{}

	var recorrer func(*nodo[string, Registro]) int
	recorrer = func(n *nodo[string, Registro]) int {
		id := contador
		contador++
		mapNodos[n] = id

		exp := NodoExport{ID: id, EsHoja: n.esHoja}
		if n.esHoja {
			exp.Claves = make([]string, len(n.entradas))
			exp.IdsRegistro = make([]int, len(n.entradas))
			for i, e := range n.entradas {
				exp.Claves[i] = truncarClave(e.clave)
				exp.IdsRegistro[i] = e.valor.Indice
			}
		} else {
			for _, s := range n.separadores {
				exp.Separadores = append(exp.Separadores, truncarClave(s))
			}
			exp.Hijos = make([]int, len(n.hijos))
			for i, hijo := range n.hijos {
				exp.Hijos[i] = recorrer(hijo)
			}
		}
		nodos = append(nodos, exp)
		return id
	}

	raizID := recorrer(t.raiz)
	return armarExportString(t.orden, t.minPorNodo, t.maxPorNodo, raizID, nodos, mapNodos, t.primeraHoja(), tipo, etiqueta, t.clavesUnicas)
}

// exportarTreeFloat serializa un árbol con claves float64 (Tempo, Danceability).
func exportarTreeFloat(t *Tree[float64, Registro], tipo, etiqueta string) EstructuraExport {
	if t == nil || t.raiz == nil {
		return EstructuraExport{TipoIndice: tipo, TipoClave: "float64", EtiquetaIndice: etiqueta, ClavesUnicas: false}
	}

	contador := 0
	nodos := []NodoExport{}
	mapNodos := map[*nodo[float64, Registro]]int{}

	var recorrer func(*nodo[float64, Registro]) int
	recorrer = func(n *nodo[float64, Registro]) int {
		id := contador
		contador++
		mapNodos[n] = id

		exp := NodoExport{ID: id, EsHoja: n.esHoja}
		if n.esHoja {
			exp.Claves = make([]string, len(n.entradas))
			exp.IdsRegistro = make([]int, len(n.entradas))
			for i, e := range n.entradas {
				exp.Claves[i] = strconv.FormatFloat(e.clave, 'f', -1, 64)
				exp.IdsRegistro[i] = e.valor.Indice
			}
		} else {
			for _, s := range n.separadores {
				exp.Separadores = append(exp.Separadores, strconv.FormatFloat(s, 'f', -1, 64))
			}
			exp.Hijos = make([]int, len(n.hijos))
			for i, hijo := range n.hijos {
				exp.Hijos[i] = recorrer(hijo)
			}
		}
		nodos = append(nodos, exp)
		return id
	}

	raizID := recorrer(t.raiz)
	return armarExportFloat(t.orden, t.minPorNodo, t.maxPorNodo, raizID, nodos, mapNodos, t.primeraHoja(), tipo, etiqueta, t.clavesUnicas)
}

//--------------Armado de estructura y utilidades---------------

// armarExport construye EstructuraExport para árboles con clave int.
func armarExport(
	orden, min, max, raizID int,
	nodos []NodoExport,
	mapNodos map[*nodo[int, Registro]]int,
	primera *nodo[int, Registro],
	tipo, tipoClave, etiqueta string,
	clavesUnicas bool,
) EstructuraExport {
	cadenaHojas := []int{}
	total := 0
	for hoja := primera; hoja != nil; hoja = hoja.siguienteHoja {
		cadenaHojas = append(cadenaHojas, mapNodos[hoja])
		total += len(hoja.entradas)
	}
	return EstructuraExport{
		Orden: orden, MinClaves: min, MaxClaves: max, MaxHijos: 2*orden + 1,
		Raiz: raizID, Nodos: nodos, CadenaHojas: cadenaHojas, TotalCanciones: total,
		TipoIndice: tipo, TipoClave: tipoClave, EtiquetaIndice: etiqueta, ClavesUnicas: clavesUnicas,
		IndicesActivos: nombresIndicesActivos(),
	}
}

// armarExportString construye EstructuraExport para árboles con clave string.
func armarExportString(orden, min, max, raizID int, nodos []NodoExport, mapNodos map[*nodo[string, Registro]]int, primera *nodo[string, Registro], tipo, etiqueta string, clavesUnicas bool) EstructuraExport {
	cadenaHojas := []int{}
	total := 0
	for hoja := primera; hoja != nil; hoja = hoja.siguienteHoja {
		cadenaHojas = append(cadenaHojas, mapNodos[hoja])
		total += len(hoja.entradas)
	}
	return EstructuraExport{
		Orden: orden, MinClaves: min, MaxClaves: max, MaxHijos: 2*orden + 1,
		Raiz: raizID, Nodos: nodos, CadenaHojas: cadenaHojas, TotalCanciones: total,
		TipoIndice: tipo, TipoClave: "string", EtiquetaIndice: etiqueta, ClavesUnicas: clavesUnicas,
		IndicesActivos: nombresIndicesActivos(),
	}
}

// armarExportFloat construye EstructuraExport para árboles con clave float64.
func armarExportFloat(orden, min, max, raizID int, nodos []NodoExport, mapNodos map[*nodo[float64, Registro]]int, primera *nodo[float64, Registro], tipo, etiqueta string, clavesUnicas bool) EstructuraExport {
	cadenaHojas := []int{}
	total := 0
	for hoja := primera; hoja != nil; hoja = hoja.siguienteHoja {
		cadenaHojas = append(cadenaHojas, mapNodos[hoja])
		total += len(hoja.entradas)
	}
	return EstructuraExport{
		Orden: orden, MinClaves: min, MaxClaves: max, MaxHijos: 2*orden + 1,
		Raiz: raizID, Nodos: nodos, CadenaHojas: cadenaHojas, TotalCanciones: total,
		TipoIndice: tipo, TipoClave: "float64", EtiquetaIndice: etiqueta, ClavesUnicas: clavesUnicas,
		IndicesActivos: nombresIndicesActivos(),
	}
}

// truncarClave acorta nombres largos para la visualización en el canvas.
// Usa runas (no bytes) para coincidir con truncarClaveVisual en Vue.
func truncarClave(s string) string {
	r := []rune(s)
	if len(r) <= 14 {
		return s
	}
	return string(r[:12]) + "…"
}

// nombresIndicesActivos lista los cinco índices disponibles en el catálogo.
func nombresIndicesActivos() []string {
	return []string{"Indice (id)", "TrackName", "Popularity", "Tempo", "Danceability"}
}

//--------------API publica de exportacion---------------

// ExportarEstructura: serializa indice primario.
// Tiempo: O(n) | Espacio: O(n) por el JSON generado
func (t *BPlusTree) ExportarEstructura() EstructuraExport {
	return exportarTreeInt(t.Tree, "indice", "Indice (id)")
}

// ExportarPorTipo serializa el árbol del tipo indicado para el canvas académico.
// Valores aceptados: indice, nombre, popularidad, tempo, danceability.
func (c *CatalogoIndices) ExportarPorTipo(tipo string) (EstructuraExport, error) {
	if c == nil {
		return EstructuraExport{}, fmt.Errorf("catálogo no inicializado")
	}
	switch tipo {
	case "", "indice", "id":
		return c.PorIndice.ExportarEstructura(), nil
	case "nombre", "trackname", "name":
		return exportarTreeString(c.PorNombre, "nombre", "TrackName (string)"), nil
	case "popularidad", "popularity":
		return exportarTreeInt(c.PorPopularidad, "popularidad", "Popularity (int)"), nil
	case "tempo":
		return exportarTreeFloat(c.PorTempo, "tempo", "Tempo (float64)"), nil
	case "danceability", "bailabilidad":
		return exportarTreeFloat(c.PorDanceability, "danceability", "Danceability (float64)"), nil
	default:
		return EstructuraExport{}, fmt.Errorf("tipo de índice desconocido: %s", tipo)
	}
}

// ExportarEstructura delega al índice primario del catálogo.
func (c *CatalogoIndices) ExportarEstructura() EstructuraExport {
	exp, _ := c.ExportarPorTipo("indice")
	return exp
}
