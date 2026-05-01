package services

import (
	"fmt"
	"sort"
	"strings"
	"time"
)

type SupportDiscountRecord struct {
	ID          int64
	Name        string
	Description string
	Status      string
	Type        string
	Priority    int
	Amount      float64
	IsActive    bool
	Email       string
	Phone       string
	Address     string
	City        string
	Country     string
	Code        string
	Reference   string
	Notes       string
	Version     int
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type CreateSupportDiscountInput struct {
	Name        string
	Description string
	Type        string
	Priority    int
	Amount      float64
	Email       string
	Phone       string
	Address     string
	City        string
	Country     string
	Code        string
	Notes       string
}

type UpdateSupportDiscountInput struct {
	Name        *string
	Description *string
	Status      *string
	Type        *string
	Priority    *int
	Amount      *float64
	Email       *string
	Phone       *string
	Notes       *string
}

type SupportDiscountSearchParams struct {
	Query     string
	Status    string
	Type      string
	Country   string
	MinAmount float64
	MaxAmount float64
	SortBy    string
	SortOrder string
	Page      int
	PageSize  int
}

type SupportDiscountProcessResult struct {
	Processed int
	Succeeded int
	Failed    int
	Errors    []string
	Duration  time.Duration
}

type SupportDiscountService struct {
	records map[int64]*SupportDiscountRecord
	nextID  int64
}

func NewSupportDiscountService() *SupportDiscountService {
	return &SupportDiscountService{
		records: make(map[int64]*SupportDiscountRecord),
		nextID:  1,
	}
}

func (s *SupportDiscountService) Create(input *CreateSupportDiscountInput) (*SupportDiscountRecord, error) {
	if strings.TrimSpace(input.Name) == "" {
		return nil, fmt.Errorf("support_discount: name is required")
	}
	if input.Amount < 0 {
		return nil, fmt.Errorf("support_discount: amount cannot be negative")
	}
	if input.Email != "" && !strings.Contains(input.Email, "@") {
		return nil, fmt.Errorf("support_discount: invalid email format")
	}
	now := time.Now()
	record := &SupportDiscountRecord{
		ID:          s.nextID,
		Name:        strings.TrimSpace(input.Name),
		Description: strings.TrimSpace(input.Description),
		Status:      "pending",
		Type:        input.Type,
		Priority:    input.Priority,
		Amount:      input.Amount,
		IsActive:    true,
		Email:       strings.ToLower(strings.TrimSpace(input.Email)),
		Phone:       strings.TrimSpace(input.Phone),
		Address:     strings.TrimSpace(input.Address),
		City:        strings.TrimSpace(input.City),
		Country:     strings.ToUpper(strings.TrimSpace(input.Country)),
		Code:        input.Code,
		Notes:       input.Notes,
		Version:     1,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	s.records[record.ID] = record
	s.nextID++
	return record, nil
}

func (s *SupportDiscountService) GetByID(id int64) (*SupportDiscountRecord, error) {
	record, ok := s.records[id]
	if !ok {
		return nil, fmt.Errorf("support_discount: record %d not found", id)
	}
	return record, nil
}

func (s *SupportDiscountService) Update(id int64, input *UpdateSupportDiscountInput) (*SupportDiscountRecord, error) {
	record, ok := s.records[id]
	if !ok {
		return nil, fmt.Errorf("support_discount: record %d not found", id)
	}
	if input.Name != nil {
		name := strings.TrimSpace(*input.Name)
		if name == "" {
			return nil, fmt.Errorf("support_discount: name cannot be empty")
		}
		record.Name = name
	}
	if input.Description != nil {
		record.Description = strings.TrimSpace(*input.Description)
	}
	if input.Status != nil {
		record.Status = *input.Status
		record.IsActive = *input.Status == "active"
	}
	if input.Type != nil {
		record.Type = *input.Type
	}
	if input.Priority != nil {
		record.Priority = *input.Priority
	}
	if input.Amount != nil {
		if *input.Amount < 0 {
			return nil, fmt.Errorf("support_discount: amount cannot be negative")
		}
		record.Amount = *input.Amount
	}
	if input.Email != nil {
		record.Email = strings.ToLower(strings.TrimSpace(*input.Email))
	}
	if input.Phone != nil {
		record.Phone = strings.TrimSpace(*input.Phone)
	}
	if input.Notes != nil {
		record.Notes = *input.Notes
	}
	record.Version++
	record.UpdatedAt = time.Now()
	return record, nil
}

func (s *SupportDiscountService) Delete(id int64) error {
	if _, ok := s.records[id]; !ok {
		return fmt.Errorf("support_discount: record %d not found", id)
	}
	delete(s.records, id)
	return nil
}

func (s *SupportDiscountService) List(params *SupportDiscountSearchParams) ([]*SupportDiscountRecord, int) {
	var results []*SupportDiscountRecord
	for _, r := range s.records {
		if params.Status != "" && r.Status != params.Status {
			continue
		}
		if params.Type != "" && r.Type != params.Type {
			continue
		}
		if params.Country != "" && r.Country != params.Country {
			continue
		}
		if params.MinAmount > 0 && r.Amount < params.MinAmount {
			continue
		}
		if params.MaxAmount > 0 && r.Amount > params.MaxAmount {
			continue
		}
		if params.Query != "" {
			q := strings.ToLower(params.Query)
			if !strings.Contains(strings.ToLower(r.Name), q) &&
				!strings.Contains(strings.ToLower(r.Description), q) &&
				!strings.Contains(strings.ToLower(r.Code), q) {
				continue
			}
		}
		results = append(results, r)
	}
	total := len(results)
	switch params.SortBy {
	case "name":
		sort.Slice(results, func(i, j int) bool {
			if params.SortOrder == "desc" {
				return results[i].Name > results[j].Name
			}
			return results[i].Name < results[j].Name
		})
	case "amount":
		sort.Slice(results, func(i, j int) bool {
			if params.SortOrder == "desc" {
				return results[i].Amount > results[j].Amount
			}
			return results[i].Amount < results[j].Amount
		})
	case "priority":
		sort.Slice(results, func(i, j int) bool {
			if params.SortOrder == "desc" {
				return results[i].Priority > results[j].Priority
			}
			return results[i].Priority < results[j].Priority
		})
	default:
		sort.Slice(results, func(i, j int) bool {
			if params.SortOrder == "asc" {
				return results[i].CreatedAt.Before(results[j].CreatedAt)
			}
			return results[i].CreatedAt.After(results[j].CreatedAt)
		})
	}
	pageSize := params.PageSize
	if pageSize <= 0 {
		pageSize = 20
	}
	offset := (params.Page - 1) * pageSize
	if offset < 0 {
		offset = 0
	}
	if offset >= len(results) {
		return nil, total
	}
	end := offset + pageSize
	if end > len(results) {
		end = len(results)
	}
	return results[offset:end], total
}

func (s *SupportDiscountService) BulkCreate(inputs []*CreateSupportDiscountInput) *SupportDiscountProcessResult {
	start := time.Now()
	result := &SupportDiscountProcessResult{}
	for _, input := range inputs {
		result.Processed++
		if _, err := s.Create(input); err != nil {
			result.Failed++
			result.Errors = append(result.Errors, err.Error())
		} else {
			result.Succeeded++
		}
	}
	result.Duration = time.Since(start)
	return result
}

func (s *SupportDiscountService) Archive(id int64) error {
	record, ok := s.records[id]
	if !ok {
		return fmt.Errorf("support_discount: record %d not found", id)
	}
	if record.Status == "archived" {
		return fmt.Errorf("support_discount: record %d is already archived", id)
	}
	record.Status = "archived"
	record.IsActive = false
	record.UpdatedAt = time.Now()
	return nil
}

func (s *SupportDiscountService) Restore(id int64) error {
	record, ok := s.records[id]
	if !ok {
		return fmt.Errorf("support_discount: record %d not found", id)
	}
	if record.Status != "archived" {
		return fmt.Errorf("support_discount: record %d is not archived", id)
	}
	record.Status = "active"
	record.IsActive = true
	record.UpdatedAt = time.Now()
	return nil
}

func (s *SupportDiscountService) Duplicate(id int64) (*SupportDiscountRecord, error) {
	original, ok := s.records[id]
	if !ok {
		return nil, fmt.Errorf("support_discount: record %d not found", id)
	}
	now := time.Now()
	clone := &SupportDiscountRecord{
		ID:          s.nextID,
		Name:        original.Name + " (copy)",
		Description: original.Description,
		Status:      "pending",
		Type:        original.Type,
		Priority:    original.Priority,
		Amount:      original.Amount,
		IsActive:    true,
		Email:       original.Email,
		Phone:       original.Phone,
		Address:     original.Address,
		City:        original.City,
		Country:     original.Country,
		Code:        original.Code,
		Notes:       original.Notes,
		Version:     1,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	s.records[clone.ID] = clone
	s.nextID++
	return clone, nil
}

func (s *SupportDiscountService) Count() int {
	return len(s.records)
}

func (s *SupportDiscountService) CountByStatus(status string) int {
	count := 0
	for _, r := range s.records {
		if r.Status == status {
			count++
		}
	}
	return count
}
