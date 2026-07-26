package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	// Handler utama untuk route /
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		dbHost := os.Getenv("DB_HOST")
		dbPort := os.Getenv("DB_PORT")
		dbUser := os.Getenv("DB_USER")
		dbPassword := os.Getenv("DB_PASSWORD")
		dbName := os.Getenv("DB_NAME")

		dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", dbUser, dbPassword, dbHost, dbPort, dbName)
		
		db, err := sql.Open("mysql", dsn)
		if err != nil {
			http.Error(w, "DB Open Error: "+err.Error(), http.StatusInternalServerError)
			return
		}
		defer db.Close()

		if err := db.Ping(); err != nil {
			fmt.Fprintf(w, "Gagal konek ke TiDB Cloud: %v", err)
			return
		}

		fmt.Fprintf(w, "SUKSES! Backend Go di Vercel terhubung ke TiDB Cloud!")
	})

	// Vercel menyediakan port secara dinamis lewat os.Getenv("PORT")
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Printf("Server berjalan di port %s\n", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		fmt.Printf("Server error: %v\n", err)
	}
}
