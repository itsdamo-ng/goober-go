package repository

import (
	"fmt"
	"strings"
	"time"
)

type ReportingBudgetEntity struct {
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

type ReportingBudgetQueryOptions struct {
	Where     []string
	Args      []interface{}
	OrderBy   string
	Direction string
	Limit     int
	Offset    int
}

type ReportingBudgetRepository struct {
	store  map[int64]*ReportingBudgetEntity
	nextID int64
}

func NewReportingBudgetRepository() *ReportingBudgetRepository {
	return &ReportingBudgetRepository{
		store:  make(map[int64]*ReportingBudgetEntity),
		nextID: 1,
	}
}

func (r *ReportingBudgetRepository) FindByID(id int64) (*ReportingBudgetEntity, error) {
	entity, ok := r.store[id]
	if !ok {
		return nil, fmt.Errorf("reporting_budget_repository: entity %d not found", id)
	}
	return entity, nil
}

func (r *ReportingBudgetRepository) FindAll(opts *ReportingBudgetQueryOptions) ([]*ReportingBudgetEntity, error) {
	var results []*ReportingBudgetEntity
	for _, entity := range r.store {
		results = append(results, entity)
	}
	if opts != nil {
		if opts.Offset > 0 && opts.Offset < len(results) {
			results = results[opts.Offset:]
		}
		if opts.Limit > 0 && opts.Limit < len(results) {
			results = results[:opts.Limit]
		}
	}
	return results, nil
}

func (r *ReportingBudgetRepository) FindByStatus(status string) ([]*ReportingBudgetEntity, error) {
	var results []*ReportingBudgetEntity
	for _, entity := range r.store {
		if entity.Status == status {
			results = append(results, entity)
		}
	}
	return results, nil
}

func (r *ReportingBudgetRepository) FindByCode(code string) (*ReportingBudgetEntity, error) {
	for _, entity := range r.store {
		if entity.Code == code {
			return entity, nil
		}
	}
	return nil, fmt.Errorf("reporting_budget_repository: entity with code %q not found", code)
}

func (r *ReportingBudgetRepository) FindByEmail(email string) (*ReportingBudgetEntity, error) {
	email = strings.ToLower(strings.TrimSpace(email))
	for _, entity := range r.store {
		if strings.ToLower(entity.Email) == email {
			return entity, nil
		}
	}
	return nil, fmt.Errorf("reporting_budget_repository: entity with email %q not found", email)
}

func (r *ReportingBudgetRepository) Save(entity *ReportingBudgetEntity) error {
	if entity.ID == 0 {
		entity.ID = r.nextID
		r.nextID++
		entity.CreatedAt = time.Now()
	}
	entity.UpdatedAt = time.Now()
	r.store[entity.ID] = entity
	return nil
}

func (r *ReportingBudgetRepository) Update(entity *ReportingBudgetEntity) error {
	if _, ok := r.store[entity.ID]; !ok {
		return fmt.Errorf("reporting_budget_repository: entity %d not found", entity.ID)
	}
	entity.UpdatedAt = time.Now()
	entity.Version++
	r.store[entity.ID] = entity
	return nil
}

func (r *ReportingBudgetRepository) Delete(id int64) error {
	if _, ok := r.store[id]; !ok {
		return fmt.Errorf("reporting_budget_repository: entity %d not found", id)
	}
	delete(r.store, id)
	return nil
}

func (r *ReportingBudgetRepository) SoftDelete(id int64) error {
	entity, ok := r.store[id]
	if !ok {
		return fmt.Errorf("reporting_budget_repository: entity %d not found", id)
	}
	entity.Status = "deleted"
	entity.IsActive = false
	entity.UpdatedAt = time.Now()
	return nil
}

func (r *ReportingBudgetRepository) Count() (int64, error) {
	return int64(len(r.store)), nil
}

func (r *ReportingBudgetRepository) CountByStatus(status string) (int64, error) {
	var count int64
	for _, entity := range r.store {
		if entity.Status == status {
			count++
		}
	}
	return count, nil
}

func (r *ReportingBudgetRepository) Exists(id int64) (bool, error) {
	_, ok := r.store[id]
	return ok, nil
}

func (r *ReportingBudgetRepository) ExistsByCode(code string) (bool, error) {
	for _, entity := range r.store {
		if entity.Code == code {
			return true, nil
		}
	}
	return false, nil
}

func (r *ReportingBudgetRepository) BulkInsert(entities []*ReportingBudgetEntity) error {
	for _, entity := range entities {
		if err := r.Save(entity); err != nil {
			return fmt.Errorf("reporting_budget_repository: bulk insert failed: %w", err)
		}
	}
	return nil
}

func (r *ReportingBudgetRepository) FindByDateRange(from, to time.Time) ([]*ReportingBudgetEntity, error) {
	var results []*ReportingBudgetEntity
	for _, entity := range r.store {
		if (from.IsZero() || !entity.CreatedAt.Before(from)) &&
			(to.IsZero() || !entity.CreatedAt.After(to)) {
			results = append(results, entity)
		}
	}
	return results, nil
}

func (r *ReportingBudgetRepository) Search(query string) ([]*ReportingBudgetEntity, error) {
	query = strings.ToLower(query)
	var results []*ReportingBudgetEntity
	for _, entity := range r.store {
		if strings.Contains(strings.ToLower(entity.Name), query) ||
			strings.Contains(strings.ToLower(entity.Description), query) ||
			strings.Contains(strings.ToLower(entity.Code), query) ||
			strings.Contains(strings.ToLower(entity.Reference), query) {
			results = append(results, entity)
		}
	}
	return results, nil
}

func (r *ReportingBudgetRepository) UpdateStatus(id int64, status string) error {
	entity, ok := r.store[id]
	if !ok {
		return fmt.Errorf("reporting_budget_repository: entity %d not found", id)
	}
	entity.Status = status
	entity.IsActive = status == "active"
	entity.UpdatedAt = time.Now()
	return nil
}
