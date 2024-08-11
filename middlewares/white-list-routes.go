package middlewares

var WhiteListRoutes = []string{
	"/",
	"/ui",
	"/metrics",
	"/favicon.ico",
	"/debug/pprof",
	"/api/v1/docs",
	"/api/v1/docs/swagger.yaml",
	"/api/v1/auth/login",
	"/api/v1/auth/refresh-token",
}
