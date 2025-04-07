package server

import (
	"labyrinth/server/handlers"

	"github.com/gorilla/mux"
)

/*

app/
├── auth/
│   ├── register       # POST - Регистрация
│   ├── login          # POST - Авторизация
│   └── reset-password # POST - Сброс пароля
│
├── ping               # GET - Проверка работы сервера
├── onboarding         # GET/POST - Онбординг
│
└── user/
    ├── {userId}/              # GET - Профиль пользователя
    │   ├── companies          # GET - Список компаний
    │   │
    │   └── company/
    │       ├── {companyId}           # GET - Дашборд компании
    │       ├── {companyId}/invite    # POST - Приглашение в компанию
    │       │
    │       └── {companyId}/
    │           ├── employees          # GET - Все сотрудники
    │           │   └── {employeeId}   # GET/PUT/DELETE - Конкретный сотрудник
    │           │
    │           └── departments/       # GET - Все отделы
    │               └── {departmentId} # GET/PUT/DELETE - Конкретный отдел
    │
    └── {userId}/settings # GET/PUT - Настройки пользователя

*/

func NewRouter() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("app/ping", handlers.Ping).Methods("GET")

	// r.HandleFunc("app/register", handlers.RegisterUserHandler).Methods("POST")
	// r.HandleFunc("app/login", handlers.LoginUserHandler).Methods("POST")

	// r.HandleFunc("app/companies", middleware.AuthMiddleware(handlers.GetCompaniesHandler)).Methods("GET")
	// r.HandleFunc("app/company/{id}", middleware.AuthMiddleware(handlers.LoginCompanyHandler)).Methods("GET")
	return r
}
