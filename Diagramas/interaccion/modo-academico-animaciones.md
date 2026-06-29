# Interacción: Modo Académico y Animaciones

Vista: `AcademicMode.vue` — canvas B+ con operaciones paso a paso.

## Componentes principales

```mermaid
flowchart TB
    AM["AcademicMode.vue"] --> CP["ControlPanel.vue"]
    AM --> TC["TreeCanvas.vue"]
    AM --> TN["TreeNode.vue"]
    AM --> OL["OperationLog.vue"]
    AM --> SB["StepBanner.vue"]
    AM --> TSP["TreeStatsPanel.vue"]
    AM --> SO["SplitOverlay.vue"]
    AM --> DO["DeleteOverlay.vue"]

    AM --> UTS["useTreeState.js"]
    AM --> API["bplustreeApi.js"]
```

## Mapa operación → animación

```mermaid
flowchart LR
    subgraph Operaciones
        O1["onBuscar"]
        O2["onInsertar"]
        O3["onEliminar"]
        O4["onRango"]
        O5["onInicializar"]
    end

    subgraph Animaciones
        A1["searchAnimation.js"]
        A2["insertAnimation.js"]
        A3["deleteAnimation.js"]
        A4["rangeAnimation.js"]
        A5["initAnimation.js"]
    end

    subgraph Layout
        L["treeLayout.js"]
        N["nodoClaves.js"]
        F["formatKeys.js"]
    end

    O1 --> A1 --> L
    O2 --> A2 --> L
    O3 --> A3 --> L
    O4 --> A4 --> L
    O5 --> A5 --> L
    L --> TC["TreeCanvas"]
```

## Ciclo de una operación animada

```mermaid
stateDiagram-v2
    [*] --> Idle
    Idle --> CalcularPasos: Usuario elige operación
    CalcularPasos --> Animando: animando=true
    Animando --> PasoN: aplicarPasoGenerico
    PasoN --> Animando: siguiente paso
    PasoN --> LlamarAPI: paso aplicar
    LlamarAPI --> Refrescar: api.* + estructura
    Refrescar --> Idle: animando=false
```

## Límite de visualización

```mermaid
flowchart TD
    C["conteo canciones"] --> L{≤ 500?}
    L -->|Sí| CAN["TreeCanvas visible<br/>LIMITE_VISUALIZACION"]
    L -->|No| STATS["Solo TreeStatsPanel<br/>sin canvas"]
```

## Vistas de árbol intercambiables

```mermaid
flowchart LR
    SEL["Selector vista"] --> I["indice / ID"]
    SEL --> N["TrackName"]
    SEL --> P["Popularidad"]
    SEL --> T["Tempo"]
    SEL --> D["Danceability"]
    I & N & P & T & D --> API["GET /api/estructura?tipo="]
    API --> RENDER["TreeCanvas renderiza árbol activo"]
```
