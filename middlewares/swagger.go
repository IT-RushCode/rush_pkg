package middlewares

import (
	"strings"

	"github.com/IT-RushCode/rush_pkg/config"
	"github.com/gofiber/contrib/swagger"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/basicauth"
)

func InitSwagger(app *fiber.App, cfg *config.Config) {
	swaggerCfg := cfg.SWAGGER
	basePath := swaggerCfg.BasePath
	if basePath == "" {
		basePath = "/api/v1"
	}
	swaggerPath := strings.Trim(swaggerCfg.Path, "/")
	if swaggerPath == "" {
		swaggerPath = "docs"
	}
	title := swaggerCfg.Title
	if title == "" {
		title = cfg.APP.NAME + " API"
	}
	cacheAge := swaggerCfg.CacheAge
	if cacheAge == 0 {
		cacheAge = 3600
	}
	filePath := swaggerCfg.FilePath
	if filePath == "" {
		filePath = "./docs/swagger.yaml"
	}

	if swaggerCfg.BasicAuthEnabled && swaggerCfg.AuthUser != "" && swaggerCfg.AuthPassword != "" {
		auth := basicauth.New(basicauth.Config{
			Users: map[string]string{
				swaggerCfg.AuthUser: swaggerCfg.AuthPassword,
			},
			Realm: "Swagger",
		})
		cleanBase := strings.TrimRight(basePath, "/")
		if cleanBase == "" {
			cleanBase = "/"
		}
		path := strings.TrimRight(cleanBase+"/"+swaggerPath, "/")
		if path == "" {
			path = "/"
		}
		app.Use(path, auth)
		app.Use(path+"/", auth)
		app.Use(path+".yaml", auth)
	}

	app.Use(swagger.New(swagger.Config{
		BasePath: basePath,
		FilePath: filePath,
		Path:     swaggerPath,
		Title:    title,
		CacheAge: cacheAge,
	}))
}
