# Subfunciones: Split y Overflow

Archivo: `bplustree/insertar.go`

## Árbol de llamadas al insertar con overflow

```mermaid
flowchart TD
    INS["Tree.Insertar"] --> BG["bajarGuardandoRuta"]
    INS --> ADD["agregarEntradaOrdenada en hoja"]
    ADD --> OF{overflow hoja?}
    OF -->|Sí| PH["partirHoja"]
    PH --> SDH["subirDivisionHoja"]
    SDH --> INI["insertarEnNodoInterno"]
    SDH --> OFI{overflow interno?}
    OFI -->|Sí| PNI["partirNodoInterno"]
    PNI --> SDI["subirDivisionInterna"]
    SDI --> NR["Nueva raíz si padres vacíos"]
```

## partirHoja — detalle

```mermaid
flowchart TD
    PH["partirHoja(hoja)"] --> M["medio = len+1 / 2"]
    M --> DER["Crear hoja derecha con entradas medio.."]
    M --> IZQ["Hoja izq conserva 0..medio-1"]
    DER --> SS["Enlazar sequence set:<br/>izq.siguienteHoja = der<br/>der.siguienteHoja = viejoSiguiente"]
    DER --> PROM["return claveCopiada = der.entradas0.clave"]
```

## partirNodoInterno — detalle

```mermaid
flowchart TD
    PNI["partirNodoInterno(n)"] --> M["medio = len separadores / 2"]
    M --> UP["claveSubida = separadores medio"]
    M --> DER["Derecha: sep medio+1.. y hijos medio+1.."]
    M --> IZQ["Izq: sep 0..medio-1 y hijos 0..medio"]
    UP --> RET(["return claveSubida, derecha"])
```

## Caso: raíz se convierte en índice

```mermaid
flowchart LR
    subgraph ANTES
        R1["Raíz HOJA llena"]
    end
    subgraph DESPUES
        R2["Raíz ÍNDICE<br/>sep = k_medio"]
        L["Hoja izq"]
        D["Hoja der"]
        R2 --> L
        R2 --> D
    end
    ANTES --> DESPUES
```

## Animación frontend (splitAnimation.js)

```mermaid
flowchart LR
    INS["insertAnimation.js"] --> DET["detectarSplits treeLayout"]
    DET --> ANA["analizarSplit"]
    ANA --> FASES["FASES_SPLIT overlay"]
    FASES --> CANVAS["SplitOverlay.vue"]
```
