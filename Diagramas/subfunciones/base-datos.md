# Subfunciones: Base de Datos SQLite

Archivo: `database/connection.go`

## Conexión y tabla

```mermaid
flowchart LR
    MAIN["main.go"] --> CON["Conectar data/dataset.db"]
    CON --> SQL["modernc.org/sqlite"]
    SQL --> T["tabla tracks<br/>21 columnas Spotify"]
```

## CargarDatosEnCatalogo

```mermaid
flowchart TD
    START(["CargarDatosEnCatalogo(cat, limite)"]) --> Q["SELECT * FROM tracks LIMIT limite"]
    Q --> LOOP["Para cada fila"]
    LOOP --> LR["leerRegistro → Registro"]
    LR --> IR["cat.InsertarRegistro(reg)"]
    IR --> LOOP
```

## Operaciones CRUD en BD

```mermaid
flowchart TB
    subgraph INSERT
        SI["SiguienteIndice MAX+1"] --> IB["InsertarEnBD reg"]
        IB --> IR2["InsertarRegistro cat"]
    end
    subgraph DELETE
        EB["EliminarDeBD id"] --> ER["EliminarRegistro cat"]
    end
    subgraph READ
        LT["ListarTracks limite"]
        CT["ContarTracks"]
    end
```

## leerRegistro — mapeo SQL → Go

```mermaid
flowchart LR
    ROW["fila SQL"] --> R["Registro"]
    R --> F1["Indice ← id"]
    R --> F2["TrackName ← track_name"]
    R --> F3["Popularity ← popularity"]
    R --> F4["Tempo ← tempo"]
    R --> F5["Danceability ← danceability"]
    R --> F6["... 16 campos más"]
```

## Secuencia insertar persistente

```mermaid
sequenceDiagram
    participant A as manejarInsertar
    participant D as database.DB
    participant C as CatalogoIndices
    participant S as SQLite

    A->>D: SiguienteIndice()
    D->>S: SELECT MAX(id)+1
    A->>D: InsertarEnBD(reg)
    D->>S: INSERT INTO tracks
    A->>C: InsertarRegistro(reg)
    C->>C: 5 × Tree.Insertar
```
