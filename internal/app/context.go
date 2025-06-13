package app

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type Context struct {
	*gin.Context
	app *Engine
}

// Cron retorna o gerenciador de cron do App
func (c *Context) Cron() *CronManager {
	return c.app.cronManager
}

// Logger retorna seu logger
func (c *Context) Logger() *zap.Logger {
	return c.app.GetLogger()
}

// NotFound é um exemplo de helper pra abortar com 404
func (c *Context) NotFound() {
	c.JSON(404, gin.H{"error": "not found"})
	c.Abort()
}

// CronContext é o contexto específico para cron jobs
type CronContext struct {
	engine *Engine
}

// Logger retorna o logger para cron jobs
func (c *CronContext) Logger() *zap.Logger {
	return c.engine.GetLogger()
}

// [adicione aqui todos os seus helpers…]
