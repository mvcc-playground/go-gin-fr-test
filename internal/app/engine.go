package app

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// Engine é o seu "engine" principal
type Engine struct {
	*gin.Engine              // embede tudo do gin.Engine
	cronManager *CronManager // gerenciador de cron personalizado
	logger      *zap.Logger  // seu logger (pode ser outro)
}

// New cria um Engine já configurado
func New(g *gin.Engine) *Engine {
	logger, _ := zap.NewProduction()
	cronManager := NewCronManager(logger)

	app := &Engine{
		Engine:      g,
		cronManager: cronManager,
		logger:      logger,
	}

	// Define a referência do engine no cronManager
	cronManager.SetEngine(app)

	// Inicia o cron automaticamente
	cronManager.Start()

	return app
}

type HandlerFunc func(*Context)

// Cron retorna o gerenciador de cron para uso direto
func (e *Engine) Cron() *CronManager {
	return e.cronManager
}

// GetLogger retorna a instância do logger
func (e *Engine) GetLogger() *zap.Logger {
	return e.logger
}

// StopCron para o scheduler de cron (chamado automaticamente no shutdown)
func (e *Engine) StopCron() {
	if e.cronManager != nil {
		e.cronManager.Stop()
	}
}

// GetCronManager retorna o gerenciador de cron
func (e *Engine) GetCronManager() *CronManager {
	return e.cronManager
}

// convertHandlers converte handlers customizados para gin.HandlerFunc
// Otimizado para melhor performance - evita alocações desnecessárias
func (e *Engine) convertHandlers(handlers []HandlerFunc) []gin.HandlerFunc {
	if len(handlers) == 0 {
		return nil
	}

	ginHandlers := make([]gin.HandlerFunc, len(handlers))
	for i, h := range handlers {
		handler := h // captura o handler para evitar closure issues
		ginHandlers[i] = func(c *gin.Context) {
			// Cria o Context customizado diretamente sem usar c.Set/c.Get
			// Isso evita alocações no map interno do gin.Context
			appCtx := &Context{Context: c, app: e}
			handler(appCtx)
		}
	}
	return ginHandlers
}

// AdaptGin converte gin.HandlerFunc para HandlerFunc customizado
// Permite usar middlewares nativos do Gin com o Context customizado
func (e *Engine) AdaptGin(ginHandler gin.HandlerFunc) HandlerFunc {
	return func(c *Context) {
		ginHandler(c.Context)
	}
}

// AdaptGinMany converte múltiplos gin.HandlerFunc para HandlerFunc customizados
func (e *Engine) AdaptGinMany(ginHandlers ...gin.HandlerFunc) []HandlerFunc {
	adapted := make([]HandlerFunc, len(ginHandlers))
	for i, h := range ginHandlers {
		handler := h // captura para evitar closure issues
		adapted[i] = func(c *Context) {
			handler(c.Context)
		}
	}
	return adapted
}

// Use adiciona middleware ao engine
func (e *Engine) Use(handlers ...HandlerFunc) {
	e.Engine.Use(e.convertHandlers(handlers)...)
}

// GET é um atalho para router.Handle("GET", path, handlers)
func (e *Engine) GET(relativePath string, handlers ...HandlerFunc) {
	e.Engine.GET(relativePath, e.convertHandlers(handlers)...)
}

// POST é um atalho para router.Handle("POST", path, handlers)
func (e *Engine) POST(relativePath string, handlers ...HandlerFunc) {
	e.Engine.POST(relativePath, e.convertHandlers(handlers)...)
}

// DELETE é um atalho para router.Handle("DELETE", path, handlers)
func (e *Engine) DELETE(relativePath string, handlers ...HandlerFunc) {
	e.Engine.DELETE(relativePath, e.convertHandlers(handlers)...)
}

// PATCH é um atalho para router.Handle("PATCH", path, handlers)
func (e *Engine) PATCH(relativePath string, handlers ...HandlerFunc) {
	e.Engine.PATCH(relativePath, e.convertHandlers(handlers)...)
}

// PUT é um atalho para router.Handle("PUT", path, handlers)
func (e *Engine) PUT(relativePath string, handlers ...HandlerFunc) {
	e.Engine.PUT(relativePath, e.convertHandlers(handlers)...)
}

// OPTIONS é um atalho para router.Handle("OPTIONS", path, handlers)
func (e *Engine) OPTIONS(relativePath string, handlers ...HandlerFunc) {
	e.Engine.OPTIONS(relativePath, e.convertHandlers(handlers)...)
}

// HEAD é um atalho para router.Handle("HEAD", path, handlers)
func (e *Engine) HEAD(relativePath string, handlers ...HandlerFunc) {
	e.Engine.HEAD(relativePath, e.convertHandlers(handlers)...)
}

// Handle registra um novo handler e middleware com o path e método especificado
func (e *Engine) Handle(httpMethod, relativePath string, handlers ...HandlerFunc) {
	e.Engine.Handle(httpMethod, relativePath, e.convertHandlers(handlers)...)
}

// Any registra uma rota que corresponde a todos os métodos HTTP
func (e *Engine) Any(relativePath string, handlers ...HandlerFunc) {
	e.Engine.Any(relativePath, e.convertHandlers(handlers)...)
}

// Match registra uma rota que corresponde aos métodos especificados
func (e *Engine) Match(methods []string, relativePath string, handlers ...HandlerFunc) {
	e.Engine.Match(methods, relativePath, e.convertHandlers(handlers)...)
}

// Group cria um novo grupo de rotas
func (e *Engine) Group(relativePath string, handlers ...HandlerFunc) *gin.RouterGroup {
	return e.Engine.Group(relativePath, e.convertHandlers(handlers)...)
}

// StaticFile registra uma única rota para servir um único arquivo
func (e *Engine) StaticFile(relativePath, filepath string) {
	e.Engine.StaticFile(relativePath, filepath)
}

// StaticFileFS funciona como StaticFile mas com um http.FileSystem customizado
func (e *Engine) StaticFileFS(relativePath, filepath string, fs http.FileSystem) {
	e.Engine.StaticFileFS(relativePath, filepath, fs)
}

// Static serve arquivos do diretório especificado
func (e *Engine) Static(relativePath, root string) {
	e.Engine.Static(relativePath, root)
}

// StaticFS funciona como Static mas com um http.FileSystem customizado
func (e *Engine) StaticFS(relativePath string, fs http.FileSystem) {
	e.Engine.StaticFS(relativePath, fs)
}
