# Proyecto Final — Árbol B+ (Spotify)

**Curso:** Algoritmos y Estructura de Datos — ESAN 2026-1  
**Estructura:** B+-Tree genérico (Comer, d≥2) + 5 índices secundarios en RAM  
**Dataset:** Catálogo Spotify (~114 000 canciones en SQLite)

## ¿Qué hace este proyecto?

Simula un **motor de indexación estilo base de datos** sobre el catálogo de Spotify usando **árboles B+** en memoria. Las canciones se cargan desde `data/dataset.db` y se indexan en cinco árboles sincronizados:

1. **ID** (índice primario, claves únicas)
2. **TrackName** (string, con prefix search y rango lexicográfico)
3. **Popularity** (int)
4. **Tempo** (float64, BPM)
5. **Danceability** (float64)

El backend expone una **API REST** en Go (`:8080`) y el frontend en **Vue 3** ofrece dos modos:

- **Modo académico:** canvas interactivo del árbol con animaciones paso a paso (insertar, buscar, split, underflow, fusión, eliminar).
- **Modo negocio:** interfaz estilo Spotify con filtros, range scan y consola de depuración.

Cuando hay más de **500 canciones**, el canvas se desactiva automáticamente y se muestran **estadísticas** del árbol (altura, nodos, entradas).

## Integrantes

| Nombre | GitHub | Rol |
|--------|--------|-----|
| Brayan Paredes Sucapuca 25100334 | [25100334-esan](https://github.com/25100334-esan) | Coordinación, API, modo negocio |
| Marcelo Miranda | [MarceloW18](https://github.com/MarceloW18) | B+-Tree genérico, índices secundarios |
| Elvin Loarte Quinteros 25100612| [Leo-nardo-14](https://github.com/Leo-nardo-14) | Frontend académico, animaciones |
| Aaron Cruzado | [aaroncruzado60-byte](https://github.com/aaroncruzado60-byte) | Tests, prefix search, integración |

## Estructura del repositorio

```
├── main.go                 # Punto de entrada del backend (:8080)
├── bplustree/              # Núcleo del B+-Tree (ver bplustree/README.md)
│   ├── tipos.go            # Estructuras + NewTree
│   ├── buscar.go           # Buscar, BuscarRango
│   ├── insertar.go         # Insertar + split
│   ├── eliminar.go         # Eliminar + underflow
│   ├── interno.go          # Recorrido interno del árbol
│   ├── registro.go         # Canción Spotify + catálogo
│   ├── prefijo.go          # Búsqueda por nombre
│   ├── serializar.go       # JSON para Vue
│   └── estadisticas.go     # Métricas sin canvas
├── api/                    # Handlers REST
├── database/               # Conexión SQLite (dataset.db)
├── data/                   # dataset.db local (no en Git)
└── frontend/               # Vue 3 + Vite
```

## Requisitos

- Go 1.21+
- Node.js 18+
- Archivo `data/dataset.db` (ver `data/README.md`)

## Ejecución

```bash
# Terminal 1 — Backend
go run .

# Terminal 2 — Frontend
cd frontend
npm install
npm run dev
```

Abre el frontend, elige el orden `d` del árbol, inicializa con un límite de canciones y explora el modo académico o negocio.

## Operaciones implementadas

| Operación | Complejidad | Descripción |
|-----------|-------------|-------------|
| `Buscar` | O(log_d n) | Búsqueda exacta por clave |
| `Insertar` | O(log_d n) | Inserción con split automático |
| `Eliminar` | O(log_d n) | Eliminación con préstamo/fusión |
| `BuscarRango` | O(log_d n + k) | Range scan por sequence set |
| `BuscarPorPrefijo` | O(log_d n + k) | Autocompletado por nombre |
| `EliminarExacto` | O(log_d n) | Elimina duplicado por clave + ID |
| Índices secundarios | — | 5 árboles sincronizados en RAM |

## Tests

```bash
# Todos los tests del proyecto
go test ./...

# Solo el paquete bplustree (con detalle)
go test ./bplustree/ -v
```

### Cobertura de tests

| Archivo | Qué prueba |
|---------|------------|
| `buscar_test.go` | Búsqueda exacta, rango, sequence set |
| `insertar_test.go` | Inserción, duplicados, splits |
| `eliminar_test.go` | Eliminación, underflow, fusión |
| `bench_test.go` | Benchmarks de rendimiento |
| `catalogo_test.go` | Los 5 índices Spotify |

## Benchmarks

```bash
go test -bench=. -benchmem ./bplustree/
```

| Benchmark | Qué mide |
|-----------|----------|
| `BenchmarkInsertar` | 10 000 inserciones secuenciales |
| `BenchmarkBuscar` | Búsquedas aleatorias en árbol poblado |
| `BenchmarkBuscarRango` | Range scan de 100 resultados |

## API REST (resumen)

| Método | Ruta | Acción |
|--------|------|--------|
| POST | `/api/inicializar` | Carga canciones desde SQLite al catálogo |
| GET | `/api/estructura` | JSON del árbol para el canvas |
| GET | `/api/estadisticas` | Métricas cuando canvas está desactivado |
| POST | `/api/insertar` | Inserta una canción |
| POST | `/api/buscar` | Búsqueda exacta |
| POST | `/api/eliminar` | Elimina por ID |
| POST | `/api/rango` | Range scan (popularidad, tempo, etc.) |
| POST | `/api/prefijo` | Prefix search por nombre |

## Referencias

- Comer (1979) — *The Ubiquitous B-Tree*
- Bayer & McCreight (1972) — B-Tree original
- Material del curso: *El Árbol B+ (B+-Tree).pdf*

## Repositorio

https://github.com/25100334-esan/Proyecto-Final-Algoritmos
