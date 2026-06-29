# Carpeta `bplustree`

Aquí está el **árbol B+** del proyecto. Todos los archivos comparten el mismo paquete Go (`package bplustree`).

## Código de producción

| Archivo | Qué hace |
|---------|----------|
| `tipos.go` | Define las estructuras del árbol (`Tree`, `nodo`, `entrada`) y `NewTree` |
| `buscar.go` | Lectura: `Buscar`, `BuscarRango`, `BuscarExactos` |
| `insertar.go` | Escritura: `Insertar` y splits cuando una hoja se llena |
| `eliminar.go` | Borrado: `Eliminar`, `EliminarSi` y balanceo (préstamo/fusión) |
| `interno.go` | Funciones privadas que bajan por el árbol (`irAHoja`, comparadores, rutas). No las llama el usuario directamente |
| `registro.go` | Modelo de canción Spotify y los 5 índices (`CatalogoIndices`) |
| `prefijo.go` | Búsqueda por nombre: prefijo y rango de letras |
| `serializar.go` | Convierte el árbol a JSON para el canvas de Vue |
| `estadisticas.go` | Conteos del árbol cuando hay muchas canciones y no se dibuja el canvas |

## Tests (mismo paquete, archivos `*_test.go`)

| Archivo | Qué prueba |
|---------|------------|
| `buscar_test.go` | Búsqueda exacta, rango y sequence set |
| `insertar_test.go` | Inserción, duplicados y splits |
| `eliminar_test.go` | Eliminación, underflow y fusión |
| `bench_test.go` | Benchmarks de rendimiento |
| `catalogo_test.go` | Los 5 índices Spotify juntos |

## Comandos

```bash
go test ./bplustree/ -v          # tests
go test -bench=. ./bplustree/    # benchmarks
```
