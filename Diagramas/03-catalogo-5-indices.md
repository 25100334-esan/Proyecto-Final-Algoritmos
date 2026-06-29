# Catálogo de 5 Índices (CatalogoIndices)

Modelo según Comer: índice primario único + secundarios con duplicados y sequence set.

```mermaid
flowchart TB
    REG["Registro<br/>Indice · TrackName · Popularity<br/>Tempo · Danceability · ..."]

    subgraph CAT["CatalogoIndices"]
        PI["PorIndice<br/>Tree int, Registro<br/>clavesUnicas = true"]
        PN["PorNombre<br/>Tree string, Registro<br/>clavesUnicas = false"]
        PP["PorPopularidad<br/>Tree int, Registro"]
        PT["PorTempo<br/>Tree float64, Registro"]
        PD["PorDanceability<br/>Tree float64, Registro"]
    end

    REG -->|InsertarRegistro| PI
    REG -->|InsertarRegistro| PN
    REG -->|InsertarRegistro| PP
    REG -->|InsertarRegistro| PT
    REG -->|InsertarRegistro| PD

    PI -->|EliminarRegistro| PN
    PI -->|EliminarRegistro| PP
    PI -->|EliminarRegistro| PT
    PI -->|EliminarRegistro| PD
```

## InsertarRegistro — cascada hacia 5 árboles

```mermaid
flowchart LR
    A["InsertarRegistro(reg)"] --> B["PorIndice.Insertar(reg.Indice, reg)"]
    A --> C["PorNombre.Insertar(reg.TrackName, reg)"]
    A --> D["PorPopularidad.Insertar(reg.Popularity, reg)"]
    A --> E["PorTempo.Insertar(reg.Tempo, reg)"]
    A --> F["PorDanceability.Insertar(reg.Danceability, reg)"]
```

## EliminarRegistro — cascada con EliminarExacto

```mermaid
flowchart LR
    A["EliminarRegistro(reg)"] --> B["PorIndice.Eliminar(reg.Indice)"]
    A --> C["EliminarExacto PorNombre<br/>TrackName + Indice"]
    A --> D["EliminarExacto PorPopularidad"]
    A --> E["EliminarExacto PorTempo"]
    A --> F["EliminarExacto PorDanceability"]
```

## Estructura de un árbol B+ (conceptual)

```mermaid
flowchart TB
    ROOT["Nodo ÍNDICE<br/>separadores + hijos"]
    L1["Hoja ← sequence set → Hoja → Hoja"]
    ROOT --> L1

    subgraph HOJA["Nodo HOJA"]
        E1["entrada: clave → valor Registro"]
        E2["entrada: clave → valor Registro"]
        SS["siguienteHoja puntero horizontal"]
    end
```
