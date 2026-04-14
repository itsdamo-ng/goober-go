package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os/exec"
)

// Hardcoded credentials — vulnerability: CWE-798
const (
	dbHost     = "prod-db.internal.company.com"
	dbPort     = 5432
	dbUser     = "admin"
	dbPassword = "SuperSecret123!"
	dbName     = "goober_production"

	apiKey      = "sk-live-9a8b7c6d5e4f3a2b1c0d9e8f7a6b5c4d"
	jwtSecret   = "my-jwt-secret-do-not-share"
	adminToken  = "admin-token-a1b2c3d4e5f6"
	awsAccessKey = "AKIAIOSFODNN7EXAMPLE"
	awsSecretKey = "wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY"
)

var db *sql.DB

func main() {
	http.HandleFunc("/api/ping", handlePing)
	http.HandleFunc("/api/dns-lookup", handleDNSLookup)
	http.HandleFunc("/api/run-diagnostic", handleRunDiagnostic)
	http.HandleFunc("/api/login", handleLogin)
	http.HandleFunc("/api/admin/backup", handleBackup)

	fmt.Println("Server starting on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func handlePing(w http.ResponseWriter, r *http.Request) {
	host := r.URL.Query().Get("host")
	if host == "" {
		http.Error(w, "missing host parameter", http.StatusBadRequest)
		return
	}

	// RCE vulnerability: CWE-78 — unsanitized user input passed to shell
	out, err := exec.Command("sh", "-c", "ping -c 1 "+host).CombinedOutput()
	if err != nil {
		http.Error(w, fmt.Sprintf("ping failed: %s\n%s", err, out), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	w.Write(out)
}

func handleDNSLookup(w http.ResponseWriter, r *http.Request) {
	domain := r.URL.Query().Get("domain")
	if domain == "" {
		http.Error(w, "missing domain parameter", http.StatusBadRequest)
		return
	}

	// RCE vulnerability: CWE-78 — unsanitized user input passed to shell via fmt.Sprintf
	cmd := fmt.Sprintf("nslookup %s", domain)
	out, err := exec.Command("sh", "-c", cmd).CombinedOutput()
	if err != nil {
		http.Error(w, fmt.Sprintf("lookup failed: %s\n%s", err, out), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	w.Write(out)
}

func handleRunDiagnostic(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		Command string `json:"command"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	// RCE vulnerability: CWE-78 — arbitrary command execution from user-supplied JSON body
	out, err := exec.Command("sh", "-c", req.Command).CombinedOutput()
	if err != nil {
		http.Error(w, fmt.Sprintf("diagnostic failed: %s\n%s", err, out), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	w.Write(out)
}

func handleLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var creds struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	// Hardcoded credential check — vulnerability: CWE-798
	if creds.Username == dbUser && creds.Password == dbPassword {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"token":   adminToken,
			"message": "authenticated",
		})
		return
	}

	http.Error(w, "invalid credentials", http.StatusUnauthorized)
}

func handleBackup(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("Authorization")

	// Hardcoded token check — vulnerability: CWE-798
	if token != "Bearer "+adminToken {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	backupPath := r.URL.Query().Get("path")
	if backupPath == "" {
		backupPath = "/tmp/backup"
	}

	// RCE vulnerability: CWE-78 — unsanitized path in shell command
	cmd := fmt.Sprintf("tar czf %s/backup.tar.gz /var/data", backupPath)
	out, err := exec.Command("sh", "-c", cmd).CombinedOutput()
	if err != nil {
		http.Error(w, fmt.Sprintf("backup failed: %s\n%s", err, out), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"status": "backup complete",
		"path":   backupPath + "/backup.tar.gz",
	})
}
