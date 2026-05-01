package models

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

type AnalyticsBooking struct {
	ID          int64             `json:"id" db:"id"`
	Name        string            `json:"name" db:"name"`
	Description string            `json:"description" db:"description"`
	Status      string            `json:"status" db:"status"`
	Type        string            `json:"type" db:"type"`
	Priority    int               `json:"priority" db:"priority"`
	Amount      float64           `json:"amount" db:"amount"`
	IsActive    bool              `json:"is_active" db:"is_active"`
	Email       string            `json:"email" db:"email"`
	Phone       string            `json:"phone" db:"phone"`
	Address     string            `json:"address" db:"address"`
	City        string            `json:"city" db:"city"`
	Country     string            `json:"country" db:"country"`
	Code        string            `json:"code" db:"code"`
	Reference   string            `json:"reference" db:"reference"`
	Notes       string            `json:"notes" db:"notes"`
	Version     int               `json:"version" db:"version"`
	Metadata    map[string]string `json:"metadata" db:"-"`
	CreatedAt   time.Time         `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time         `json:"updated_at" db:"updated_at"`
}

type AnalyticsBookingFilter struct {
	IDs       []int64   `json:"ids"`
	Status    string    `json:"status"`
	Type      string    `json:"type"`
	Search    string    `json:"search"`
	MinAmount float64   `json:"min_amount"`
	MaxAmount float64   `json:"max_amount"`
	Country   string    `json:"country"`
	IsActive  *bool     `json:"is_active"`
	FromDate  time.Time `json:"from_date"`
	ToDate    time.Time `json:"to_date"`
	SortBy    string    `json:"sort_by"`
	SortOrder string    `json:"sort_order"`
	Limit     int       `json:"limit"`
	Offset    int       `json:"offset"`
}

type AnalyticsBookingListResponse struct {
	Data       []*AnalyticsBooking `json:"data"`
	Total      int64    `json:"total"`
	Page       int      `json:"page"`
	PageSize   int      `json:"page_size"`
	TotalPages int      `json:"total_pages"`
	HasMore    bool     `json:"has_more"`
}

type AnalyticsBookingStats struct {
	TotalCount   int64   `json:"total_count"`
	ActiveCount  int64   `json:"active_count"`
	PendingCount int64   `json:"pending_count"`
	TotalAmount  float64 `json:"total_amount"`
	AvgAmount    float64 `json:"avg_amount"`
	MaxAmount    float64 `json:"max_amount"`
	MinAmount    float64 `json:"min_amount"`
}

func NewAnalyticsBooking(name string) *AnalyticsBooking {
	now := time.Now()
	return &AnalyticsBooking{
		Name:      name,
		Status:    "pending",
		IsActive:  true,
		Version:   1,
		Metadata:  make(map[string]string),
		CreatedAt: now,
		UpdatedAt: now,
	}
}

func (m *AnalyticsBooking) Validate() error {
	if strings.TrimSpace(m.Name) == "" {
		return fmt.Errorf("analytics_booking: name is required")
	}
	if len(m.Name) > 500 {
		return fmt.Errorf("analytics_booking: name exceeds maximum length of 500 characters")
	}
	if len(m.Description) > 5000 {
		return fmt.Errorf("analytics_booking: description exceeds maximum length of 5000 characters")
	}
	validStatuses := map[string]bool{
		"active": true, "inactive": true, "pending": true,
		"completed": true, "cancelled": true, "archived": true,
	}
	if m.Status != "" && !validStatuses[m.Status] {
		return fmt.Errorf("analytics_booking: invalid status %q", m.Status)
	}
	if m.Amount < 0 {
		return fmt.Errorf("analytics_booking: amount cannot be negative")
	}
	if m.Priority < 0 || m.Priority > 10 {
		return fmt.Errorf("analytics_booking: priority must be between 0 and 10")
	}
	if m.Email != "" && !strings.Contains(m.Email, "@") {
		return fmt.Errorf("analytics_booking: invalid email format")
	}
	if m.Version < 1 {
		return fmt.Errorf("analytics_booking: version must be at least 1")
	}
	return nil
}

func (m *AnalyticsBooking) ToJSON() ([]byte, error) {
	data, err := json.Marshal(m)
	if err != nil {
		return nil, fmt.Errorf("analytics_booking: marshal failed: %w", err)
	}
	return data, nil
}

func (m *AnalyticsBooking) FromJSON(data []byte) error {
	if err := json.Unmarshal(data, m); err != nil {
		return fmt.Errorf("analytics_booking: unmarshal failed: %w", err)
	}
	return nil
}

func (m *AnalyticsBooking) Clone() *AnalyticsBooking {
	clone := *m
	clone.ID = 0
	clone.Version = 1
	clone.Metadata = make(map[string]string)
	for k, v := range m.Metadata {
		clone.Metadata[k] = v
	}
	clone.CreatedAt = time.Now()
	clone.UpdatedAt = time.Now()
	return &clone
}

func (m *AnalyticsBooking) SetStatus(status string) error {
	transitions := map[string][]string{
		"pending":   {"active", "cancelled"},
		"active":    {"inactive", "completed", "suspended"},
		"inactive":  {"active", "archived"},
		"suspended": {"active", "cancelled"},
	}
	allowed, ok := transitions[m.Status]
	if !ok {
		return fmt.Errorf("analytics_booking: no transitions from status %q", m.Status)
	}
	for _, a := range allowed {
		if a == status {
			m.Status = status
			m.UpdatedAt = time.Now()
			m.IsActive = status == "active"
			return nil
		}
	}
	return fmt.Errorf("analytics_booking: transition from %q to %q not allowed", m.Status, status)
}

func (m *AnalyticsBooking) Sanitize() {
	m.Name = strings.TrimSpace(m.Name)
	m.Description = strings.TrimSpace(m.Description)
	m.Email = strings.ToLower(strings.TrimSpace(m.Email))
	m.Phone = strings.TrimSpace(m.Phone)
	m.Address = strings.TrimSpace(m.Address)
	m.City = strings.TrimSpace(m.City)
	m.Country = strings.ToUpper(strings.TrimSpace(m.Country))
	m.Code = strings.TrimSpace(m.Code)
	m.Reference = strings.TrimSpace(m.Reference)
	if m.Status == "" {
		m.Status = "pending"
	}
}

func (m *AnalyticsBooking) Merge(other *AnalyticsBooking) {
	if other.Name != "" {
		m.Name = other.Name
	}
	if other.Description != "" {
		m.Description = other.Description
	}
	if other.Status != "" {
		m.Status = other.Status
	}
	if other.Type != "" {
		m.Type = other.Type
	}
	if other.Email != "" {
		m.Email = other.Email
	}
	if other.Phone != "" {
		m.Phone = other.Phone
	}
	if other.Address != "" {
		m.Address = other.Address
	}
	if other.Amount != 0 {
		m.Amount = other.Amount
	}
	if other.Priority != 0 {
		m.Priority = other.Priority
	}
	for k, v := range other.Metadata {
		m.Metadata[k] = v
	}
	m.Version++
	m.UpdatedAt = time.Now()
}

func (m *AnalyticsBooking) String() string {
	return fmt.Sprintf("AnalyticsBooking[id=%d name=%q status=%s amount=%.2f]", m.ID, m.Name, m.Status, m.Amount)
}

func (m *AnalyticsBooking) MatchesFilter(f *AnalyticsBookingFilter) bool {
	if f.Status != "" && m.Status != f.Status {
		return false
	}
	if f.Type != "" && m.Type != f.Type {
		return false
	}
	if f.Country != "" && m.Country != f.Country {
		return false
	}
	if f.IsActive != nil && m.IsActive != *f.IsActive {
		return false
	}
	if f.MinAmount > 0 && m.Amount < f.MinAmount {
		return false
	}
	if f.MaxAmount > 0 && m.Amount > f.MaxAmount {
		return false
	}
	if f.Search != "" {
		search := strings.ToLower(f.Search)
		if !strings.Contains(strings.ToLower(m.Name), search) &&
			!strings.Contains(strings.ToLower(m.Description), search) {
			return false
		}
	}
	if !f.FromDate.IsZero() && m.CreatedAt.Before(f.FromDate) {
		return false
	}
	if !f.ToDate.IsZero() && m.CreatedAt.After(f.ToDate) {
		return false
	}
	return true
}

func FilterAnalyticsBookingList(items []*AnalyticsBooking, f *AnalyticsBookingFilter) *AnalyticsBookingListResponse {
	var filtered []*AnalyticsBooking
	for _, item := range items {
		if item.MatchesFilter(f) {
			filtered = append(filtered, item)
		}
	}
	total := len(filtered)
	if f.Offset > 0 && f.Offset < total {
		filtered = filtered[f.Offset:]
	}
	pageSize := f.Limit
	if pageSize <= 0 {
		pageSize = 20
	}
	if pageSize < len(filtered) {
		filtered = filtered[:pageSize]
	}
	totalPages := total / pageSize
	if total%pageSize != 0 {
		totalPages++
	}
	page := 1
	if f.Offset > 0 && pageSize > 0 {
		page = f.Offset/pageSize + 1
	}
	return &AnalyticsBookingListResponse{
		Data:       filtered,
		Total:      int64(total),
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
		HasMore:    f.Offset+pageSize < total,
	}
}

func CalculateAnalyticsBookingStats(items []*AnalyticsBooking) *AnalyticsBookingStats {
	stats := &AnalyticsBookingStats{}
	if len(items) == 0 {
		return stats
	}
	stats.MinAmount = items[0].Amount
	stats.MaxAmount = items[0].Amount
	for _, item := range items {
		stats.TotalCount++
		if item.IsActive {
			stats.ActiveCount++
		}
		if item.Status == "pending" {
			stats.PendingCount++
		}
		stats.TotalAmount += item.Amount
		if item.Amount > stats.MaxAmount {
			stats.MaxAmount = item.Amount
		}
		if item.Amount < stats.MinAmount {
			stats.MinAmount = item.Amount
		}
	}
	if stats.TotalCount > 0 {
		stats.AvgAmount = stats.TotalAmount / float64(stats.TotalCount)
	}
	return stats
}
