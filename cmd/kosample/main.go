package main

import (
	"context"
	"kosample"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/google/uuid"
)

func init() {
	db := kosample.GetConnect()
	defer db.Close()

	tableExistSQL := "SELECT EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name = 'users');"
	var result string
	err := db.Get(&result, tableExistSQL)
	if err != nil {
		panic(err)
	}
	// ユーザテーブルがなければ作る
	if result == "false" {
		createTableSQL := `
		CREATE TABLE users (id TEXT PRIMARY KEY, name TEXT);
		`
		_, err = db.Exec(createTableSQL)
		if err != nil {
			panic(err)
		}
		log.Println("テーブル: usersが作成されました。")

		users := []struct {
			ID   string
			Name string
		}{
			{uuid.NewString(), "Alice"},
			{uuid.NewString(), "Bob"},
		}

		// バルクインサートで初回のユーザを登録する。NamedExecを利用するときは、構造体変数のタグにDBタグを利用すること。
		query := `INSERT INTO users(id, name) VALUES(:id, :name)`
		_, err := db.NamedExec(query, users)
		if err != nil {
			panic(err)
		}
	}

}

func main() {
	server := kosample.NewServer()
	defer server.DBPool.Close()
	ur := kosample.NewUserResource(server.DBPool)
	server.Handler = ur.Router()

	// Server run context
	serverCtx, serverStopCtx := context.WithCancel(context.Background())

	// Listen for syscall signals for process to interrupt/quit
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	go func() {
		<-sig

		// Shutdown signal with grace period of 30 seconds
		shutdownCtx, cancel := context.WithTimeout(serverCtx, 30*time.Second)
		defer cancel()

		go func() {
			<-shutdownCtx.Done()
			if shutdownCtx.Err() == context.DeadlineExceeded {
				log.Fatal("graceful shutdown timed out.. forcing exit.")
			}
		}()

		// Trigger graceful shutdown
		err := server.Shutdown(shutdownCtx)
		if err != nil {
			log.Fatal(err)
		}
		serverStopCtx()
	}()

	// Run the server
	err := server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}

	// Wait for server context to be stopped
	<-serverCtx.Done()
}
