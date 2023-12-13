package main

import (
	"context"
	"embed"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"

	"foodia-be/common"
	"foodia-be/configs"
	"foodia-be/enums"
	"foodia-be/routers"

	"github.com/goccy/go-json"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/filesystem"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"

	"github.com/knadh/koanf/parsers/dotenv"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
	"github.com/wneessen/go-mail"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

//go:embed templates/*
var templateFS embed.FS

func init() {
	createDirIfNotExists := func(path string) error {
		if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
			if err = os.MkdirAll(path, os.ModePerm); err != nil {
				return err
			}

			log.Println(fmt.Sprintf("%s directory created", path))
		}
		return nil
	}

	logsPath := "logs"
	if err := createDirIfNotExists(logsPath); err != nil {
		log.Fatalf("Error creating %s directory:%v", logsPath, err)
	}

	detonatorPath := "storage/detonator"
	if err := createDirIfNotExists(detonatorPath); err != nil {
		log.Fatalf("Error creating %s directory:%v", detonatorPath, err)
	}

	merchantPath := "storage/merchant"
	if err := createDirIfNotExists(merchantPath); err != nil {
		log.Fatalf("Error creating %s directory:%v", merchantPath, err)
	}

	campaignPath := "storage/campaign"
	if err := createDirIfNotExists(campaignPath); err != nil {
		log.Fatalf("Error creating %s directory:%v", campaignPath, err)
	}

	productPath := "storage/product"
	if err := createDirIfNotExists(productPath); err != nil {
		log.Fatalf("Error creating %s directory:%v", productPath, err)
	}
}

func main() {
	conf := koanf.New(".")
	if err := conf.Load(file.Provider(".env"), dotenv.Parser()); err != nil {
		log.Fatalf("error loading .env file: %v", err)
	}

	config := new(configs.EnvConfig)
	if err := conf.UnmarshalWithConf("", &config, koanf.UnmarshalConf{FlatPaths: false}); err != nil {
		log.Fatalf("failed to read .env file: %v", err)
	}

	db, err := gorm.Open(mysql.Open(config.GetDSN()), &gorm.Config{})
	if err != nil {
		log.Fatal(err.Error())
	}

	smtp, err := mail.NewClient(
		config.SmtpHost,
		mail.WithPort(config.SmtpPort),
		mail.WithSMTPAuth(mail.SMTPAuthPlain),
		mail.WithUsername(config.SmtpUser),
		mail.WithPassword(config.SmtpPass),
		mail.WithSSL(),
	)
	if err != nil {
		// log.Fatalf("failed connect to smtp server with error: %v", err)
	}

	logfile, err := common.NewLogger(config.LogFile)
	if err != nil {
		log.Fatalf("failed connect to logger service with error: %v", err)
	}

	ctx := context.WithValue(context.Background(), enums.GormCtxKey, db)
	ctx = context.WithValue(ctx, enums.ConfigCtxKey, config)
	ctx = context.WithValue(ctx, enums.LoggerCtxKey, logfile)
	ctx = context.WithValue(ctx, enums.TemplateCtxKey, templateFS)
	ctx = context.WithValue(ctx, enums.SmtpCtxKey, smtp)

	app := fiber.New(fiber.Config{
		ProxyHeader: fiber.HeaderXForwardedFor,
		AppName:     config.AppName,
		JSONEncoder: json.Marshal,
		JSONDecoder: json.Unmarshal,
	})

	app.Use(logger.New())
	app.Use(recover.New())

	app.Use(cors.New(cors.Config{
		AllowOrigins:     "*",
		AllowMethods:     "GET,POST,PUT,DELETE",
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization",
		ExposeHeaders:    "Content-Length",
		AllowCredentials: true,
	}))

	// File System
	app.Use("/storage", filesystem.New(filesystem.Config{
		Root:   http.Dir("./storage"),
		Browse: false,
	}))

	routers.UseRouter(ctx, app)

	if err = app.Listen(fmt.Sprintf(":%d", config.AppPort)); err != nil {
		log.Fatal(err.Error())
	}
}
