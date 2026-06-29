# Operación: Búsqueda por Prefijo

**API:** `BuscarPorPrefijo(t, prefijo)` — O(log d n + k)

Implementación como range scan semiabierto sobre `PorNombre`.

```mermaid
flowchart TD
    START(["BuscarPorPrefijo(prefijo)"]) --> T["Trim + validar no vacío"]
    T --> N["prefijoNorm = ToLower(prefijo)"]
    N --> LIM["hasta = limiteSuperiorPrefijo(prefijoNorm)"]
    LIM --> SEMI["buscarRangoStringSemi(t, prefijoNorm, hasta)"]
    SEMI --> IRP["irAHojaPrefijo(t, prefijoNorm)"]
    IRP --> LOOP["Recorrer hojas vía sequence set"]
    LOOP --> CMP{claveNorm >= prefijo<br/>y claveNorm < hasta?}
    CMP -->|Sí| ADD["append Registro"]
    CMP -->|No, clave >= hasta| RET(["return"])
    ADD --> LOOP
```

## limiteSuperiorPrefijo

```mermaid
flowchart LR
    P1["prefijo = com<br/>→ comee"] 
    P2["prefijo = o<br/>→ p"]
    P3["prefijo con 0xFF<br/>→ prefijo + U+10FFFF"]
```

## Ejemplo: prefijo "Can"

```mermaid
flowchart LR
    R["Rango semiabierto"] --> A["desde = can"]
    R --> B["hasta = cao"]
    A --> M["Match: Canción Simulada..."]
    A --> M2["Match: Can't Help..."]
    B --> X["Stop antes de claves >= cao"]
```

## Flujo API + Modo Negocio

```mermaid
sequenceDiagram
    participant U as Usuario
    participant B as BusinessMode.vue
    participant A as bplustreeApi
    participant H as manejarBuscarPrefijo
    participant P as BuscarPorPrefijoNombre

    U->>B: Escribe "com" en búsqueda
    B->>B: debounce 280ms
    B->>A: GET /api/buscar-prefijo?prefijo=com
    A->>H: query prefijo
    H->>P: catalogo.BuscarPorPrefijoNombre
    P->>P: BuscarPorPrefijo(PorNombre, "com")
    P-->>H: []Registro
    H-->>B: JSON canciones
    B->>B: SongGridCard render
```
