package app

import (
	"github.com/robfig/cron/v3"
	"go.uber.org/zap"
)

// CronManager gerencia todas as funcionalidades relacionadas ao cron
type CronManager struct {
	cron   *cron.Cron
	logger *zap.Logger
	jobs   map[string]cron.EntryID // mapeamento de nome para ID do job
	engine *Engine                 // referência ao engine para usar no context
}

// CronJobFunc é a função que será executada pelo cron job
type CronJobFunc func(*CronContext)

// NewCronManager cria um novo gerenciador de cron
func NewCronManager(logger *zap.Logger) *CronManager {
	return &CronManager{
		cron:   cron.New(cron.WithSeconds()), // permite usar segundos
		logger: logger,
		jobs:   make(map[string]cron.EntryID),
		engine: nil, // será definido depois
	}
}

// SetEngine define a referência do engine
func (cm *CronManager) SetEngine(engine *Engine) {
	cm.engine = engine
}

// Add adiciona um novo job de cron sem retornar erro
// spec: especificação do cron (ex: "*/10 * * * * *" para cada 10 segundos)
// name: nome do job (opcional, usado para remover depois)
// fn: função a ser executada
func (cm *CronManager) Add(spec, name string, fn CronJobFunc) {
	// Verifica se já existe um job com o mesmo nome
	if name != "" {
		if _, exists := cm.jobs[name]; exists {
			// panic("PANIC: Cron job with name '" + name + "' already exists!")
			cm.logger.Warn("Cron job with same name already exists - skipping addition",
				zap.String("name", name),
				zap.String("spec", spec),
			)
			return // Não adiciona o job e retorna
		}
	}

	entryID, err := cm.cron.AddFunc(spec, func() {
		// cm.logger.Info("Executing cron job",
		// 	zap.String("name", name),
		// 	zap.String("spec", spec),
		// 	zap.Time("time", time.Now()),
		// )

		// Cria o CronContext específico para cron jobs
		cronCtx := &CronContext{
			engine: cm.engine,
		}

		// Executa a função do cron job
		fn(cronCtx)
	})

	if err != nil {
		cm.logger.Error("Failed to add cron job - job will not run",
			zap.String("name", name),
			zap.String("spec", spec),
			zap.Error(err),
		)
		return // Não para o servidor, apenas não adiciona o job
	}

	// Armazena o ID do job se um nome foi fornecido
	if name != "" {
		cm.jobs[name] = entryID
	}

	cm.logger.Info("Cron job added successfully",
		zap.String("name", name),
		zap.String("spec", spec),
		zap.Int("entryID", int(entryID)),
	)
}

// Remove remove um job de cron pelo nome
func (cm *CronManager) Remove(name string) {
	if entryID, exists := cm.jobs[name]; exists {
		cm.cron.Remove(entryID)
		delete(cm.jobs, name)
		cm.logger.Info("Cron job removed", zap.String("name", name))
	} else {
		cm.logger.Warn("Cron job not found", zap.String("name", name))
	}
}

// Start inicia o scheduler de cron
func (cm *CronManager) Start() {
	cm.cron.Start()
	cm.logger.Info("Cron scheduler started")
}

// Stop para o scheduler de cron
func (cm *CronManager) Stop() {
	cm.cron.Stop()
	cm.logger.Info("Cron scheduler stopped")
}

// GetCron retorna a instância do cron para uso avançado
func (cm *CronManager) GetCron() *cron.Cron {
	return cm.cron
}

// ListJobs lista todos os jobs ativos
func (cm *CronManager) ListJobs() []string {
	var jobNames []string
	for name := range cm.jobs {
		jobNames = append(jobNames, name)
	}
	return jobNames
}

// GetJobsCount retorna o número de jobs ativos
func (cm *CronManager) GetJobsCount() int {
	return len(cm.jobs)
}
