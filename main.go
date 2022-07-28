package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/S-S-Group/Vaccinator/src/controllers"
	"github.com/S-S-Group/Vaccinator/src/data"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {

	l := log.New(os.Stdout, "vaccinator", log.LstdFlags)

	var err error
	err = godotenv.Load()
	if err != nil {
		l.Fatalf("Error getting env, not comming through %v", err)
	}

	client := data.Connect2(os.Getenv("DB_DRIVER"), os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_PORT"), os.Getenv("DB_HOST"), os.Getenv("DB_NAME"), l)
	if client != nil {
		//data.Load(client, l)
		sm := controllers.Startup(l)

		s := &http.Server{
			Addr:         ":8080",
			Handler:      sm,
			ErrorLog:     l,
			ReadTimeout:  5 * time.Second,
			WriteTimeout: 10 * time.Second,
			IdleTimeout:  120 * time.Second,
		}

		go func() {
			l.Println("Starting server on port 8080")
			err := s.ListenAndServe()
			if err != nil {
				l.Printf("Error starting server: %s\n", err)
				os.Exit(1)
			}
		}()
		sigChan := make(chan os.Signal)
		signal.Notify(sigChan, os.Interrupt)
		signal.Notify(sigChan, os.Kill)

		sig := <-sigChan
		l.Println("Received terminate, graceful shutdown", sig)

		defer client.Close()
		tc, _ := context.WithTimeout(context.Background(), 30*time.Second)
		s.Shutdown(tc)
	}
}
