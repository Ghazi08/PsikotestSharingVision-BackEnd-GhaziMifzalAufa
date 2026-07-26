package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

type Post struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Content     string `json:"content"`
	Category    string `json:"category"`
	CreatedDate string `json:"created_date"`
	UpdatedDate string `json:"updated_date"`
	Status      string `json:"status"` // Publish | Draft | Thrash
}

var db *sql.DB

func main() {
	var err error

	dbHost := getEnv("MYSQLHOST", "127.0.0.1")
	dbPort := getEnv("MYSQLPORT", "3306")
	dbUser := getEnv("MYSQLUSER", "root")
	dbPass := getEnv("MYSQLPASSWORD", "")
	dbName := getEnv("MYSQLDATABASE", "db_article")

	// Koneksi awal untuk membuat database jika belum ada
	dsnNoDB := fmt.Sprintf("%s:%s@tcp(%s:%s)/?parseTime=true", dbUser, dbPass, dbHost, dbPort)
	tempDB, err := sql.Open("mysql", dsnNoDB)
	if err != nil {
		log.Fatalf("Gagal konek ke server MySQL: %v", err)
	}
	defer tempDB.Close()

	_, err = tempDB.Exec(fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %s;", dbName))
	if err != nil {
		log.Fatalf("Gagal membuat database: %v", err)
	}

	// Koneksi utama ke database
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", dbUser, dbPass, dbHost, dbPort, dbName)
	db, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("Gagal konek ke database: %v", err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Fatalf("Database tidak merespon: %v", err)
	}
	fmt.Println("Berhasil terhubung ke Database MySQL!")

	// Jalankan migrasi tabel
	runMigration()

	// Routing Endpoint CRUD
	http.HandleFunc("/posts", postsHandler)
	http.HandleFunc("/posts/", postDetailHandler)

	port := getEnv("PORT", "8080")
	fmt.Printf("Server berjalan di port %s...\n", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}

func runMigration() {
	query := `
	CREATE TABLE IF NOT EXISTS posts (
		id INT AUTO_INCREMENT PRIMARY KEY,
		title VARCHAR(200) NOT NULL,
		content TEXT NOT NULL,
		category VARCHAR(100) NOT NULL,
		created_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
		status VARCHAR(100) NOT NULL
	);`

	_, err := db.Exec(query)
	if err != nil {
		log.Fatalf("Migrasi tabel gagal: %v", err)
	}
	fmt.Println("Migrasi tabel 'posts' berhasil dijalankan!")
}

// Handler untuk GET (Semua Post) & POST (Tambah Post Baru)
func postsHandler(w http.ResponseWriter, r *http.Request) {
	setCORS(w, r)
	if r.Method == "OPTIONS" {
		return
	}

	switch r.Method {
	case "GET":
		rows, err := db.Query("SELECT id, title, content, category, created_date, updated_date, status FROM posts")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var posts []Post
		for rows.Next() {
			var p Post
			var created, updated sql.NullTime
			if err := rows.Scan(&p.ID, &p.Title, &p.Content, &p.Category, &created, &updated, &p.Status); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			if created.Valid {
				p.CreatedDate = created.Time.Format("2006-01-02 15:04:05")
			}
			if updated.Valid {
				p.UpdatedDate = updated.Time.Format("2006-01-02 15:04:05")
			}
			posts = append(posts, p)
		}
		json.NewEncoder(w).Encode(posts)

	case "POST":
		var p Post
		if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		query := "INSERT INTO posts (title, content, category, status) VALUES (?, ?, ?, ?)"
		result, err := db.Exec(query, p.Title, p.Content, p.Category, p.Status)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		id, _ := result.LastInsertId()
		p.ID = int(id)
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(p)

	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

// Handler untuk Detail Post (PUT/Update & DELETE berdasarkan ID)
func postDetailHandler(w http.ResponseWriter, r *http.Request) {
	setCORS(w, r)
	if r.Method == "OPTIONS" {
		return
	}

	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 3 || parts[2] == "" {
		http.NotFound(w, r)
		return
	}
	idStr := parts[2]

	switch r.Method {
	case "PUT":
		var p Post
		if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		query := "UPDATE posts SET title = ?, content = ?, category = ?, status = ? WHERE id = ?"
		_, err := db.Exec(query, p.Title, p.Content, p.Category, p.Status, idStr)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(map[string]string{"message": "Post updated successfully"})

	case "DELETE":
		query := "DELETE FROM posts WHERE id = ?"
		_, err := db.Exec(query, idStr)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(map[string]string{"message": "Post deleted successfully"})

	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func setCORS(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
}