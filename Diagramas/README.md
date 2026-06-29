# Diagramas del Proyecto B+ Tree

Diagramas en formato **Mermaid** listos para [Mermaid Live Editor](https://mermaid.live).

## Cómo usarlos

1. Abre https://mermaid.live
2. Copia el bloque `mermaid` del archivo que necesites
3. Pégalo en el editor y exporta PNG/SVG si lo requieres

## Índice

### Estructura y arquitectura

| Archivo | Contenido |
|---------|-----------|
| [01-repositorio-estructura.md](./01-repositorio-estructura.md) | Árbol de carpetas y archivos del repo |
| [02-arquitectura-sistema.md](./02-arquitectura-sistema.md) | Capas: Vue → API Go → B+ Tree → SQLite |
| [03-catalogo-5-indices.md](./03-catalogo-5-indices.md) | Catálogo `CatalogoIndices` y sus 5 árboles |
| [04-interaccion-completa.md](./04-interaccion-completa.md) | Vista unificada de todas las capas |

### Operaciones básicas

| Archivo | Operación | Complejidad |
|---------|-----------|-------------|
| [operaciones/buscar.md](./operaciones/buscar.md) | `Tree.Buscar` / `BuscarExactos` | O(log d n) |
| [operaciones/insertar.md](./operaciones/insertar.md) | `Tree.Insertar` + split | O(log d n) |
| [operaciones/eliminar.md](./operaciones/eliminar.md) | `Eliminar` / `EliminarExacto` | O(log d n) |
| [operaciones/range-scan.md](./operaciones/range-scan.md) | `BuscarRango` + sequence set | O(log d n + k) |
| [operaciones/prefijo.md](./operaciones/prefijo.md) | `BuscarPorPrefijo` | O(log d n + k) |

### Subfunciones internas

| Archivo | Tema |
|---------|------|
| [subfunciones/navegacion-interna.md](./subfunciones/navegacion-interna.md) | `irAHoja`, `indiceHijo`, `bajarGuardandoRuta` |
| [subfunciones/split-overflow.md](./subfunciones/split-overflow.md) | `partirHoja`, `partirNodoInterno`, promoción |
| [subfunciones/underflow-prestamo-fusion.md](./subfunciones/underflow-prestamo-fusion.md) | Préstamo, fusión, `ajustarRaiz` |
| [subfunciones/base-datos.md](./subfunciones/base-datos.md) | SQLite: carga, insert, delete |
| [subfunciones/serializacion-estadisticas.md](./subfunciones/serializacion-estadisticas.md) | JSON canvas + estadísticas |

### Interacción entre capas

| Archivo | Tema |
|---------|------|
| [interaccion/flujo-general-api.md](./interaccion/flujo-general-api.md) | Rutas REST y handlers |
| [interaccion/flujo-insertar-completo.md](./interaccion/flujo-insertar-completo.md) | POST /api/insertar end-to-end |
| [interaccion/flujo-eliminar-cascada.md](./interaccion/flujo-eliminar-cascada.md) | DELETE /api/eliminar en 5 índices |
| [interaccion/modo-academico-animaciones.md](./interaccion/modo-academico-animaciones.md) | Canvas + animaciones paso a paso |
| [interaccion/modo-negocio-consultas.md](./interaccion/modo-negocio-consultas.md) | UI Spotify + prefijo/filtros |
