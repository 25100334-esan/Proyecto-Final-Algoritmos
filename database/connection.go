package database

import (
	"database/sql"
	"fmt"

	"bplustree-proyecto/bplustree"

	_ "modernc.org/sqlite"
)

// DB mantiene la conexión al archivo dataset.db en disco.
type DB struct {
	conexion *sql.DB
}

// Conectar abre el archivo SQLite. Complejidad: O(1) temporal, O(1) espacial.
func Conectar(ruta string) (*DB, error) {
	conexion, err := sql.Open("sqlite", ruta)
	if err != nil {
		return nil, fmt.Errorf("abrir base de datos: %w", err)
	}
	if err := conexion.Ping(); err != nil {
		return nil, fmt.Errorf("conectar base de datos: %w", err)
	}
	return &DB{conexion: conexion}, nil
}

// Cerrar libera la conexión al disco. Complejidad: O(1).
func (db *DB) Cerrar() error {
	return db.conexion.Close()
}

// CargarDatosEnArbol lee hasta limite canciones del disco y las inserta en el árbol en RAM.
// Complejidad: O(n log n) temporal, O(n) espacial.
func (db *DB) CargarDatosEnArbol(arbol *bplustree.BPlusTree, limite int) (int, error) {
	return db.cargarTracks(limite, func(reg bplustree.Registro) {
		arbol.Insertar(reg)
	})
}

// CargarDatosEnCatalogo inserta cada fila en todos los índices secundarios del catálogo.
func (db *DB) CargarDatosEnCatalogo(cat *bplustree.CatalogoIndices, limite int) (int, error) {
	return db.cargarTracks(limite, cat.InsertarRegistro)
}

func (db *DB) cargarTracks(limite int, insertar func(bplustree.Registro)) (int, error) {
	consulta := `
		SELECT id, track_id, artists, album_name, track_name, popularity,
		       duration_ms, explicit, danceability, energy, musical_key, loudness,
		       musical_mode, speechiness, acousticness, instrumentalness, liveness,
		       valence, tempo, time_signature, track_genre
		FROM tracks
		ORDER BY id
		LIMIT ?`

	filas, err := db.conexion.Query(consulta, limite)
	if err != nil {
		return 0, fmt.Errorf("consultar tracks: %w", err)
	}
	defer filas.Close()

	cargados := 0
	for filas.Next() {
		reg, err := leerRegistro(filas)
		if err != nil {
			return cargados, fmt.Errorf("leer fila %d: %w", cargados, err)
		}
		insertar(reg)
		cargados++
	}
	if err := filas.Err(); err != nil {
		return cargados, err
	}

	return cargados, nil
}

// ListarTracks devuelve hasta limite canciones ordenadas por id (sin insertar en el árbol).
func (db *DB) ListarTracks(limite int) ([]bplustree.Registro, error) {
	consulta := `
		SELECT id, track_id, artists, album_name, track_name, popularity,
		       duration_ms, explicit, danceability, energy, musical_key, loudness,
		       musical_mode, speechiness, acousticness, instrumentalness, liveness,
		       valence, tempo, time_signature, track_genre
		FROM tracks
		ORDER BY id
		LIMIT ?`

	filas, err := db.conexion.Query(consulta, limite)
	if err != nil {
		return nil, fmt.Errorf("consultar tracks: %w", err)
	}
	defer filas.Close()

	var lista []bplustree.Registro
	for filas.Next() {
		reg, err := leerRegistro(filas)
		if err != nil {
			return lista, fmt.Errorf("leer fila: %w", err)
		}
		lista = append(lista, reg)
	}
	return lista, filas.Err()
}

// ContarTracks devuelve el total de filas en la tabla tracks.
func (db *DB) ContarTracks() (int, error) {
	var total int
	err := db.conexion.QueryRow(`SELECT COUNT(*) FROM tracks`).Scan(&total)
	if err != nil {
		return 0, fmt.Errorf("contar tracks: %w", err)
	}
	return total, nil
}

// SiguienteIndice devuelve el próximo id libre en tracks (MAX(id) + 1).
func (db *DB) SiguienteIndice() (int, error) {
	var maxID sql.NullInt64
	err := db.conexion.QueryRow(`SELECT MAX(id) FROM tracks`).Scan(&maxID)
	if err != nil {
		return 0, fmt.Errorf("obtener siguiente índice: %w", err)
	}
	if !maxID.Valid {
		return 1, nil
	}
	return int(maxID.Int64) + 1, nil
}

// InsertarEnBD guarda una canción en el disco. Complejidad: O(1) temporal, O(1) espacial.
func (db *DB) InsertarEnBD(reg bplustree.Registro) error {
	_, err := db.conexion.Exec(`
		INSERT INTO tracks (
			id, track_id, artists, album_name, track_name, popularity,
			duration_ms, explicit, danceability, energy, musical_key, loudness,
			musical_mode, speechiness, acousticness, instrumentalness, liveness,
			valence, tempo, time_signature, track_genre
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		reg.Indice, reg.TrackID, reg.Artists, reg.AlbumName, reg.TrackName, reg.Popularity,
		reg.DurationMs, reg.Explicit, reg.Danceability, reg.Energy, reg.Key, reg.Loudness,
		reg.Mode, reg.Speechiness, reg.Acousticness, reg.Instrumentalness, reg.Liveness,
		reg.Valence, reg.Tempo, reg.TimeSignature, reg.TrackGenre,
	)
	if err != nil {
		return fmt.Errorf("insertar en tracks: %w", err)
	}
	return nil
}

// EliminarDeBD borra una canción del disco por su Indice (campo id). Complejidad: O(1).
func (db *DB) EliminarDeBD(indice int) error {
	resultado, err := db.conexion.Exec(`DELETE FROM tracks WHERE id = ?`, indice)
	if err != nil {
		return fmt.Errorf("eliminar de tracks: %w", err)
	}
	afectados, err := resultado.RowsAffected()
	if err != nil {
		return err
	}
	if afectados == 0 {
		return fmt.Errorf("no existe canción con índice %d", indice)
	}
	return nil
}

// leerRegistro convierte una fila SQL en un Registro del árbol.
func leerRegistro(filas *sql.Rows) (bplustree.Registro, error) {
	var reg bplustree.Registro
	var album, genero sql.NullString
	var clave, modo, compas sql.NullInt64
	var explicito int

	err := filas.Scan(
		&reg.Indice, &reg.TrackID, &reg.Artists, &album, &reg.TrackName, &reg.Popularity,
		&reg.DurationMs, &explicito, &reg.Danceability, &reg.Energy, &clave, &reg.Loudness,
		&modo, &reg.Speechiness, &reg.Acousticness, &reg.Instrumentalness, &reg.Liveness,
		&reg.Valence, &reg.Tempo, &compas, &genero,
	)
	if err != nil {
		return reg, err
	}

	reg.Explicit = explicito == 1
	reg.AlbumName = album.String
	reg.TrackGenre = genero.String
	if clave.Valid {
		reg.Key = int(clave.Int64)
	}
	if modo.Valid {
		reg.Mode = int(modo.Int64)
	}
	if compas.Valid {
		reg.TimeSignature = int(compas.Int64)
	}

	return reg, nil
}
