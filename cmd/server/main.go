package main

import (
	"net/http"

	"github.com/felipedias-dev/fullcycle-go-expert-basic-api/configs"
	_ "github.com/felipedias-dev/fullcycle-go-expert-basic-api/docs"
	"github.com/felipedias-dev/fullcycle-go-expert-basic-api/internal/entity"
	"github.com/felipedias-dev/fullcycle-go-expert-basic-api/internal/infra/database"
	"github.com/felipedias-dev/fullcycle-go-expert-basic-api/internal/infra/webserver/handlers"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/jwtauth"
	httpSwagger "github.com/swaggo/http-swagger"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// @title Full Cycle Go Expert Basic API
// @description Product API with Authentication
// @version 1.0.0
// @license.name MIT
// @license.url https://github.com/git/git-scm.com/blob/main/MIT-LICENSE.txt
// @termsOfService http://swagger.io/terms/
// @host localhost:8000
// @BasePath /
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func main() {
	config, err := configs.LoadConfig(".")
	if err != nil {
		panic(err)
	}
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	db.AutoMigrate(&entity.Product{}, &entity.User{})

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.WithValue("jwt", config.TokenAuthKey))
	r.Use(middleware.WithValue("jwtExpiration", config.JWTExpiration))

	productDB := database.NewProduct(db)
	productHandler := handlers.NewProductHandler(productDB)

	r.Route("/products", func(r chi.Router) {
		r.Use(jwtauth.Verifier(config.TokenAuthKey))
		r.Use(jwtauth.Authenticator)
		r.Post("/", productHandler.CreateProduct)
		r.Get("/", productHandler.GetProducts)
		r.Get("/{id}", productHandler.GetProduct)
		r.Put("/{id}", productHandler.UpdateProduct)
		r.Delete("/{id}", productHandler.DeleteProduct)
	})

	userDB := database.NewUser(db)
	userHandler := handlers.NewUserHandler(userDB)

	r.Post("/users", userHandler.CreateUser)
	r.Post("/users/auth", userHandler.GetJWTHandler)

	r.Get("/docs/*", httpSwagger.Handler(httpSwagger.URL("http://localhost:8000/docs/doc.json")))

	http.ListenAndServe(":8000", r)
}
