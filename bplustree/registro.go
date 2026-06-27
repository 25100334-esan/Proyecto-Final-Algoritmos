package bplustree

//--------------Modelo de cancion (Spotify)---------------

// Registro: una fila del dataset. Indice = id en SQLite.
type Registro struct {
	Indice           int
	TrackID          string
	Artists          string
	AlbumName        string
	TrackName        string
	Popularity       int
	DurationMs       int
	Explicit         bool
	Danceability     float64
	Energy           float64
	Key              int // musical_key en la BD
	Loudness         float64
	Mode             int // musical_mode en la BD
	Speechiness      float64
	Acousticness     float64
	Instrumentalness float64
	Liveness         float64
	Valence          float64
	Tempo            float64
	TimeSignature    int
	TrackGenre       string
}

//--------------Indice primario por ID---------------

// BPlusTree: indice principal (claves unicas por Indice).
type BPlusTree struct {
	*Tree[int, Registro]
}

// New: arbol vacio indexado por ID.
// Tiempo: O(1) | Espacio: O(1)
func New(orden int) *BPlusTree {
	return &BPlusTree{Tree: NewTree[int, Registro](orden, true)}
}

// Insertar: agrega cancion por Indice (ignora duplicados).
// Tiempo: O(log_d n) | Espacio: O(log_d n)
func (t *BPlusTree) Insertar(reg Registro) {
	t.Tree.Insertar(reg.Indice, reg)
}

// Buscar: localiza por Indice.
// Tiempo: O(log_d n) | Espacio: O(1)
func (t *BPlusTree) Buscar(indice int) (*Registro, bool) {
	v, ok := t.Tree.Buscar(indice)
	if !ok {
		return nil, false
	}
	return &v, true
}

// BuscarRango: range scan por ID.
// Tiempo: O(log_d n + k) | Espacio: O(k)
func (t *BPlusTree) BuscarRango(desde, hasta int) []Registro {
	return t.Tree.BuscarRango(desde, hasta)
}

// Eliminar: borra por Indice.
// Tiempo: O(log_d n) | Espacio: O(h)
func (t *BPlusTree) Eliminar(indice int) {
	t.Tree.Eliminar(indice)
}

//--------------Catalogo de 5 indices en RAM---------------

// CatalogoIndices: un B+-Tree por criterio (Comer).
type CatalogoIndices struct {
	Orden           int
	PorIndice       *BPlusTree
	PorNombre       *Tree[string, Registro]
	PorPopularidad  *Tree[int, Registro]
	PorTempo        *Tree[float64, Registro]
	PorDanceability *Tree[float64, Registro]
}

func NuevoCatalogo(orden int) *CatalogoIndices {
	if orden < 2 {
		orden = 2
	}
	return &CatalogoIndices{
		Orden:           orden,
		PorIndice:       New(orden),
		PorNombre:       NewTree[string, Registro](orden, false),
		PorPopularidad:  NewTree[int, Registro](orden, false),
		PorTempo:        NewTree[float64, Registro](orden, false),
		PorDanceability: NewTree[float64, Registro](orden, false),
	}
}

// InsertarRegistro: inserta en los 5 arboles a la vez.
// Tiempo: O(5 * log_d n) | Espacio: O(log_d n)
func (c *CatalogoIndices) InsertarRegistro(reg Registro) {
	c.PorIndice.Insertar(reg)
	c.PorNombre.Insertar(reg.TrackName, reg)
	c.PorPopularidad.Insertar(reg.Popularity, reg)
	c.PorTempo.Insertar(reg.Tempo, reg)
	c.PorDanceability.Insertar(reg.Danceability, reg)
}

func (c *CatalogoIndices) EliminarPorIndice(indice int) bool {
	reg, ok := c.PorIndice.Buscar(indice)
	if !ok {
		return false
	}
	c.EliminarRegistro(*reg)
	return true
}

// EliminarRegistro: cascada en primario + secundarios.
// Tiempo: O(5 * log_d n) | Espacio: O(h)
func (c *CatalogoIndices) EliminarRegistro(reg Registro) {
	c.PorIndice.Eliminar(reg.Indice)
	EliminarExacto(c.PorNombre, reg.TrackName, reg.Indice)
	EliminarExacto(c.PorPopularidad, reg.Popularity, reg.Indice)
	EliminarExacto(c.PorTempo, reg.Tempo, reg.Indice)
	EliminarExacto(c.PorDanceability, reg.Danceability, reg.Indice)
}

func coincideIndice(indice int) func(Registro) bool {
	return func(r Registro) bool { return r.Indice == indice }
}

//--------------Consultas del catalogo---------------

// BuscarRangoPopularidad: Tiempo O(log_d n + k) | Espacio O(k)
func (c *CatalogoIndices) BuscarRangoPopularidad(desde, hasta int) []Registro {
	return c.PorPopularidad.BuscarRango(desde, hasta)
}

func (c *CatalogoIndices) BuscarRangoTempo(desde, hasta float64) []Registro {
	return c.PorTempo.BuscarRango(desde, hasta)
}

func (c *CatalogoIndices) BuscarRangoDanceability(desde, hasta float64) []Registro {
	return c.PorDanceability.BuscarRango(desde, hasta)
}

func (c *CatalogoIndices) BuscarPorPrefijoNombre(prefijo string) []Registro {
	return BuscarPorPrefijo(c.PorNombre, prefijo)
}

func (c *CatalogoIndices) BuscarRangoNombre(desde, hasta string) []Registro {
	return BuscarRangoString(c.PorNombre, desde, hasta)
}
