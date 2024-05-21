package config

// ------------ GLOBAL CONFIG ------------
type Config struct {
	APP      AppConfig      `json:"app"`
	DB       DBConfig       `json:"db"`
	JWT      JwtConfig      `json:"jwt"`
	SERVICES ServicesConfig `json:"services"`
	DATETIME DateTimeConfig `json:"datetime"`
	REDIS    RedisConfig    `json:"redis"`
	KAFKA    KafKaConfig    `json:"kafka"`
}

// ------------ SERVICES ------------
type AppConfig struct {
	ENV            string `json:"env"`
	DEBUG          bool   `json:"debug"`
	NAME           string `json:"name"`
	HOST           string `json:"host"`
	PORT           string `json:"port"`
	MAX_CONNECTION string `json:"max_connection"`
}

type ServiceConfig struct {
	NAME string `json:"name"`
	HOST string `json:"host"`
	PORT string `json:"port"`
}

type ServicesConfig struct {
	AUTH         ServiceConfig `json:"auth"`
	COLLECTIONS  ServiceConfig `json:"collections"`
	CLIENT       ServiceConfig `json:"client"`
	PLEDGE       ServiceConfig `json:"pledge"`
	LOAN         ServiceConfig `json:"loan"`
	CALCULATOR   ServiceConfig `json:"calculator"`
	NOTIFICATION ServiceConfig `json:"notification"`
	STORAGE      ServiceConfig `json:"storage"`
	SCORING      ServiceConfig `json:"scoring"`
	BP           ServiceConfig `json:"bp"`
	FIN_ANALYS   ServiceConfig `json:"fin_analys"`
}

// ------------ JWT ------------
type JwtConfig struct {
	JWT_SECRET  string `json:"jwt_secret"`
	JWT_TTL     int64  `json:"jwt_ttl"`
	REFRESH_TTL int64  `json:"refresh_ttl"`
}

// ------------ KAFKA ------------
type KafKaConfig struct {
	HOST1 string `json:"host1"`
	HOST2 string `json:"host2"`
	HOST3 string `json:"host3"`
}

// ------------ DATETIME ------------
type DateTimeConfig struct {
	Datetime string `json:"datetime"`
	Date     string `json:"date"`
	Time     string `json:"time"`
}

// ------------ REDIS ------------
type RedisConfig struct {
	HOST string `json:"host"`
	PORT string `json:"port"`
	PASS string `json:"password"`
	DB   int    `json:"db"`
}

// ------------ MONGO DB ------------
type MongoDBConfig struct {
	URI string `json:"uri"`
}

// ------------ DATABASE ------------
type DatabaseConfig struct {
	Title            string `json:"title"`
	MigrationEnabled bool   `json:"migration_enabled"`
	Host             string `json:"dbhost"`
	Port             int64  `json:"dbport"`
	User             string `json:"dbuser"`
	Pass             string `json:"dbpass"`
	Name             string `json:"dbname"`
	CHARSET          string `json:"charset"`
}

type DBConfig struct {
	MSSQL          DatabaseConfig `json:"MSSQL"`
	MYSQL          DatabaseConfig `json:"MYSQL"`
	PSQL           DatabaseConfig `json:"PSQL"`
	MSSQL_ABS      DatabaseConfig `json:"MSSQL_ABS"`
	MSSQL_CONVEYOR DatabaseConfig `json:"MSSQL_CONVEYOR"`
}
