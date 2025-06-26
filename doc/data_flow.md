## Runtime Flow ##
```mermaid
graph LR
subgraph ETL Flow
    A1[pharmacies.json] --> B1(ETL Parser)
    A2[user.json] --> B1
    B1 --> C1((Go Structs))
    C1 --> D1[Database Migration Code]
    D1 --> E1[(PostgreSQL)]
end

subgraph API Runtime
    F[Client HTTP Request] --> G[Middleware Layer]
    G --> H[Gin Controller]
    H --> I[GORM Query]
    I --> E1
    E1 --> H --> J[JSON Response]
end