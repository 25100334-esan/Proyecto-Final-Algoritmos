# Interacción: Flujo Eliminar en Cascada

Eliminar por ID en índice primario → cascada en 4 secundarios con `EliminarExacto`.

```mermaid
sequenceDiagram
    participant U as Usuario
    participant A as AcademicMode.vue
    participant API as bplustreeApi
    participant H as manejarEliminar
    participant D as database.DB
    participant C as CatalogoIndices

    U->>A: Eliminar ID=114035 (vista TrackName)
    A->>API: GET /api/buscar?indice=114035
    API-->>A: Registro {TrackName, ...}
    A->>A: calcularSecuenciaEliminacion<br/>clave + idCancion
    A->>A: deleteAnimation + leafLocate
    A->>API: DELETE /api/eliminar?indice=114035
    API->>H: query indice
    H->>C: PorIndice.Buscar(114035)
    H->>D: EliminarDeBD(114035)
    H->>C: EliminarRegistro(reg)
    Note over C: PorIndice.Eliminar<br/>EliminarExacto × 4 secundarios
    H-->>API: {ok, cascada: true}
    A->>A: refrescarEstructura + estadísticas
```

## EliminarRegistro detalle

```mermaid
flowchart TD
    ER["EliminarRegistro(reg)"] --> E0["PorIndice.Eliminar reg.Indice"]
    ER --> E1["EliminarExacto PorNombre<br/>reg.TrackName, reg.Indice"]
    ER --> E2["EliminarExacto PorPopularidad"]
    ER --> E3["EliminarExacto PorTempo"]
    ER --> E4["EliminarExacto PorDanceability"]

    E1 --> ES["EliminarSi + coincideIndice"]
    ES --> UF["Underflow si aplica"]
```

## Animación en vista secundaria

```mermaid
flowchart LR
    ID["Usuario ingresa ID"] --> BUS["api.buscar → Registro"]
    BUS --> CLV["claveVista.claveDeCancionPorEstructura"]
    CLV --> LOC["leafLocate.localizarEntradaEnHojas<br/>por idsRegistro"]
    LOC --> DEL["deleteAnimation.calcularSecuenciaEliminacion"]
    DEL --> API["api.eliminar"]
```

## Frontend utils involucrados

```mermaid
flowchart TB
    AM["AcademicMode.vue"] --> CV["claveVista.js"]
    AM --> LL["leafLocate.js"]
    AM --> DA["deleteAnimation.js"]
    AM --> TL["treeLayout.js"]
    DA --> DO["DeleteOverlay.vue"]
```
