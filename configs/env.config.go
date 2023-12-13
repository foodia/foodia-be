package configs

import "time"

type EnvConfig struct {
	AppName               string        `koanf:"APP_NAME"`
	AppEnv                string        `koanf:"APP_ENV"`
	AppPort               uint32        `koanf:"APP_PORT"`
	AppVersion            float32       `koanf:"APP_VERSION"`
	DBHost                string        `koanf:"DB_HOST"`
	DBPort                string        `koanf:"DB_PORT"`
	DBUser                string        `koanf:"DB_USER"`
	DBPassword            string        `koanf:"DB_PASSWORD"`
	DBName                string        `koanf:"DB_NAME"`
	JWTSecret             string        `koanf:"JWT_SECRET"`
	JWTExpirationDuration time.Duration `koanf:"JWT_EXPIRATION_DURATION"`
	LogFile               string        `koanf:"LOGFILE"`
	SmtpHost              string        `koanf:"SMTP_HOST"`
	SmtpPort              int           `koanf:"SMTP_PORT"`
	SmtpUser              string        `koanf:"SMTP_USER"`
	SmtpPass              string        `koanf:"SMTP_PASS"`
	SmtpSender            string        `koanf:"SMTP_SENDER"`
}
