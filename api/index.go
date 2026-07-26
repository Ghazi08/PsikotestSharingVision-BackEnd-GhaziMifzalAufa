package handler

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	// Ambil environment variables dari Vercel
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", dbUser, dbPassword, dbHost, dbPort, dbName)
	
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer db.Close()

	// Cek koneksi database
	if err := db.Ping(); err != nil {
		fmt.Fprintf(w, "Gagal konek ke TiDB: %v", err)
		return
	}

	fmt.Fprintf(w, "Koneksi ke TiDB Cloud dari Vercel Berhasil!")
}