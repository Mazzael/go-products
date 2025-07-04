package main

import (
	"net/http"

	"github.com/Mazzael/go-api/configs"
	_ "github.com/Mazzael/go-api/docs"
	"github.com/Mazzael/go-api/internal/entity"
	"github.com/Mazzael/go-api/internal/infra/database"
	"github.com/Mazzael/go-api/internal/infra/webserver/handlers"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth"
	httpSwagger "github.com/swaggo/http-swagger"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// @title           Go API Example
// @version         1.0
// @description     Product API with auhtentication

// @contact.name   Vinicius Trujillo Mazza
// @contact.email  vtrujillomazza@gmail.com

// @host      localhost:8080
// @BasePath  /
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization

func main() {
	configs, err := configs.LoadConfig(".")
	if err != nil {
		panic("Failed to load configuration: " + err.Error())
	}

	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic("Failed to connect to the database: " + err.Error())
	}

	db.AutoMigrate(&entity.Product{}, &entity.User{})

	gormProductRepository := database.NewProduct(db)
	productHandler := handlers.NewProductHandler(gormProductRepository)

	gormUserRepository := database.NewUser(db)
	userHandler := handlers.NewUserHandler(gormUserRepository)

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.WithValue("token", configs.TokenAuth))
	r.Use(middleware.WithValue("JwtExpiresIn", configs.JWTExpiresIn))

	r.Route("/products", func(r chi.Router) {
		r.Use(jwtauth.Verifier(configs.TokenAuth))
		r.Use(jwtauth.Authenticator)
		r.Post("/", productHandler.CreateProduct)
		r.Get("/{id}", productHandler.GetProduct)
		r.Get("/", productHandler.GetProducts)
		r.Put("/{id}", productHandler.UpdateProduct)
		r.Delete("/{id}", productHandler.DeleteProduct)
	})

	r.Post("/users", userHandler.CreateUser)
	r.Post("/users/auth", userHandler.Login)

	r.Get("/docs/*", httpSwagger.Handler(httpSwagger.URL("http://localhost:8080/docs/doc.json")))

	http.ListenAndServe(":8080", r)
}
