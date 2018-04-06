# HTTPS Server

Um servidor HTTPS com capacidade para servir múltiplos domínios simultaneamente

## Como compilar

Para compilar basta o seguinte:

```bash
go build
```

Supondo que você esta em usando Windows ou macOS e quer compilar para Linux, faça assim:

```bash
GOOS=linux CGO_ENABLED=0 go build
```