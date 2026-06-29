# Operación: Range Scan (BuscarRango)

**API:** `Tree[K,V].BuscarRango(desde, hasta)` — O(log d n + k)

```mermaid
flowchart TD
    START(["BuscarRango(desde, hasta)"]) --> V{desde > hasta?}
    V -->|Sí| NIL["return nil"]
    V -->|No| H["hoja = irAHoja(desde)"]
    H --> LOOP["Para cada hoja en sequence set"]
    LOOP --> ENTR["Para cada entrada en hoja"]
    ENTR --> LT{clave < desde?}
    LT -->|Sí| ENTR
    LT -->|No| GT{clave > hasta?}
    GT -->|Sí| RET(["return resultado"])
    GT -->|No| ADD["append valor"]
    ADD --> ENTR
    ENTR -->|fin hoja| UL{última clave > hasta?}
    UL -->|Sí| RET
    UL -->|No| NEXT["hoja = hoja.siguienteHoja"]
    NEXT --> LOOP
```

## Sequence set (enlace horizontal entre hojas)

```mermaid
flowchart LR
    H1["Hoja 1<br/>k1 k2 k3"] -->|siguienteHoja| H2["Hoja 2<br/>k4 k5"]
    H2 -->|siguienteHoja| H3["Hoja 3<br/>k6 k7 k8"]
```

## Range scan por tipo de índice

```mermaid
flowchart TB
    API["GET /api/rango?campo=&inicio=&fin="] --> SW{campo}
    SW -->|indice| I["PorIndice.BuscarRango"]
    SW -->|popularidad| P["BuscarRangoPopularidad"]
    SW -->|tempo| T["BuscarRangoTempo"]
    SW -->|danceability| D["BuscarRangoDanceability"]
    SW -->|nombre| N["BuscarRangoNombre"]
    N --> BS["BuscarRangoString"]
    I & P & T & D --> BR["Tree.BuscarRango genérico"]
    BS --> SEMI["Semiabierto o inclusivo<br/>según longitud del rango"]
```

## BuscarRangoString — lógica de strings

```mermaid
flowchart TD
    BRS["BuscarRangoString(desde, hasta)"] --> N["Normalizar ToLower"]
    N --> L1{desde y hasta<br/>son 1 letra?}
    L1 -->|Sí| SEMI["buscarRangoStringSemi<br/>[M, limiteSuperior O)"]
    L1 -->|No| INC["buscarRangoStringInclusivo<br/>[desde, hasta]"]
    SEMI --> IRP["irAHojaPrefijo + recorrido SS"]
    INC --> IRP
```
