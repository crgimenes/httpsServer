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

## Configurando

A configuração é propositalmente bem simples, apenas um arquivo JSON como no exemplo:

```json
[
    {
        "pattern": "example1.com/",
        "root": "assets.example1.com",
        "key": {
            "certFile": "/etc/letsencrypt/live/example1.com/fullchain.pem",
            "keyFile": "/etc/letsencrypt/live/example1.com/privkey.pem"
        }
    },
    {
        "pattern": "example2.com/",
        "root": "example2.com",
        "key": {
            "certFile": "/etc/letsencrypt/live/example2.com/fullchain.pem",
            "keyFile": "/etc/letsencrypt/live/example2.com/privkey.pem"
        }
    }
]
```

## Executando

Para executar basta chamar `./httpsServer` isso vai subir o servidor mas para deixar rodando em segundo plano faça da seguinte maneira:

```bash
nohup ./httpsServer&
```
