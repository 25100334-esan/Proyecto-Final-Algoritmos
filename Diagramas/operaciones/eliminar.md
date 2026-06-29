# Operación: Eliminar

**API:** `Tree.Eliminar(clave)` / `EliminarExacto(t, clave, id)` — O(log d n)

```mermaid
flowchart TD
    START(["EliminarSi(clave, predicado)"]) --> H["hoja = irAHoja(clave)"]
    H --> LOOP{Recorrer hojas<br/>sequence set}
    LOOP --> POS["pos = posicionEntrada"]
    POS --> SCAN["Para i con clave == buscada"]
    SCAN --> PRED{predicado valor?}
    PRED -->|No| SCAN
    PRED -->|Sí| DEL["Quitar entrada[i]"]
    DEL --> MIN{len >= minPorNodo<br/>o es raíz?}
    MIN -->|Sí| OK(["return true"])
    MIN -->|No| UF["corregirUnderflowHoja"]
    UF --> OK
    SCAN -->|Sin match| NEXT{última clave == buscada?}
    NEXT -->|Sí| LOOP
    NEXT -->|No| FAIL(["return false"])
```

## EliminarExacto en índices secundarios

```mermaid
flowchart LR
    EE["EliminarExacto(t, clave, indice)"] --> ES["EliminarSi(clave, coincideIndice)"]
    ES --> MATCH["Coincide clave Y reg.Indice == id"]
    MATCH --> ONE["Elimina solo esa entrada<br/>no todas con misma clave"]
```

## Underflow en hoja — decisión

```mermaid
flowchart TD
    CU["corregirUnderflowHoja"] --> PD{Hermano der<br/>tiene extras?}
    PD -->|Sí| PRD["prestarDerecha"]
    PD -->|No| PI{Hermano izq<br/>tiene extras?}
    PI -->|Sí| PRI["prestarIzquierda"]
    PI -->|No| FUS{¿Cuál fusionar?}
    FUS --> FD["fusionarConDerecha"]
    FUS --> FI["fusionarConIzquierda"]
    PRD & PRI & FD & FI --> AR["ajustarRaiz"]
    AR --> CI{padre underflow?}
    CI -->|Sí| CUI["corregirUnderflowInterno"]
```

## Cascada EliminarRegistro

```mermaid
flowchart TB
    API["DELETE /api/eliminar?indice=N"] --> B["PorIndice.Buscar(N)"]
    B --> DB["db.EliminarDeBD(N)"]
    DB --> ER["EliminarRegistro(reg)"]
    ER --> E1["PorIndice.Eliminar"]
    ER --> E2["EliminarExacto PorNombre"]
    ER --> E3["EliminarExacto PorPopularidad"]
    ER --> E4["EliminarExacto PorTempo"]
    ER --> E5["EliminarExacto PorDanceability"]
```
