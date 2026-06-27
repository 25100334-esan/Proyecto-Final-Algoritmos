package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"bplustree-proyecto/bplustree"
	"bplustree-proyecto/database"
)

// Servidor expone el catálogo de índices B+ en memoria como API JSON para Vue.js.
type Servidor struct {
	catalogo *bplustree.CatalogoIndices
	db       *database.DB
}

// configInicial recibe orden y límite desde la pantalla de configuración de Vue.js.
type configInicial struct {
	Orden  int `json:"orden"`
	Limite int `json:"limite"`
}

// NuevoServidor crea el director de orquesta. Los índices nacen vacíos hasta /api/inicializar.
func NuevoServidor(db *database.DB) *Servidor {
	return &Servidor{db: db}
}

// Iniciar registra las rutas y escucha peticiones HTTP.
func (s *Servidor) Iniciar(puerto string) error {
	http.HandleFunc("/api/health", s.habilitarCORS(s.manejarHealth))
	http.HandleFunc("/api/inicializar", s.habilitarCORS(s.manejarInicializar))
	http.HandleFunc("/api/inicializar-demo", s.habilitarCORS(s.manejarInicializarDemo))
	http.HandleFunc("/api/insertar-arbol", s.habilitarCORS(s.manejarInsertarArbol))
	http.HandleFunc("/api/buscar", s.habilitarCORS(s.manejarBuscar))
	http.HandleFunc("/api/buscar-nombre", s.habilitarCORS(s.manejarBuscarNombre))
	http.HandleFunc("/api/buscar-prefijo", s.habilitarCORS(s.manejarBuscarPrefijo))
	http.HandleFunc("/api/buscar-tempo", s.habilitarCORS(s.manejarBuscarTempoExacto))
	http.HandleFunc("/api/rango", s.habilitarCORS(s.manejarRango))
	http.HandleFunc("/api/rango-nombre", s.habilitarCORS(s.manejarRangoNombre))
	http.HandleFunc("/api/rango-popularidad", s.habilitarCORS(s.manejarRangoPopularidad))
	http.HandleFunc("/api/rango-tempo", s.habilitarCORS(s.manejarRangoTempo))
	http.HandleFunc("/api/rango-danceability", s.habilitarCORS(s.manejarRangoDanceability))
	http.HandleFunc("/api/insertar", s.habilitarCORS(s.manejarInsertar))
	http.HandleFunc("/api/eliminar", s.habilitarCORS(s.manejarEliminar))
	http.HandleFunc("/api/estructura", s.habilitarCORS(s.manejarEstructura))
	http.HandleFunc("/api/estadisticas", s.habilitarCORS(s.manejarEstadisticas))
	http.HandleFunc("/api/conteo-bd", s.habilitarCORS(s.manejarConteoBD))
	http.HandleFunc("/api/siguiente-indice", s.habilitarCORS(s.manejarSiguienteIndice))
	http.HandleFunc("/api/indices", s.habilitarCORS(s.manejarIndices))

	return http.ListenAndServe(puerto, nil)
}

func (s *Servidor) manejarHealth(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		responderError(w, http.StatusMethodNotAllowed, "use GET")
		return
	}
	responderJSON(w, http.StatusOK, map[string]interface{}{
		"ok":      true,
		"orden":   "d>=2 (Comer)",
		"generico": true,
	})
}

func (s *Servidor) manejarInicializar(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		responderError(w, http.StatusMethodNotAllowed, "use POST")
		return
	}

	var config configInicial
	if err := json.NewDecoder(r.Body).Decode(&config); err != nil {
		responderError(w, http.StatusBadRequest, "JSON inválido")
		return
	}
	if config.Orden < 2 {
		responderError(w, http.StatusBadRequest, "orden (d) debe ser al menos 2 según Comer")
		return
	}
	if config.Limite < 1 {
		responderError(w, http.StatusBadRequest, "limite debe ser al menos 1")
		return
	}

	s.catalogo = bplustree.NuevoCatalogo(config.Orden)

	total, err := s.db.CargarDatosEnCatalogo(s.catalogo, config.Limite)
	if err != nil {
		s.catalogo = nil
		responderError(w, http.StatusInternalServerError, err.Error())
		return
	}

	responderJSON(w, http.StatusOK, map[string]interface{}{
		"mensaje": "catálogo de índices inicializado",
		"orden":   config.Orden,
		"limite":  config.Limite,
		"total":   total,
		"indices": s.nombresIndices(),
	})
}

const limiteDemoInicialMax = 50

func (s *Servidor) manejarInicializarDemo(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		responderError(w, http.StatusMethodNotAllowed, "use POST")
		return
	}

	var config configInicial
	if err := json.NewDecoder(r.Body).Decode(&config); err != nil {
		responderError(w, http.StatusBadRequest, "JSON inválido")
		return
	}
	if config.Orden < 2 {
		responderError(w, http.StatusBadRequest, "orden (d) debe ser al menos 2 según Comer")
		return
	}
	if config.Limite < 1 || config.Limite > limiteDemoInicialMax {
		responderError(w, http.StatusBadRequest,
			fmt.Sprintf("demo paso a paso: limite debe estar entre 1 y %d", limiteDemoInicialMax))
		return
	}

	s.catalogo = bplustree.NuevoCatalogo(config.Orden)

	canciones, err := s.db.ListarTracks(config.Limite)
	if err != nil {
		s.catalogo = nil
		responderError(w, http.StatusInternalServerError, err.Error())
		return
	}

	responderJSON(w, http.StatusOK, map[string]interface{}{
		"mensaje":   "árbol vacío listo para demo",
		"orden":     config.Orden,
		"limite":    config.Limite,
		"total":     len(canciones),
		"canciones": canciones,
		"indices":   s.nombresIndices(),
	})
}

func (s *Servidor) manejarInsertarArbol(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		responderError(w, http.StatusMethodNotAllowed, "use POST")
		return
	}
	if !s.verificarCatalogo(w) {
		return
	}

	var reg bplustree.Registro
	if err := json.NewDecoder(r.Body).Decode(&reg); err != nil {
		responderError(w, http.StatusBadRequest, "JSON inválido")
		return
	}
	if reg.Indice < 0 {
		responderError(w, http.StatusBadRequest, "Indice inválido (debe ser ≥ 0)")
		return
	}

	s.catalogo.InsertarRegistro(reg)

	responderJSON(w, http.StatusOK, map[string]interface{}{
		"mensaje": "canción insertada en catálogo de índices",
		"indice":  reg.Indice,
	})
}

func (s *Servidor) manejarBuscar(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		responderError(w, http.StatusMethodNotAllowed, "use GET")
		return
	}
	if !s.verificarCatalogo(w) {
		return
	}

	indice, err := leerEnteroQuery(r, "indice")
	if err != nil {
		responderError(w, http.StatusBadRequest, "parámetro indice inválido")
		return
	}

	reg, ok := s.catalogo.PorIndice.Buscar(indice)
	if !ok {
		responderError(w, http.StatusNotFound, "canción no encontrada")
		return
	}

	responderJSON(w, http.StatusOK, reg)
}

func (s *Servidor) manejarBuscarNombre(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		responderError(w, http.StatusMethodNotAllowed, "use GET")
		return
	}
	if !s.verificarCatalogo(w) {
		return
	}

	nombre := strings.TrimSpace(r.URL.Query().Get("nombre"))
	if nombre == "" {
		responderError(w, http.StatusBadRequest, "parámetro nombre requerido")
		return
	}

	todos := s.catalogo.PorNombre.BuscarExactos(nombre)
	if len(todos) == 0 {
		responderError(w, http.StatusNotFound, "canción no encontrada")
		return
	}

	responderJSON(w, http.StatusOK, map[string]interface{}{
		"total":     len(todos),
		"canciones": todos,
		"cancion":   todos[0],
		"metodo":    "BuscarRango(k,k) + sequence set — todas las coincidencias exactas",
	})
}

func (s *Servidor) manejarBuscarPrefijo(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		responderError(w, http.StatusMethodNotAllowed, "use GET")
		return
	}
	if !s.verificarCatalogo(w) {
		return
	}

	prefijo := strings.TrimSpace(r.URL.Query().Get("prefijo"))
	if prefijo == "" {
		responderError(w, http.StatusBadRequest, "parámetro prefijo requerido")
		return
	}

	todos := s.catalogo.BuscarPorPrefijoNombre(prefijo)

	responderJSON(w, http.StatusOK, map[string]interface{}{
		"total":     len(todos),
		"canciones": todos,
		"prefijo":   prefijo,
		"metodo":    "irAHoja(prefijo) + sequence set + strings.HasPrefix — Prefix B+-Tree (Comer)",
	})
}

func (s *Servidor) manejarBuscarTempoExacto(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		responderError(w, http.StatusMethodNotAllowed, "use GET")
		return
	}
	if !s.verificarCatalogo(w) {
		return
	}

	tempo, err := leerFloatQuery(r, "tempo")
	if err != nil {
		responderError(w, http.StatusBadRequest, "parámetro tempo inválido")
		return
	}

	todos := s.catalogo.PorTempo.BuscarExactos(tempo)
	if len(todos) == 0 {
		responderError(w, http.StatusNotFound, "ninguna canción con ese tempo")
		return
	}

	responderJSON(w, http.StatusOK, map[string]interface{}{
		"total":     len(todos),
		"canciones": todos,
		"tempo":     tempo,
		"metodo":    "BuscarRango(tempo,tempo) — claves duplicadas vía sequence set",
	})
}

func (s *Servidor) manejarRango(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		responderError(w, http.StatusMethodNotAllowed, "use GET")
		return
	}
	if !s.verificarCatalogo(w) {
		return
	}

	campo := strings.ToLower(strings.TrimSpace(r.URL.Query().Get("campo")))

	switch campo {
	case "popularidad":
		inicio, err := leerEnteroQuery(r, "inicio")
		if err != nil {
			responderError(w, http.StatusBadRequest, "parámetro inicio inválido")
			return
		}
		fin, err := leerEnteroQuery(r, "fin")
		if err != nil {
			responderError(w, http.StatusBadRequest, "parámetro fin inválido")
			return
		}
		resultado := s.catalogo.BuscarRangoPopularidad(inicio, fin)
		responderJSON(w, http.StatusOK, map[string]interface{}{
			"campo":     "popularidad",
			"indice":    "Popularity",
			"inicio":    inicio,
			"fin":       fin,
			"total":     len(resultado),
			"canciones": resultado,
			"metodo":    "BuscarRango + recorrido horizontal del sequence set (Comer)",
		})
		return
	case "tempo":
		inicio, err := leerFloatQuery(r, "inicio")
		if err != nil {
			responderError(w, http.StatusBadRequest, "parámetro inicio inválido")
			return
		}
		fin, err := leerFloatQuery(r, "fin")
		if err != nil {
			responderError(w, http.StatusBadRequest, "parámetro fin inválido")
			return
		}
		resultado := s.catalogo.BuscarRangoTempo(inicio, fin)
		responderJSON(w, http.StatusOK, map[string]interface{}{
			"campo":     "tempo",
			"indice":    "Tempo",
			"inicio":    inicio,
			"fin":       fin,
			"total":     len(resultado),
			"canciones": resultado,
			"metodo":    "BuscarRango + recorrido horizontal del sequence set (Comer)",
		})
		return
	case "danceability":
		inicio, err := leerFloatQuery(r, "inicio")
		if err != nil {
			responderError(w, http.StatusBadRequest, "parámetro inicio inválido")
			return
		}
		fin, err := leerFloatQuery(r, "fin")
		if err != nil {
			responderError(w, http.StatusBadRequest, "parámetro fin inválido")
			return
		}
		resultado := s.catalogo.BuscarRangoDanceability(inicio, fin)
		responderJSON(w, http.StatusOK, map[string]interface{}{
			"campo":     "danceability",
			"indice":    "Danceability",
			"inicio":    inicio,
			"fin":       fin,
			"total":     len(resultado),
			"canciones": resultado,
			"metodo":    "BuscarRango + recorrido horizontal del sequence set (Comer)",
		})
		return
	case "nombre", "trackname":
		desde := strings.TrimSpace(r.URL.Query().Get("inicio"))
		if desde == "" {
			desde = strings.TrimSpace(r.URL.Query().Get("desde"))
		}
		hasta := strings.TrimSpace(r.URL.Query().Get("fin"))
		if hasta == "" {
			hasta = strings.TrimSpace(r.URL.Query().Get("hasta"))
		}
		if desde == "" || hasta == "" {
			responderError(w, http.StatusBadRequest, "parámetros inicio/fin (desde/hasta) requeridos para campo nombre")
			return
		}
		resultado := s.catalogo.BuscarRangoNombre(desde, hasta)
		responderJSON(w, http.StatusOK, map[string]interface{}{
			"campo":     "nombre",
			"indice":    "TrackName",
			"desde":     desde,
			"hasta":     hasta,
			"inicio":    desde,
			"fin":       hasta,
			"total":     len(resultado),
			"canciones": resultado,
			"metodo":    "BuscarRangoString + sequence set — orden lexicográfico (Prefix B+-Tree, Comer)",
		})
		return
	}

	inicio, err := leerEnteroQuery(r, "inicio")
	if err != nil {
		responderError(w, http.StatusBadRequest, "parámetro inicio inválido")
		return
	}
	fin, err := leerEnteroQuery(r, "fin")
	if err != nil {
		responderError(w, http.StatusBadRequest, "parámetro fin inválido")
		return
	}

	resultado := s.catalogo.PorIndice.BuscarRango(inicio, fin)
	responderJSON(w, http.StatusOK, map[string]interface{}{
		"campo":     "indice",
		"indice":    "Indice",
		"total":     len(resultado),
		"canciones": resultado,
		"metodo":    "BuscarRango + recorrido horizontal del sequence set (Comer)",
	})
}

func (s *Servidor) manejarRangoNombre(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		responderError(w, http.StatusMethodNotAllowed, "use GET")
		return
	}
	if !s.verificarCatalogo(w) {
		return
	}

	desde := strings.TrimSpace(r.URL.Query().Get("desde"))
	if desde == "" {
		desde = strings.TrimSpace(r.URL.Query().Get("inicio"))
	}
	hasta := strings.TrimSpace(r.URL.Query().Get("hasta"))
	if hasta == "" {
		hasta = strings.TrimSpace(r.URL.Query().Get("fin"))
	}
	if desde == "" || hasta == "" {
		responderError(w, http.StatusBadRequest, "parámetros desde/hasta requeridos")
		return
	}

	resultado := s.catalogo.BuscarRangoNombre(desde, hasta)
	responderJSON(w, http.StatusOK, map[string]interface{}{
		"campo":     "nombre",
		"indice":    "TrackName",
		"desde":     desde,
		"hasta":     hasta,
		"total":     len(resultado),
		"canciones": resultado,
		"metodo":    "BuscarRangoString + sequence set — orden lexicográfico (Prefix B+-Tree, Comer)",
	})
}

func (s *Servidor) manejarRangoPopularidad(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		responderError(w, http.StatusMethodNotAllowed, "use GET")
		return
	}
	if !s.verificarCatalogo(w) {
		return
	}

	inicio, err := leerEnteroQuery(r, "inicio")
	if err != nil {
		responderError(w, http.StatusBadRequest, "parámetro inicio inválido")
		return
	}
	fin, err := leerEnteroQuery(r, "fin")
	if err != nil {
		responderError(w, http.StatusBadRequest, "parámetro fin inválido")
		return
	}

	resultado := s.catalogo.BuscarRangoPopularidad(inicio, fin)
	responderJSON(w, http.StatusOK, map[string]interface{}{
		"indice":    "Popularity",
		"inicio":    inicio,
		"fin":       fin,
		"total":     len(resultado),
		"canciones": resultado,
		"metodo":    "BuscarRango + recorrido horizontal del sequence set (Comer)",
	})
}

func (s *Servidor) manejarRangoTempo(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		responderError(w, http.StatusMethodNotAllowed, "use GET")
		return
	}
	if !s.verificarCatalogo(w) {
		return
	}

	inicio, err := leerFloatQuery(r, "inicio")
	if err != nil {
		responderError(w, http.StatusBadRequest, "parámetro inicio inválido")
		return
	}
	fin, err := leerFloatQuery(r, "fin")
	if err != nil {
		responderError(w, http.StatusBadRequest, "parámetro fin inválido")
		return
	}

	resultado := s.catalogo.BuscarRangoTempo(inicio, fin)
	responderJSON(w, http.StatusOK, map[string]interface{}{
		"indice":    "Tempo",
		"inicio":    inicio,
		"fin":       fin,
		"total":     len(resultado),
		"canciones": resultado,
		"metodo":    "BuscarRango + recorrido horizontal del sequence set (Comer)",
	})
}

func (s *Servidor) manejarRangoDanceability(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		responderError(w, http.StatusMethodNotAllowed, "use GET")
		return
	}
	if !s.verificarCatalogo(w) {
		return
	}

	inicio, err := leerFloatQuery(r, "inicio")
	if err != nil {
		responderError(w, http.StatusBadRequest, "parámetro inicio inválido")
		return
	}
	fin, err := leerFloatQuery(r, "fin")
	if err != nil {
		responderError(w, http.StatusBadRequest, "parámetro fin inválido")
		return
	}

	resultado := s.catalogo.BuscarRangoDanceability(inicio, fin)
	responderJSON(w, http.StatusOK, map[string]interface{}{
		"indice":    "Danceability",
		"inicio":    inicio,
		"fin":       fin,
		"total":     len(resultado),
		"canciones": resultado,
		"metodo":    "BuscarRango[float64] + sequence set (Comer)",
	})
}

func (s *Servidor) manejarInsertar(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		responderError(w, http.StatusMethodNotAllowed, "use POST")
		return
	}
	if !s.verificarCatalogo(w) {
		return
	}

	var reg bplustree.Registro
	if err := json.NewDecoder(r.Body).Decode(&reg); err != nil {
		responderError(w, http.StatusBadRequest, "JSON inválido")
		return
	}

	siguiente, err := s.db.SiguienteIndice()
	if err != nil {
		responderError(w, http.StatusInternalServerError, err.Error())
		return
	}
	reg.Indice = siguiente

	if strings.TrimSpace(reg.TrackID) == "" || strings.HasPrefix(reg.TrackID, "sim") {
		reg.TrackID = fmt.Sprintf("sim_%d", reg.Indice)
	}

	if err := s.db.InsertarEnBD(reg); err != nil {
		responderError(w, http.StatusInternalServerError, err.Error())
		return
	}

	s.catalogo.InsertarRegistro(reg)

	responderJSON(w, http.StatusOK, map[string]interface{}{
		"mensaje": "canción insertada",
		"indice":  reg.Indice,
		"cancion": reg,
	})
}

func (s *Servidor) manejarEliminar(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		responderError(w, http.StatusMethodNotAllowed, "use DELETE")
		return
	}
	if !s.verificarCatalogo(w) {
		return
	}

	indice, err := leerEnteroQuery(r, "indice")
	if err != nil {
		responderError(w, http.StatusBadRequest, "parámetro indice inválido")
		return
	}

	// Orquestador (Modo Académico): ID → Registro completo → SQLite → cascada RAM (clave + ID en secundarios).
	reg, ok := s.catalogo.PorIndice.Buscar(indice)
	if !ok {
		responderError(w, http.StatusNotFound, "canción no encontrada en arbolPorIndice")
		return
	}

	if err := s.db.EliminarDeBD(indice); err != nil {
		responderError(w, http.StatusNotFound, err.Error())
		return
	}

	// EliminarRegistro: PorIndice.Eliminar(id) + EliminarExacto(clave, id) en cada secundario.
	s.catalogo.EliminarRegistro(*reg)

	responderJSON(w, http.StatusOK, map[string]interface{}{
		"mensaje": "canción eliminada de BD y de los 5 índices B+ en RAM",
		"indice":  indice,
		"cancion": reg,
		"cascada": map[string]interface{}{
			"indice":        reg.Indice,
			"trackName":     reg.TrackName,
			"popularity":    reg.Popularity,
			"tempo":         reg.Tempo,
			"danceability":  reg.Danceability,
			"metodoExacto":  "EliminarExacto(claveSecundaria, reg.Indice) con sequence set",
		},
		"metodo": "PorIndice.Buscar → EliminarDeBD → EliminarRegistro (cascada con ID en secundarios)",
	})
}

func (s *Servidor) manejarSiguienteIndice(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		responderError(w, http.StatusMethodNotAllowed, "use GET")
		return
	}

	indice, err := s.db.SiguienteIndice()
	if err != nil {
		responderError(w, http.StatusInternalServerError, err.Error())
		return
	}

	responderJSON(w, http.StatusOK, map[string]int{"indice": indice})
}

func (s *Servidor) manejarConteoBD(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		responderError(w, http.StatusMethodNotAllowed, "use GET")
		return
	}
	total, err := s.db.ContarTracks()
	if err != nil {
		responderError(w, http.StatusInternalServerError, err.Error())
		return
	}
	responderJSON(w, http.StatusOK, map[string]int{"total": total})
}

func (s *Servidor) manejarEstadisticas(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		responderError(w, http.StatusMethodNotAllowed, "use GET")
		return
	}
	if !s.verificarCatalogo(w) {
		return
	}
	tipo := strings.TrimSpace(r.URL.Query().Get("tipo"))
	if tipo == "" {
		tipo = "indice"
	}
	stats, err := s.catalogo.EstadisticasPorTipo(tipo)
	if err != nil {
		responderError(w, http.StatusBadRequest, err.Error())
		return
	}
	responderJSON(w, http.StatusOK, stats)
}

func (s *Servidor) manejarEstructura(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		responderError(w, http.StatusMethodNotAllowed, "use GET")
		return
	}
	if !s.verificarCatalogo(w) {
		return
	}

	tipo := strings.TrimSpace(r.URL.Query().Get("tipo"))
	exp, err := s.catalogo.ExportarPorTipo(tipo)
	if err != nil {
		responderError(w, http.StatusBadRequest, err.Error())
		return
	}

	responderJSON(w, http.StatusOK, exp)
}

func (s *Servidor) manejarIndices(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		responderError(w, http.StatusMethodNotAllowed, "use GET")
		return
	}
	if !s.verificarCatalogo(w) {
		return
	}

	responderJSON(w, http.StatusOK, map[string]interface{}{
		"indices": s.nombresIndices(),
		"orden":   s.catalogo.Orden,
	})
}

func (s *Servidor) nombresIndices() []map[string]string {
	return []map[string]string{
		{"campo": "Indice", "tipo": "int", "uso": "búsqueda exacta y range scan por id"},
		{"campo": "TrackName", "tipo": "string", "uso": "Prefix B+-Tree: búsqueda por prefijo y range scan lexicográfico (BuscarRangoString)"},
		{"campo": "Popularity", "tipo": "int", "uso": "range scan — top canciones"},
		{"campo": "Tempo", "tipo": "float64", "uso": "range scan — filtro BPM para DJs"},
		{"campo": "Danceability", "tipo": "float64", "uso": "range scan — mezcla por bailabilidad"},
	}
}

func (s *Servidor) verificarCatalogo(w http.ResponseWriter) bool {
	if s.catalogo == nil {
		responderError(w, http.StatusServiceUnavailable, "índices no inicializados, use POST /api/inicializar")
		return false
	}
	return true
}

func (s *Servidor) habilitarCORS(siguiente http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}
		siguiente(w, r)
	}
}

func leerEnteroQuery(r *http.Request, nombre string) (int, error) {
	valor := r.URL.Query().Get(nombre)
	return strconv.Atoi(valor)
}

func leerFloatQuery(r *http.Request, nombre string) (float64, error) {
	valor := r.URL.Query().Get(nombre)
	return strconv.ParseFloat(valor, 64)
}

func responderJSON(w http.ResponseWriter, codigo int, datos interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(codigo)
	json.NewEncoder(w).Encode(datos)
}

func responderError(w http.ResponseWriter, codigo int, mensaje string) {
	responderJSON(w, codigo, map[string]string{"error": mensaje})
}
