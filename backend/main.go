package main

import (
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/blastrain/vitess-sqlparser/sqlparser"
)

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello from Go backend!")
}

func uploadSQLHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	err := r.ParseMultipartForm(10 << 20) // 10 MB max memory
	if err != nil {
		http.Error(w, "Error parsing multipart form", http.StatusBadRequest)
		return
	}
	file, header, err := r.FormFile("sqlfile")
	if err != nil {
		http.Error(w, "Error retrieving the file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	fmt.Printf("Uploaded File: %+v\n", header.Filename)
	fmt.Printf("File Size: %+v\n", header.Size)
	fmt.Printf("MIME Header: %+v\n", header.Header)

	sqlContent, err := io.ReadAll(file)
	if err != nil {
		http.Error(w, "Error reading the file", http.StatusInternalServerError)
		return
	}

	if processSQLFile(string(sqlContent)) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "SQL file processed successfully")
	} else {
		http.Error(w, "Error processing SQL file", http.StatusInternalServerError)
	}
}

func processSQLFile(sqlString string) bool {
	// Trim unnecessary whitespace
	sqlString = strings.TrimSpace(sqlString)

	// Check and remove BOM if present
	if strings.HasPrefix(sqlString, "\ufeff") {
		sqlString = strings.TrimPrefix(sqlString, "\ufeff")
	}

	// Log the raw content for debugging
	fmt.Printf("Processed SQL Content: %q\n", sqlString)

	// Parse the SQL string
	stmt, err := sqlparser.Parse(sqlString)
	if err != nil {
		fmt.Printf("Error parsing SQL: %v\n", err)
		return false
	}

	// Successfully parsed SQL
	fmt.Printf("Parsed Statement: %+v\n", stmt)
	return true
}

func main() {
	http.HandleFunc("/data", corsMiddleware(handler))
	http.HandleFunc("/api/upload-sql", corsMiddleware(uploadSQLHandler))
	http.ListenAndServe(":5000", nil)
}

func corsMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	}
}
