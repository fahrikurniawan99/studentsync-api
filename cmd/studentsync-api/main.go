package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres" // Ganti dengan driver database Anda
	"gorm.io/gorm"

	"github.com/fahrikurniawan99/studentsync-api/internal/app/handler"
	"github.com/fahrikurniawan99/studentsync-api/internal/app/repository"
	"github.com/fahrikurniawan99/studentsync-api/internal/app/service"
	"github.com/fahrikurniawan99/studentsync-api/internal/config"
	"github.com/fahrikurniawan99/studentsync-api/internal/model"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Gagal memuat konfigurasi: %v", err)
	}

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.DB_HOST,
		cfg.DB_PORT,
		cfg.DB_USER,
		cfg.DB_PASSWORD,
		cfg.DB_NAME,
	)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Gagal terhubung ke database: %v", err)
	}
	log.Println("Berhasil terhubung ke database!")

	err = db.AutoMigrate(&model.User{})
	if err != nil {
		log.Fatalf("Gagal melakukan migrasi database: %v", err)
	}
	log.Println("Migrasi database berhasil!")

	userRepo := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepo)
	userHandler := handler.NewUserHandler(userService)

	router := gin.Default()

	userHandler.RegisterRoutes(router)

	serverAddr := ":" + cfg.PORT
	fmt.Printf("Server berjalan di %s\n", serverAddr)
	log.Fatal(http.ListenAndServe(serverAddr, router))
}
