package config

import "time"

var Conf = &Config{
	Network: Network{
		IP:            "127.0.0.1",         // IP-адрес сервера
		Port:          "8080",              // Порт сервера
		PathToTLSCert: "/path/to/tls/cert", // Путь к TLS-сертификату
		PathToTLSKey:  "/path/to/tls/key",  // Путь к TLS-ключу
		KeepAlive:     true,                // Включение keep-alive
		Timeout:       30 * time.Second,    // Таймаут соединения
		AliveTime:     60 * time.Second,    // Время жизни соединения
	},
	PostgreSQL: PostgreSQL{
		Host:         "localhost", // Хост PostgreSQL
		Port:         "5432",      // Порт PostgreSQL
		Username:     "postgres",  // Имя пользователя
		Password:     "0000",      // Пароль
		DatabaseName: "labyrinth", // Имя базы данных
		SSLMode:      "disable",   // Режим SSL
	},
	Mongo: Mongo{
		URI:          "mongodb://localhost:27017", // URI MongoDB
		DatabaseName: "mydb",                      // Имя базы данных
		Username:     "user",                      // Имя пользователя
		Password:     "0000",                      // Пароль
		AuthSource:   "admin",                     // Источник аутентификации
	},
	Redis: Redis{
		Addr:     "localhost:6379", // Адрес Redis
		Password: "password",       // Пароль Redis
		DB:       0,                // Номер базы данных
	},
	Minio: Minio{
		Endpoint:  "localhost:9000", // Адрес MinIO
		AccessKey: "minioadmin",     // Ключ доступа
		SecretKey: "minioadmin",     // Секретный ключ
		UseSSL:    false,            // Использование SSL
		Bucket:    "mybucket",       // Имя бакета
	},
}

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
	SSLMode      string `json:"ssl_mode"`
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
