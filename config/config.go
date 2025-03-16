package config

import "time"

type Config struct {
	Network    Network    `json:"network"`
	PostgreSQL PostgreSQL `json:"postgresql"`
	Mongo      Mongo      `json:"mongo"`
	Redis      Redis      `json:"redis"`
	Minio      Minio      `json:"minio"`
}

type Network struct {
	IP            string        `json:"ip"`
	Port          string        `json:"port"`
	PathToTLSCert string        `json:"path_to_tls_cert"`
	PathToTLSKey  string        `json:"path_to_tls_key"`
	KeepAlive     bool          `json:"keep_alive"`
	Timeout       time.Duration `json:"timeout"`
	AliveTime     time.Duration `json:"alive_time"`
}

type PostgreSQL struct {
	Host         string `json:"host"`
	Port         string `json:"port"`
	Username     string `json:"username"`
	Password     string `json:"password"`
	DatabaseName string `json:"database_name"`
	SSLMode      bool   `json:"ssl_mode"`
}

type Mongo struct {
	URI          string `json:"uri"`
	DatabaseName string `json:"database_name"`
	Username     string `json:"username"`
	Password     string `json:"password"`
	AuthSource   string `json:"auth_source"`
}

type Redis struct {
	Addr     string `json:"addr"`
	Password string `json:"password"`
	DB       int    `json:"db"`
}

type Minio struct {
	Endpoint  string `json:"endpoint"`
	AccessKey string `json:"access_key"`
	SecretKey string `json:"secret_key"`
	UseSSL    bool   `json:"use_ssl"`
	Bucket    string `json:"bucket"`
}
