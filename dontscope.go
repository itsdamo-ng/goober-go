package main

import (
	"crypto/md5"
	"database/sql"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

func init() {
	http.HandleFunc("/api/user", handleUserLookup)
	http.HandleFunc("/api/search", handleSearch)
	http.HandleFunc("/api/file", handleFileRead)
	http.HandleFunc("/api/fetch", handleFetch)
	http.HandleFunc("/api/register", handleRegister)
}

// SQL Injection — CWE-89: unsanitized user input concatenated into query
func handleUserLookup(w http.ResponseWriter, r *http.Request) {
	username := r.URL.Query().Get("username")
	if username == "" {
		http.Error(w, "missing username parameter", http.StatusBadRequest)
		return
	}

	query := fmt.Sprintf("SELECT id, username, email FROM users WHERE username = '%s'", username)
	rows, err := db.Query(query)
	if err != nil {
		http.Error(w, fmt.Sprintf("query failed: %s", err), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var results []map[string]string
	for rows.Next() {
		var id, uname, email string
		if err := rows.Scan(&id, &uname, &email); err != nil {
			continue
		}
		results = append(results, map[string]string{"id": id, "username": uname, "email": email})
	}

	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, `{"results": %d}`, len(results))
}

// Reflected XSS — CWE-79: user input rendered directly into HTML response
func handleSearch(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	if query == "" {
		http.Error(w, "missing q parameter", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "text/html")
	fmt.Fprintf(w, `<html><body>
		<h1>Search Results</h1>
		<p>You searched for: %s</p>
		<p>No results found.</p>
	</body></html>`, query)
}

// Path Traversal — CWE-22: unsanitized file path allows reading arbitrary files
func handleFileRead(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	if name == "" {
		http.Error(w, "missing name parameter", http.StatusBadRequest)
		return
	}

	filePath := filepath.Join("/var/data/reports", name)
	data, err := os.ReadFile(filePath)
	if err != nil {
		http.Error(w, fmt.Sprintf("file not found: %s", err), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	w.Write(data)
}

// SSRF — CWE-918: server makes HTTP request to user-supplied URL
func handleFetch(w http.ResponseWriter, r *http.Request) {
	target := r.URL.Query().Get("url")
	if target == "" {
		http.Error(w, "missing url parameter", http.StatusBadRequest)
		return
	}

	resp, err := http.Get(target)
	if err != nil {
		http.Error(w, fmt.Sprintf("fetch failed: %s", err), http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	w.Header().Set("Content-Type", resp.Header.Get("Content-Type"))
	w.Write(body)
}

// Weak Crypto — CWE-327: MD5 used for password hashing
func handleRegister(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	username := r.FormValue("username")
	password := r.FormValue("password")
	if username == "" || password == "" {
		http.Error(w, "missing username or password", http.StatusBadRequest)
		return
	}

	hash := md5.Sum([]byte(password))
	passwordHash := hex.EncodeToString(hash[:])

	_, err := db.Exec(
		fmt.Sprintf("INSERT INTO users (username, password_hash) VALUES ('%s', '%s')", username, passwordHash),
	)
	if err != nil {
		http.Error(w, fmt.Sprintf("registration failed: %s", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, `{"status": "registered", "username": "%s"}`, username)
}

// Ensure db is usable from this file (declared in main.go)
var _ *sql.DB = db
