# Subfunciones: Navegación Interna

Archivo: `bplustree/interno.go` — usado por buscar, insertar y eliminar.

## indiceHijo — elegir puntero en nodo índice

```mermaid
flowchart TD
    START(["indiceHijo(nodo, clave)"]) --> UQ{clavesUnicas?}
    UQ -->|Sí primario| A["i++ mientras clave >= separadores i<br/>criterio >= separador"]
    UQ -->|No secundario| B["i++ mientras clave > separadores i<br/>criterio > separador"]
    A --> RET(["return i"])
    B --> RET
```

## irAHoja — bajar hasta hoja candidata

```mermaid
flowchart TD
    START(["irAHoja(clave)"]) --> N["n = raíz"]
    N --> L{es hoja?}
    L -->|No| I["n = n.hijos indiceHijo n, clave"]
    I --> L
    L -->|Sí| RET(["return n"])
```

## bajarGuardandoRuta — para insertar con split

```mermaid
flowchart TD
    START(["bajarGuardandoRuta(clave)"]) --> INIT["padres=[], posiciones=[], n=raíz"]
    INIT --> L{es hoja?}
    L -->|No| I["i = indiceHijo(n, clave)"]
    I --> SAVE["append padres, posiciones"]
    SAVE --> DESC["n = n.hijos[i]"]
    DESC --> L
    L -->|Sí| RET(["return n, padres, posiciones"])
```

## posicionEntrada — búsqueda binaria lineal en hoja

```mermaid
flowchart TD
    START(["posicionEntrada(entradas, clave)"]) --> LOOP["Para i, e en entradas"]
    LOOP --> CMP{e.clave >= clave?}
    CMP -->|Sí| RET(["return i"])
    CMP -->|No| LOOP
    LOOP -->|fin| END(["return len entradas"])
```

## rutaAHoja — reconstruir padres tras eliminar

```mermaid
flowchart TD
    START(["rutaAHoja(target)"]) --> DFS["DFS desde raíz"]
    DFS --> FOUND{encontró target?}
    FOUND -->|Sí| REV["Invertir padres y posiciones"]
    REV --> RET(["return padres, posiciones"])
    FOUND -->|No| NIL(["return nil, nil"])
```

## Helpers de comparación

```mermaid
flowchart LR
    compare["compare(a,b)"] --> cmp["cmp.Compare"]
    menor["menor"] --> compare
    mayor["mayor"] --> compare
    noMenor["noMenor"] --> compare
    posicion["posicionEntrada"] --> menor
```
