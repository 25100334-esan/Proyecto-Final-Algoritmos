package bplustree_test

import (
	"fmt"
	"testing"

	"bplustree-proyecto/bplustree"
)

func TestIndicesSecundariosGenericos(t *testing.T) {
	cat := bplustree.NuevoCatalogo(4)

	reg := bplustree.Registro{
		Indice: 1, TrackName: "Comedy", Artists: "Gen Hoshino",
		Popularity: 73, Tempo: 87.9, Danceability: 0.67,
	}
	cat.InsertarRegistro(reg)

	if _, ok := cat.PorIndice.Buscar(1); !ok {
		t.Fatal("índice por id falló")
	}

	valores, ok := cat.PorNombre.Buscar("Comedy")
	if !ok {
		t.Fatal("índice por nombre falló")
	}
	if valores.Indice != 1 {
		t.Fatalf("nombre devolvió id %d", valores.Indice)
	}

	rango := cat.PorPopularidad.BuscarRango(70, 80)
	if len(rango) != 1 || rango[0].TrackName != "Comedy" {
		t.Fatalf("rango popularidad incorrecto: %+v", rango)
	}

	rangoTempo := cat.PorTempo.BuscarRango(85.0, 90.0)
	if len(rangoTempo) != 1 {
		t.Fatalf("rango tempo incorrecto: %d", len(rangoTempo))
	}

	cat.EliminarPorIndice(1)
	if _, ok := cat.PorIndice.Buscar(1); ok {
		t.Fatal("eliminar no sincronizó índice principal")
	}
	if _, ok := cat.PorNombre.Buscar("Comedy"); ok {
		t.Fatal("eliminar no sincronizó índice por nombre")
	}
}

func TestClavesDuplicadasTempo(t *testing.T) {
	cat := bplustree.NuevoCatalogo(2)
	tempo := 120.0

	for i := 0; i < 5; i++ {
		cat.InsertarRegistro(bplustree.Registro{
			Indice: i, TrackName: fmt.Sprintf("Track-%d", i),
			Tempo: tempo, Popularity: 50 + i,
		})
	}

	exactos := cat.PorTempo.BuscarExactos(tempo)
	if len(exactos) != 5 {
		t.Fatalf("BuscarExactos(tempo) devolvió %d registros, esperaba 5", len(exactos))
	}

	rango := cat.PorTempo.BuscarRango(tempo, tempo)
	if len(rango) != 5 {
		t.Fatalf("BuscarRango(%v,%v) devolvió %d registros, esperaba 5", tempo, tempo, len(rango))
	}

	ids := map[int]bool{}
	for _, r := range exactos {
		ids[r.Indice] = true
	}
	for i := 0; i < 5; i++ {
		if !ids[i] {
			t.Fatalf("falta canción con id %d en resultados de tempo duplicado", i)
		}
	}
}

func TestEliminarExactoTempoDuplicado(t *testing.T) {
	cat := bplustree.NuevoCatalogo(2)
	tempo := 120.0

	for i := 0; i < 5; i++ {
		cat.InsertarRegistro(bplustree.Registro{
			Indice: i, TrackName: fmt.Sprintf("Track-%d", i),
			Tempo: tempo, Popularity: 50 + i,
		})
	}

	if !bplustree.EliminarExacto(cat.PorTempo, tempo, 3) {
		t.Fatal("EliminarExacto(tempo, id=3) debe devolver true")
	}
	if len(cat.PorTempo.BuscarExactos(tempo)) != 4 {
		t.Fatal("debe quedar 4 registros con tempo 120.0")
	}
	for _, r := range cat.PorTempo.BuscarExactos(tempo) {
		if r.Indice == 3 {
			t.Fatal("id 3 no debe permanecer en índice tempo")
		}
	}
	if _, ok := cat.PorIndice.Buscar(3); !ok {
		t.Fatal("id 3 debe seguir en índice primario tras EliminarExacto solo en tempo")
	}

	if !cat.EliminarPorIndice(3) {
		t.Fatal("EliminarPorIndice(3) debe eliminar en cascada")
	}
	if len(cat.PorTempo.BuscarExactos(tempo)) != 4 {
		t.Fatalf("cascada debe dejar 4 en tempo, obtuvo %d", len(cat.PorTempo.BuscarExactos(tempo)))
	}
}

func TestBuscarPorPrefijoNombre(t *testing.T) {
	cat := bplustree.NuevoCatalogo(2)
	nombres := []string{"Comedy", "Comedown", "Comedy Central", "Comet", "Fade"}
	for i, nombre := range nombres {
		cat.InsertarRegistro(bplustree.Registro{
			Indice: i, TrackName: nombre, Tempo: 100,
		})
	}

	resultado := cat.BuscarPorPrefijoNombre("comed")
	if len(resultado) != 3 {
		t.Fatalf("BuscarPorPrefijo(comed) devolvió %d, esperaba 3: %+v", len(resultado), nombresDe(resultado))
	}
	// Orden alfabético en hojas (sequence set): Comedown < Comedy < Comedy Central
	if resultado[0].TrackName != "Comedown" || resultado[1].TrackName != "Comedy" || resultado[2].TrackName != "Comedy Central" {
		t.Fatalf("orden o contenido incorrecto: %+v", nombresDe(resultado))
	}

	if len(cat.BuscarPorPrefijoNombre("Comet")) != 1 {
		t.Fatal("prefijo Comet debe devolver solo Comet")
	}
	if len(cat.BuscarPorPrefijoNombre("fade")) != 1 {
		t.Fatal("prefijo fade debe ser case-insensitive")
	}
}

func TestBuscarPorPrefijoTrasClavesMayores(t *testing.T) {
	// Regresión: hoja con ultima clave >= prefijo pero sin coincidencias reales (p. ej. "texas" > "test").
	cat := bplustree.NuevoCatalogo(2)
	nombres := []string{"texas", "text", "thank you", "testing"}
	for i, nombre := range nombres {
		cat.InsertarRegistro(bplustree.Registro{Indice: i, TrackName: nombre, Tempo: 100})
	}
	cat.InsertarRegistro(bplustree.Registro{Indice: 99, TrackName: "test", Tempo: 120})

	if len(cat.BuscarPorPrefijoNombre("t")) < 1 {
		t.Fatal("prefijo t debe incluir test")
	}
	resultado := cat.BuscarPorPrefijoNombre("test")
	nombresTest := nombresDe(resultado)
	if len(resultado) == 0 || nombresTest[0] != "test" {
		t.Fatalf("prefijo test debe incluir la canción test, obtuvo: %+v", nombresTest)
	}
	if len(cat.BuscarPorPrefijoNombre("testing")) != 1 {
		t.Fatal("prefijo testing debe devolver 1")
	}
}

func TestBuscarRangoNombreLetras(t *testing.T) {
	cat := bplustree.NuevoCatalogo(2)
	canciones := []string{"Apple", "Macarena", "Mamma Mia", "November Rain", "Oasis", "Paint it Black", "Zebra"}
	for i, nombre := range canciones {
		cat.InsertarRegistro(bplustree.Registro{Indice: i, TrackName: nombre, Tempo: 100})
	}

	rango := cat.BuscarRangoNombre("M", "O")
	nombres := nombresDe(rango)
	if len(nombres) != 4 {
		t.Fatalf("rango M-O esperaba 4 (Macarena, Mamma Mia, November Rain, Oasis), obtuvo %d: %v", len(nombres), nombres)
	}
	for _, n := range nombres {
		if n == "Apple" || n == "Paint it Black" || n == "Zebra" {
			t.Fatalf("canción fuera de rango M-O: %s", n)
		}
	}

	// Prefijo «comed» = range scan acotado (Prefix B+-Tree)
	cat2 := bplustree.NuevoCatalogo(2)
	for i, nombre := range []string{"Comedy", "Comedown", "Comet"} {
		cat2.InsertarRegistro(bplustree.Registro{Indice: i, TrackName: nombre, Tempo: 100})
	}
	if len(cat2.BuscarRangoNombre("Comed", "Comedy")) != 2 {
		t.Fatalf("rango Comed-Comedy esperaba Comedown y Comedy")
	}
}

func TestBuscarRangoDanceability(t *testing.T) {
	cat := bplustree.NuevoCatalogo(2)
	datos := []struct {
		id   int
		name string
		d    float64
	}{
		{0, "Low", 0.3},
		{1, "Mid", 0.65},
		{2, "High", 0.85},
	}
	for _, d := range datos {
		cat.InsertarRegistro(bplustree.Registro{
			Indice: d.id, TrackName: d.name, Danceability: d.d, Tempo: 100,
		})
	}

	rango := cat.BuscarRangoDanceability(0.5, 0.8)
	if len(rango) != 1 || rango[0].TrackName != "Mid" {
		t.Fatalf("rango danceability [0.5,0.8] esperaba Mid, obtuvo: %+v", nombresDe(rango))
	}
}

func TestEliminarExactoNombreDuplicado(t *testing.T) {
	cat := bplustree.NuevoCatalogo(2)
	nombre := "Shape of You"

	for i := 0; i < 3; i++ {
		cat.InsertarRegistro(bplustree.Registro{
			Indice: i, TrackName: nombre, Tempo: 100 + float64(i),
		})
	}

	if !bplustree.EliminarExacto(cat.PorNombre, nombre, 1) {
		t.Fatal("EliminarExacto(nombre, id=1) debe devolver true")
	}
	resto := cat.PorNombre.BuscarExactos(nombre)
	if len(resto) != 2 {
		t.Fatalf("deben quedar 2 con mismo nombre, obtuvo %d", len(resto))
	}
	for _, r := range resto {
		if r.Indice == 1 {
			t.Fatal("id 1 no debe permanecer en índice nombre")
		}
	}
}

func TestEliminarCascadaIndices(t *testing.T) {
	cat := bplustree.NuevoCatalogo(2)
	reg := bplustree.Registro{
		Indice: 150, TrackName: "Shape of You", Popularity: 80,
		Tempo: 95.5, Danceability: 0.75,
	}
	cat.InsertarRegistro(reg)

	if !cat.EliminarPorIndice(150) {
		t.Fatal("EliminarPorIndice debe devolver true")
	}
	if _, ok := cat.PorIndice.Buscar(150); ok {
		t.Fatal("debe eliminar del índice principal")
	}
	if _, ok := cat.PorNombre.Buscar("Shape of You"); ok {
		t.Fatal("debe eliminar del índice por nombre")
	}
	if len(cat.PorTempo.BuscarExactos(95.5)) != 0 {
		t.Fatal("debe eliminar del índice por tempo")
	}
	if cat.EliminarPorIndice(999) {
		t.Fatal("EliminarPorIndice inexistente debe devolver false")
	}
}

func nombresDe(regs []bplustree.Registro) []string {
	out := make([]string, len(regs))
	for i, r := range regs {
		out[i] = r.TrackName
	}
	return out
}
