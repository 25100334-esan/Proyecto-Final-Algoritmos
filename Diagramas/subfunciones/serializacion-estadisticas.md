# Subfunciones: Serialización y Estadísticas

Archivos: `bplustree/serializar.go`, `bplustree/estadisticas.go`

## ExportarEstructura — JSON para canvas

```mermaid
flowchart TD
    EE["ExportarEstructura / ExportarPorTipo"] --> BFS["Recorrido BFS del árbol"]
    BFS --> NOD["Por cada nodo crear NodoExport"]
    NOD --> TRUNC["truncarClave runas UTF-8"]
    NOD --> IDS["idsRegistro desde valor.Indice"]
    BFS --> CAD["cadenaHojas via sequence set"]
    CAD --> JSON["EstructuraExport JSON"]
```

## Campos EstructuraExport

```mermaid
classDiagram
    class EstructuraExport {
        +tipoIndice string
        +tipoClave string
        +clavesUnicas bool
        +orden int
        +raiz int
        +cadenaHojas int[]
        +nodos NodoExport[]
    }
    class NodoExport {
        +id int
        +esHoja bool
        +claves any[]
        +separadores any[]
        +hijos int[]
        +idsRegistro int[]
    }
    EstructuraExport --> NodoExport
```

## EstadisticasArbol

```mermaid
flowchart LR
    ES["EstadisticasPorTipo"] --> C["Contar nodos hoja e índice"]
    C --> H["Calcular altura"]
    C --> K["Total claves"]
    C --> O["Ocupación promedio"]
```

## Frontend consume estructura

```mermaid
flowchart LR
    API["GET /api/estructura"] --> STATE["useTreeState.arbolOperaciones"]
    STATE --> LAYOUT["treeLayout.calcularLayout"]
    LAYOUT --> CANVAS["TreeCanvas.vue"]
    STATE --> ANIM["*Animation.js"]
```
