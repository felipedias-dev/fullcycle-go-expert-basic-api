package main

import (
	"net/http"

	"github.com/felipedias-dev/fullcycle-go-expert-basic-api/configs"
	"github.com/felipedias-dev/fullcycle-go-expert-basic-api/internal/entity"
	"github.com/felipedias-dev/fullcycle-go-expert-basic-api/internal/infra/database"
	"github.com/felipedias-dev/fullcycle-go-expert-basic-api/internal/infra/webserver/handlers"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

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

	productDB := database.NewProduct(db)
	productHandler := handlers.NewProductHandler(productDB)

	r.Post("/products", productHandler.CreateProduct)
	r.Get("/products", productHandler.GetProducts)
	r.Get("/products/{id}", productHandler.GetProduct)
	r.Put("/products/{id}", productHandler.UpdateProduct)
	r.Delete("/products/{id}", productHandler.DeleteProduct)

	userDB := database.NewUser(db)
	userHandler := handlers.NewUserHandler(userDB, config.TokenAuthKey, config.JWTExpiration)

	r.Post("/users", userHandler.CreateUser)
	r.Post("/users/auth", userHandler.GetJWTHandler)

	http.ListenAndServe(":8000", r)
}
