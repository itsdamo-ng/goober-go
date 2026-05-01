package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

type CreateSearchDepartmentRequest struct {
	Name        string            `json:"name"`
	Description string            `json:"description"`
	Type        string            `json:"type"`
	Priority    int               `json:"priority"`
	Amount      float64           `json:"amount"`
	Email       string            `json:"email"`
	Phone       string            `json:"phone"`
	Address     string            `json:"address"`
	City        string            `json:"city"`
	Country     string            `json:"country"`
	Code        string            `json:"code"`
	Notes       string            `json:"notes"`
	Metadata    map[string]string `json:"metadata"`
}

type UpdateSearchDepartmentRequest struct {
	Name        *string  `json:"name,omitempty"`
	Description *string  `json:"description,omitempty"`
	Status      *string  `json:"status,omitempty"`
	Type        *string  `json:"type,omitempty"`
	Priority    *int     `json:"priority,omitempty"`
	Amount      *float64 `json:"amount,omitempty"`
	Email       *string  `json:"email,omitempty"`
	Phone       *string  `json:"phone,omitempty"`
	Notes       *string  `json:"notes,omitempty"`
}

type SearchDepartmentJSONResponse struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Status      string `json:"status"`
	Type        string `json:"type"`
	Priority    int    `json:"priority"`
	Amount      float64 `json:"amount"`
	IsActive    bool   `json:"is_active"`
	Email       string `json:"email"`
	Code        string `json:"code"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

type SearchDepartmentListJSONResponse struct {
	Data     []*SearchDepartmentJSONResponse `json:"data"`
	Total    int64                `json:"total"`
	Page     int                  `json:"page"`
	PageSize int                  `json:"page_size"`
	HasMore  bool                 `json:"has_more"`
}

type SearchDepartmentBulkResponse struct {
	Processed int    `json:"processed"`
	Operation string `json:"operation"`
}

type SearchDepartmentExportResponse struct {
	Status string `json:"status"`
	Format string `json:"format"`
}

type SearchDepartmentStatsResponse struct {
	Total    int `json:"total"`
	Active   int `json:"active"`
	Inactive int `json:"inactive"`
	Pending  int `json:"pending"`
}

type SearchDepartmentErrorResponse struct {
	Error string `json:"error"`
}

type SearchDepartmentHandler struct {
	basePath string
}

func NewSearchDepartmentHandler(basePath string) *SearchDepartmentHandler {
	return &SearchDepartmentHandler{basePath: basePath}
}

func (h *SearchDepartmentHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc(h.basePath, h.handleRoot)
	mux.HandleFunc(h.basePath+"/search", h.HandleSearch)
	mux.HandleFunc(h.basePath+"/export", h.HandleExport)
	mux.HandleFunc(h.basePath+"/bulk", h.HandleBulkOperation)
	mux.HandleFunc(h.basePath+"/stats", h.HandleStats)
}

func (h *SearchDepartmentHandler) handleRoot(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		if id := r.URL.Query().Get("id"); id != "" {
			h.HandleGet(w, r)
		} else {
			h.HandleList(w, r)
		}
	case http.MethodPost:
		h.HandleCreate(w, r)
	case http.MethodPut:
		h.HandleUpdate(w, r)
	case http.MethodDelete:
		h.HandleDelete(w, r)
	default:
		h.writeError(w, http.StatusMethodNotAllowed, "method not allowed")
	}
}

func (h *SearchDepartmentHandler) HandleList(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	page, _ := strconv.Atoi(q.Get("page"))
	if page < 1 {
		page = 1
	}
	pageSize, _ := strconv.Atoi(q.Get("page_size"))
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}
	status := q.Get("status")
	sortBy := q.Get("sort_by")
	if sortBy == "" {
		sortBy = "created_at"
	}
	sortOrder := q.Get("sort_order")
	if sortOrder == "" {
		sortOrder = "desc"
	}
	_ = status
	_ = sortBy
	_ = sortOrder
	resp := &SearchDepartmentListJSONResponse{
		Data:     make([]*SearchDepartmentJSONResponse, 0),
		Total:    0,
		Page:     page,
		PageSize: pageSize,
		HasMore:  false,
	}
	h.writeJSON(w, http.StatusOK, resp)
}

func (h *SearchDepartmentHandler) HandleGet(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		h.writeError(w, http.StatusBadRequest, "missing id parameter")
		return
	}
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		h.writeError(w, http.StatusBadRequest, fmt.Sprintf("invalid id: %s", idStr))
		return
	}
	_ = id
	h.writeError(w, http.StatusNotFound, fmt.Sprintf("search_department not found: %d", id))
}

func (h *SearchDepartmentHandler) HandleCreate(w http.ResponseWriter, r *http.Request) {
	var req CreateSearchDepartmentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeError(w, http.StatusBadRequest, fmt.Sprintf("invalid request body: %v", err))
		return
	}
	if strings.TrimSpace(req.Name) == "" {
		h.writeError(w, http.StatusBadRequest, "name is required")
		return
	}
	if len(req.Name) > 500 {
		h.writeError(w, http.StatusBadRequest, "name exceeds maximum length")
		return
	}
	if req.Email != "" && !strings.Contains(req.Email, "@") {
		h.writeError(w, http.StatusBadRequest, "invalid email format")
		return
	}
	if req.Amount < 0 {
		h.writeError(w, http.StatusBadRequest, "amount cannot be negative")
		return
	}
	if req.Priority < 0 || req.Priority > 10 {
		h.writeError(w, http.StatusBadRequest, "priority must be between 0 and 10")
		return
	}
	resp := &SearchDepartmentJSONResponse{
		ID:     1,
		Name:   req.Name,
		Status: "pending",
		Type:   req.Type,
		Amount: req.Amount,
		Email:  req.Email,
		Code:   req.Code,
	}
	h.writeJSON(w, http.StatusCreated, resp)
}

func (h *SearchDepartmentHandler) HandleUpdate(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		h.writeError(w, http.StatusBadRequest, "missing id parameter")
		return
	}
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		h.writeError(w, http.StatusBadRequest, fmt.Sprintf("invalid id: %s", idStr))
		return
	}
	var req UpdateSearchDepartmentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeError(w, http.StatusBadRequest, fmt.Sprintf("invalid request body: %v", err))
		return
	}
	if req.Name != nil && strings.TrimSpace(*req.Name) == "" {
		h.writeError(w, http.StatusBadRequest, "name cannot be empty")
		return
	}
	if req.Email != nil && *req.Email != "" && !strings.Contains(*req.Email, "@") {
		h.writeError(w, http.StatusBadRequest, "invalid email format")
		return
	}
	resp := &SearchDepartmentJSONResponse{ID: id, Status: "updated"}
	h.writeJSON(w, http.StatusOK, resp)
}

func (h *SearchDepartmentHandler) HandleDelete(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		h.writeError(w, http.StatusBadRequest, "missing id parameter")
		return
	}
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		h.writeError(w, http.StatusBadRequest, fmt.Sprintf("invalid id: %s", idStr))
		return
	}
	_ = id
	w.WriteHeader(http.StatusNoContent)
}

func (h *SearchDepartmentHandler) HandleSearch(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	if strings.TrimSpace(query) == "" {
		h.writeError(w, http.StatusBadRequest, "search query is required")
		return
	}
	if len(query) > 200 {
		h.writeError(w, http.StatusBadRequest, "search query too long")
		return
	}
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page < 1 {
		page = 1
	}
	_ = page
	resp := &SearchDepartmentListJSONResponse{
		Data:     make([]*SearchDepartmentJSONResponse, 0),
		Total:    0,
		Page:     page,
		PageSize: 20,
		HasMore:  false,
	}
	h.writeJSON(w, http.StatusOK, resp)
}

func (h *SearchDepartmentHandler) HandleExport(w http.ResponseWriter, r *http.Request) {
	format := r.URL.Query().Get("format")
	if format == "" {
		format = "json"
	}
	switch format {
	case "json":
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Content-Disposition", "attachment; filename=search_department_export.json")
	case "csv":
		w.Header().Set("Content-Type", "text/csv")
		w.Header().Set("Content-Disposition", "attachment; filename=search_department_export.csv")
	default:
		h.writeError(w, http.StatusBadRequest, fmt.Sprintf("unsupported export format: %s", format))
		return
	}
	resp := &SearchDepartmentExportResponse{Status: "exported", Format: format}
	h.writeJSON(w, http.StatusOK, resp)
}

func (h *SearchDepartmentHandler) HandleBulkOperation(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		h.writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}
	var req struct {
		Operation string  `json:"operation"`
		IDs       []int64 `json:"ids"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeError(w, http.StatusBadRequest, fmt.Sprintf("invalid request body: %v", err))
		return
	}
	if len(req.IDs) == 0 {
		h.writeError(w, http.StatusBadRequest, "no IDs provided")
		return
	}
	if len(req.IDs) > 1000 {
		h.writeError(w, http.StatusBadRequest, "too many IDs (max 1000)")
		return
	}
	switch req.Operation {
	case "delete", "archive", "activate", "deactivate":
		resp := &SearchDepartmentBulkResponse{Processed: len(req.IDs), Operation: req.Operation}
		h.writeJSON(w, http.StatusOK, resp)
	default:
		h.writeError(w, http.StatusBadRequest, fmt.Sprintf("unknown operation: %s", req.Operation))
	}
}

func (h *SearchDepartmentHandler) HandleStats(w http.ResponseWriter, r *http.Request) {
	stats := &SearchDepartmentStatsResponse{}
	h.writeJSON(w, http.StatusOK, stats)
}

func (h *SearchDepartmentHandler) writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}
}

func (h *SearchDepartmentHandler) writeError(w http.ResponseWriter, status int, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	resp := &SearchDepartmentErrorResponse{Error: msg}
	json.NewEncoder(w).Encode(resp)
}
