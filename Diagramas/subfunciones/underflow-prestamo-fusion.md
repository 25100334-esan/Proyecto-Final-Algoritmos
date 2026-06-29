# Subfunciones: Underflow — Préstamo y Fusión

Archivo: `bplustree/eliminar.go`

## corregirUnderflowHoja

```mermaid
flowchart TD
    CU["corregirUnderflowHoja"] --> CHK{len >= min o sin padre?}
    CHK -->|Sí| END(["return"])
    CHK -->|No| PD["Hermano derecho existe?"]
    PD -->|Sí y tiene extra| PRD["prestarDerecha"]
    PD -->|No| PI["Hermano izquierdo tiene extra?"]
    PI -->|Sí| PRI["prestarIzquierda"]
    PI -->|No| FUS["Fusionar con hermano"]
    PRD & PRI & FUS --> AR["ajustarRaiz"]
    AR --> CI{padre bajo mínimo?}
    CI -->|Sí| CUI["corregirUnderflowInterno recursivo"]
```

## prestarDerecha

```mermaid
flowchart LR
    A["Mover primera entrada<br/>hermano der → hoja"] --> B["Actualizar separador padre<br/>con nueva primera clave der"]
```

## prestarIzquierda

```mermaid
flowchart LR
    A["Mover última entrada<br/>hermano izq → inicio hoja"] --> B["Actualizar separador padre<br/>con primera clave hoja"]
```

## fusionarConDerecha

```mermaid
flowchart TD
    F["fusionarConDerecha"] --> M1["hoja += entradas hermano der"]
    M1 --> M2["hoja.siguienteHoja = hermano.siguienteHoja"]
    M2 --> M3["Quitar separador e hijo del padre"]
```

## corregirUnderflowInterno

```mermaid
flowchart TD
    CUI["corregirUnderflowInterno"] --> PD["prestarSeparadorDerecha?"]
    PD -->|No| PI["prestarSeparadorIzquierda?"]
    PI -->|No| FI["fusionarInternoConDerecha/Izquierda"]
    FI --> AR["ajustarRaiz"]
    AR --> REC["Recursión hacia arriba si padre underflow"]
```

## ajustarRaiz — colapsar raíz vacía

```mermaid
flowchart TD
    AR["ajustarRaiz(n)"] --> C1{n == raíz y es índice?}
    C1 -->|No| END(["return"])
    C1 -->|Sí| C2{0 separadores y 1 hijo?}
    C2 -->|Sí| NEW["raíz = único hijo"]
    C2 -->|No| END
```

## EliminarSi + sequence set

```mermaid
flowchart LR
    ES["EliminarSi"] --> H1["irAHoja"]
    H1 --> SCAN["Escanear duplicados en hoja"]
    SCAN --> SS["Si última clave == buscada<br/>→ siguienteHoja"]
    SS --> SCAN
```
