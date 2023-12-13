package enums

type ContextKey string

const (
	ConfigCtxKey   ContextKey = "config.ctx.key"
	GormCtxKey     ContextKey = "gorm.ctx.key"
	LoggerCtxKey   ContextKey = "logger.ctx.key"
	SmtpCtxKey     ContextKey = "smtp.ctx.key"
	TemplateCtxKey ContextKey = "template.ctx.key"
)
