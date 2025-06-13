# Go gin Project

Um projeto de teste em Go que demonstra a implementação de um engine web customizado com funcionalidades de cron jobs integradas.

## 📋 Descrição

Este é um projeto experimental desenvolvido para testar e demonstrar:

- **Engine Web Customizado**: Wrapper personalizado sobre o Gin framework
- **Sistema de Cron Jobs**: Gerenciamento de tarefas agendadas com robfig/cron
- **Logger Estruturado**: Sistema de logging com Zap (console + arquivo)
- **Context Personalizado**: Contextos customizados para handlers e cron jobs

## 🚀 Funcionalidades

### Engine Web

- Wrapper sobre Gin com API familiar
- Context customizado para handlers
- Métodos HTTP padrão (GET, POST, PUT, DELETE, etc.)
- Suporte a middlewares

### Cron Manager

- Agendamento de tarefas com sintaxe cron
- Suporte a segundos (`*/5 * * * * *`)
- Gerenciamento de jobs por nome
- Prevenção de jobs duplicados
- Context específico para cron jobs

### Sistema de Logging

- Logs coloridos no console
- Logs estruturados em arquivo (JSON)
- Formato personalizado: `LEVEL CALLER TIME MESSAGE`
- Separação de responsabilidades

## 📁 Estrutura do Projeto

```
.
├── main.go                 # Ponto de entrada
├── app.log                # Arquivo de logs
├── internal/app/
│   ├── engine.go          # Engine principal
│   ├── context.go         # Contextos customizados
│   ├── cron.go           # Gerenciador de cron
│   └── logger.go         # Configuração do logger
└── README.md             # Este arquivo
```

## 🛠️ Tecnologias Utilizadas

- **[Gin](https://github.com/gin-gonic/gin)**: Framework web
- **[Robfig Cron](https://github.com/robfig/cron)**: Scheduler de cron jobs
- **[Zap](https://github.com/uber-go/zap)**: Logger estruturado

## 📖 Exemplo de Uso

```go
package main

import (
    "github.com/gin-gonic/gin"
    "github.com/matheusvcouto/go-test/internal/app"
)

func main() {
    r := gin.Default()
    engine := app.New(r)

    // Rota web
    engine.GET("/", func(c *app.Context) {
        c.JSON(200, gin.H{"message": "Hello World!"})
    })

    // Cron job
    engine.Cron().Add("*/5 * * * * *", "test-job", func(ctx *app.CronContext) {
        ctx.GetLogger().Info("Cron job executado!")
    })

    engine.Run(":8080")
}
```

## 🎯 Objetivo

Este projeto serve como base para experimentação e teste de conceitos em Go, incluindo:

- Padrões de design para engines web
- Integração de sistemas de cron
- Configuração avançada de logging
- Arquitetura modular e extensível

## ⚠️ Aviso

Este é um **projeto de teste** e não deve ser usado em produção sem as devidas adaptações e testes de segurança.

---

_Projeto desenvolvido para fins educacionais e experimentação com Go._
