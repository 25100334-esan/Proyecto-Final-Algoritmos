# Interacción: Modo Negocio — Consultas

Vista: `BusinessMode.vue` — UI estilo Spotify sin canvas de árbol.

## Componentes

```mermaid
flowchart TB
    BM["BusinessMode.vue"] --> FS["FilterSidebar.vue"]
    BM --> SG["SongGridCard.vue"]
    BM --> SC["SongCard.vue"]
    BM --> PV["PlaylistView.vue"]
    BM --> DC["DebugConsole.vue"]
    BM --> PL["SpotifyLogo.vue"]
    BM --> PERF["perfLog.js"]
    BM --> API["bplustreeApi.js"]
```

## Flujos de consulta

```mermaid
flowchart TD
    subgraph Busqueda
        Q["searchQuery"] --> DEB["debounce 280ms"]
        DEB --> PF["api.buscarPorPrefijo"]
    end

    subgraph Filtros
        POP["Slider popularidad"] --> RP["api.buscarRangoPopularidad"]
        TEM["Slider tempo"] --> RT["api.buscarRangoTempo"]
        NOM["Rango nombre A-Z"] --> RN["api.buscarRangoNombre"]
    end

    subgraph ID
        IDIN["Input ID"] --> BID["api.buscar"]
    end

    PF & RP & RT & RN & BID --> GRID["SongGridCard lista"]
```

## Secuencia autocompletado por prefijo

```mermaid
sequenceDiagram
    participant U as Usuario
    participant B as BusinessMode
    participant A as bplustreeApi
    participant H as manejarBuscarPrefijo
    participant P as BuscarPorPrefijo

    U->>B: Escribe "ghost"
    B->>B: debounce 280ms
    B->>A: GET /api/buscar-prefijo?prefijo=ghost
    A->>H: prefijo query
    H->>P: BuscarPorPrefijoNombre
    P-->>H: []Registro
    H-->>B: JSON resultados
    B->>B: perfLog tiempo ms
    B->>B: Actualizar grid
```

## Comparación Modo Académico vs Negocio

```mermaid
flowchart LR
    subgraph ACAD["Modo Académico"]
        A1["Visualiza árbol B+"]
        A2["Animaciones paso a paso"]
        A3["Operaciones CRUD demo"]
    end

    subgraph NEG["Modo Negocio"]
        N1["UI producto Spotify"]
        N2["Prefijo + filtros"]
        N3["Consola debug rendimiento"]
        N4["Sin canvas árbol"]
    end

    BOTH["Comparten"] --> API["misma API REST"]
    BOTH --> CAT["mismo CatalogoIndices RAM"]
```
