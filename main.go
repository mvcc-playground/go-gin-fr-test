package main

import (
	"fmt"
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/matheusvcouto/go-test/internal/app"
	"go.uber.org/zap"
)

func main() {
	// Cria o engine customizado - cron inicia automaticamente
	r := app.New(gin.New())

	// Middleware personalizado
	r.Use(func(c *app.Context) {
		c.Logger().Info("Request received", zap.String("ip", c.RemoteIP()))
		c.Next()
	})

	// Rotas GET
	r.GET("/", func(c *app.Context) {
		// c.Cron().Remove("frequent-task") // Exemplo de remoção de um job de cron
		c.JSON(200, gin.H{
			"message": "Hello, World!",
			"jobs":    c.Cron().GetJobsCount(),
		})
	})

	// Run the cron job every 5 seconds
	r.Cron().Add("*/5 * * * * *", "frequent-task", func(c *app.CronContext) {
		fmt.Println("Cron job executed at", time.Now().Format(time.RFC3339))
		c.Logger().Info("Cron job executed")
	})
	r.Cron().Add("*/10 * * * * *", "frequent-tasks", func(c *app.CronContext) {
		fmt.Println("Cron job executed at", time.Now().Format(time.RFC3339))
	})

	// Inicia o servidor
	log.Println("Servidor iniciando na porta 8080...")
	r.Run(":8080")
}
