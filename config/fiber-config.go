package config

import (
	"time"

	"github.com/gofiber/fiber/v2"
)

// ------------ FIBER CONFIG ------------
type FiberConfig struct {
	Prefork                      bool     `mapstructure:"Prefork"`
	ServerHeader                 string   `mapstructure:"ServerHeader"`
	StrictRouting                bool     `mapstructure:"StrictRouting"`
	CaseSensitive                bool     `mapstructure:"CaseSensitive"`
	Immutable                    bool     `mapstructure:"Immutable"`
	UnescapePath                 bool     `mapstructure:"UnescapePath"`
	ETag                         bool     `mapstructure:"ETag"`
	BodyLimit                    int      `mapstructure:"BodyLimit"`
	Concurrency                  int      `mapstructure:"Concurrency"`
	ReadTimeout                  int      `mapstructure:"ReadTimeout"`
	WriteTimeout                 int      `mapstructure:"WriteTimeout"`
	IdleTimeout                  int      `mapstructure:"IdleTimeout"`
	ReadBufferSize               int      `mapstructure:"ReadBufferSize"`
	WriteBufferSize              int      `mapstructure:"WriteBufferSize"`
	CompressedFileSuffix         string   `mapstructure:"CompressedFileSuffix"`
	DisableKeepalive             bool     `mapstructure:"DisableKeepalive"`
	DisableDefaultDate           bool     `mapstructure:"DisableDefaultDate"`
	DisableDefaultContentType    bool     `mapstructure:"DisableDefaultContentType"`
	DisableHeaderNormalizing     bool     `mapstructure:"DisableHeaderNormalizing"`
	DisableStartupMessage        bool     `mapstructure:"DisableStartupMessage"`
	AppName                      string   `mapstructure:"AppName"`
	StreamRequestBody            bool     `mapstructure:"StreamRequestBody"`
	DisablePreParseMultipartForm bool     `mapstructure:"DisablePreParseMultipartForm"`
	ReduceMemoryUsage            bool     `mapstructure:"ReduceMemoryUsage"`
	Network                      string   `mapstructure:"Network"`
	EnableTrustedProxyCheck      bool     `mapstructure:"EnableTrustedProxyCheck"`
	TrustedProxies               []string `mapstructure:"TrustedProxies"`
	EnableIPValidation           bool     `mapstructure:"EnableIPValidation"`
	EnablePrintRoutes            bool     `mapstructure:"EnablePrintRoutes"`
	ProxyHeader                  string   `mapstructure:"ProxyHeader"`
	GETOnly                      bool     `mapstructure:"GETOnly"`
	RequestMethods               []string `mapstructure:"RequestMethods"`
	EnableSplittingOnParsers     bool     `mapstructure:"EnableSplittingOnParsers"`
}

// InitFiberConfig создает конфигурацию Fiber из глобальной конфигурации.
func InitFiberConfig(cfg *Config) fiber.Config {
	return fiber.Config{
		Prefork:                      cfg.FIBER.Prefork,
		ServerHeader:                 cfg.FIBER.ServerHeader,
		StrictRouting:                cfg.FIBER.StrictRouting,
		CaseSensitive:                cfg.FIBER.CaseSensitive,
		Immutable:                    cfg.FIBER.Immutable,
		UnescapePath:                 cfg.FIBER.UnescapePath,
		ETag:                         cfg.FIBER.ETag,
		BodyLimit:                    cfg.FIBER.BodyLimit,
		Concurrency:                  cfg.FIBER.Concurrency,
		ReadTimeout:                  time.Duration(cfg.FIBER.ReadTimeout),
		WriteTimeout:                 time.Duration(cfg.FIBER.WriteTimeout),
		IdleTimeout:                  time.Duration(cfg.FIBER.IdleTimeout),
		ReadBufferSize:               cfg.FIBER.ReadBufferSize,
		WriteBufferSize:              cfg.FIBER.WriteBufferSize,
		CompressedFileSuffix:         cfg.FIBER.CompressedFileSuffix,
		DisableKeepalive:             cfg.FIBER.DisableKeepalive,
		DisableDefaultDate:           cfg.FIBER.DisableDefaultDate,
		DisableDefaultContentType:    cfg.FIBER.DisableDefaultContentType,
		DisableHeaderNormalizing:     cfg.FIBER.DisableHeaderNormalizing,
		DisableStartupMessage:        cfg.FIBER.DisableStartupMessage,
		AppName:                      cfg.FIBER.AppName,
		StreamRequestBody:            cfg.FIBER.StreamRequestBody,
		DisablePreParseMultipartForm: cfg.FIBER.DisablePreParseMultipartForm,
		ReduceMemoryUsage:            cfg.FIBER.ReduceMemoryUsage,
		Network:                      cfg.FIBER.Network,
		EnableTrustedProxyCheck:      cfg.FIBER.EnableTrustedProxyCheck,
		TrustedProxies:               cfg.FIBER.TrustedProxies,
		EnableIPValidation:           cfg.FIBER.EnableIPValidation,
		EnablePrintRoutes:            cfg.FIBER.EnablePrintRoutes,
		ProxyHeader:                  cfg.FIBER.ProxyHeader,
		GETOnly:                      cfg.FIBER.GETOnly,
		RequestMethods:               cfg.FIBER.RequestMethods,
		EnableSplittingOnParsers:     cfg.FIBER.EnableSplittingOnParsers,
	}
}
