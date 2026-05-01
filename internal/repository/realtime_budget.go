package repository

import (
	"fmt"
	"strings"
	"time"
)

type RealtimeBudgetEntity struct {
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

type RealtimeBudgetQueryOptions struct {
	Where     []string
	Args      []interface{}
	OrderBy   string
	Direction string
	Limit     int
	Offset    int
}

type RealtimeBudgetRepository struct {
	store  map[int64]*RealtimeBudgetEntity
	nextID int64
}

func NewRealtimeBudgetRepository() *RealtimeBudgetRepository {
	return &RealtimeBudgetRepository{
		store:  make(map[int64]*RealtimeBudgetEntity),
		nextID: 1,
	}
}

func (r *RealtimeBudgetRepository) FindByID(id int64) (*RealtimeBudgetEntity, error) {
	entity, ok := r.store[id]
	if !ok {
		return nil, fmt.Errorf("realtime_budget_repository: entity %d not found", id)
	}
	return entity, nil
}

func (r *RealtimeBudgetRepository) FindAll(opts *RealtimeBudgetQueryOptions) ([]*RealtimeBudgetEntity, error) {
	var results []*RealtimeBudgetEntity
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

func (r *RealtimeBudgetRepository) FindByStatus(status string) ([]*RealtimeBudgetEntity, error) {
	var results []*RealtimeBudgetEntity
	for _, entity := range r.store {
		if entity.Status == status {
			results = append(results, entity)
		}
	}
	return results, nil
}

func (r *RealtimeBudgetRepository) FindByCode(code string) (*RealtimeBudgetEntity, error) {
	for _, entity := range r.store {
		if entity.Code == code {
			return entity, nil
		}
	}
	return nil, fmt.Errorf("realtime_budget_repository: entity with code %q not found", code)
}

func (r *RealtimeBudgetRepository) FindByEmail(email string) (*RealtimeBudgetEntity, error) {
	email = strings.ToLower(strings.TrimSpace(email))
	for _, entity := range r.store {
		if strings.ToLower(entity.Email) == email {
			return entity, nil
		}
	}
	return nil, fmt.Errorf("realtime_budget_repository: entity with email %q not found", email)
}

func (r *RealtimeBudgetRepository) Save(entity *RealtimeBudgetEntity) error {
	if entity.ID == 0 {
		entity.ID = r.nextID
		r.nextID++
		entity.CreatedAt = time.Now()
	}
	entity.UpdatedAt = time.Now()
	r.store[entity.ID] = entity
	return nil
}

func (r *RealtimeBudgetRepository) Update(entity *RealtimeBudgetEntity) error {
	if _, ok := r.store[entity.ID]; !ok {
		return fmt.Errorf("realtime_budget_repository: entity %d not found", entity.ID)
	}
	entity.UpdatedAt = time.Now()
	entity.Version++
	r.store[entity.ID] = entity
	return nil
}

func (r *RealtimeBudgetRepository) Delete(id int64) error {
	if _, ok := r.store[id]; !ok {
		return fmt.Errorf("realtime_budget_repository: entity %d not found", id)
	}
	delete(r.store, id)
	return nil
}

func (r *RealtimeBudgetRepository) SoftDelete(id int64) error {
	entity, ok := r.store[id]
	if !ok {
		return fmt.Errorf("realtime_budget_repository: entity %d not found", id)
	}
	entity.Status = "deleted"
	entity.IsActive = false
	entity.UpdatedAt = time.Now()
	return nil
}

func (r *RealtimeBudgetRepository) Count() (int64, error) {
	return int64(len(r.store)), nil
}

func (r *RealtimeBudgetRepository) CountByStatus(status string) (int64, error) {
	var count int64
	for _, entity := range r.store {
		if entity.Status == status {
			count++
		}
	}
	return count, nil
}

func (r *RealtimeBudgetRepository) Exists(id int64) (bool, error) {
	_, ok := r.store[id]
	return ok, nil
}

func (r *RealtimeBudgetRepository) ExistsByCode(code string) (bool, error) {
	for _, entity := range r.store {
		if entity.Code == code {
			return true, nil
		}
	}
	return false, nil
}

func (r *RealtimeBudgetRepository) BulkInsert(entities []*RealtimeBudgetEntity) error {
	for _, entity := range entities {
		if err := r.Save(entity); err != nil {
			return fmt.Errorf("realtime_budget_repository: bulk insert failed: %w", err)
		}
	}
	return nil
}

func (r *RealtimeBudgetRepository) FindByDateRange(from, to time.Time) ([]*RealtimeBudgetEntity, error) {
	var results []*RealtimeBudgetEntity
	for _, entity := range r.store {
		if (from.IsZero() || !entity.CreatedAt.Before(from)) &&
			(to.IsZero() || !entity.CreatedAt.After(to)) {
			results = append(results, entity)
		}
	}
	return results, nil
}

func (r *RealtimeBudgetRepository) Search(query string) ([]*RealtimeBudgetEntity, error) {
	query = strings.ToLower(query)
	var results []*RealtimeBudgetEntity
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

func (r *RealtimeBudgetRepository) UpdateStatus(id int64, status string) error {
	entity, ok := r.store[id]
	if !ok {
		return fmt.Errorf("realtime_budget_repository: entity %d not found", id)
	}
	entity.Status = status
	entity.IsActive = status == "active"
	entity.UpdatedAt = time.Now()
	return nil
}
