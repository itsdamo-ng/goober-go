package repository

import (
	"fmt"
	"strings"
	"time"
)

type AdminAccountEntity struct {
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

type AdminAccountQueryOptions struct {
	Where     []string
	Args      []interface{}
	OrderBy   string
	Direction string
	Limit     int
	Offset    int
}

type AdminAccountRepository struct {
	store  map[int64]*AdminAccountEntity
	nextID int64
}

func NewAdminAccountRepository() *AdminAccountRepository {
	return &AdminAccountRepository{
		store:  make(map[int64]*AdminAccountEntity),
		nextID: 1,
	}
}

func (r *AdminAccountRepository) FindByID(id int64) (*AdminAccountEntity, error) {
	entity, ok := r.store[id]
	if !ok {
		return nil, fmt.Errorf("admin_account_repository: entity %d not found", id)
	}
	return entity, nil
}

func (r *AdminAccountRepository) FindAll(opts *AdminAccountQueryOptions) ([]*AdminAccountEntity, error) {
	var results []*AdminAccountEntity
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

func (r *AdminAccountRepository) FindByStatus(status string) ([]*AdminAccountEntity, error) {
	var results []*AdminAccountEntity
	for _, entity := range r.store {
		if entity.Status == status {
			results = append(results, entity)
		}
	}
	return results, nil
}

func (r *AdminAccountRepository) FindByCode(code string) (*AdminAccountEntity, error) {
	for _, entity := range r.store {
		if entity.Code == code {
			return entity, nil
		}
	}
	return nil, fmt.Errorf("admin_account_repository: entity with code %q not found", code)
}

func (r *AdminAccountRepository) FindByEmail(email string) (*AdminAccountEntity, error) {
	email = strings.ToLower(strings.TrimSpace(email))
	for _, entity := range r.store {
		if strings.ToLower(entity.Email) == email {
			return entity, nil
		}
	}
	return nil, fmt.Errorf("admin_account_repository: entity with email %q not found", email)
}

func (r *AdminAccountRepository) Save(entity *AdminAccountEntity) error {
	if entity.ID == 0 {
		entity.ID = r.nextID
		r.nextID++
		entity.CreatedAt = time.Now()
	}
	entity.UpdatedAt = time.Now()
	r.store[entity.ID] = entity
	return nil
}

func (r *AdminAccountRepository) Update(entity *AdminAccountEntity) error {
	if _, ok := r.store[entity.ID]; !ok {
		return fmt.Errorf("admin_account_repository: entity %d not found", entity.ID)
	}
	entity.UpdatedAt = time.Now()
	entity.Version++
	r.store[entity.ID] = entity
	return nil
}

func (r *AdminAccountRepository) Delete(id int64) error {
	if _, ok := r.store[id]; !ok {
		return fmt.Errorf("admin_account_repository: entity %d not found", id)
	}
	delete(r.store, id)
	return nil
}

func (r *AdminAccountRepository) SoftDelete(id int64) error {
	entity, ok := r.store[id]
	if !ok {
		return fmt.Errorf("admin_account_repository: entity %d not found", id)
	}
	entity.Status = "deleted"
	entity.IsActive = false
	entity.UpdatedAt = time.Now()
	return nil
}

func (r *AdminAccountRepository) Count() (int64, error) {
	return int64(len(r.store)), nil
}

func (r *AdminAccountRepository) CountByStatus(status string) (int64, error) {
	var count int64
	for _, entity := range r.store {
		if entity.Status == status {
			count++
		}
	}
	return count, nil
}

func (r *AdminAccountRepository) Exists(id int64) (bool, error) {
	_, ok := r.store[id]
	return ok, nil
}

func (r *AdminAccountRepository) ExistsByCode(code string) (bool, error) {
	for _, entity := range r.store {
		if entity.Code == code {
			return true, nil
		}
	}
	return false, nil
}

func (r *AdminAccountRepository) BulkInsert(entities []*AdminAccountEntity) error {
	for _, entity := range entities {
		if err := r.Save(entity); err != nil {
			return fmt.Errorf("admin_account_repository: bulk insert failed: %w", err)
		}
	}
	return nil
}

func (r *AdminAccountRepository) FindByDateRange(from, to time.Time) ([]*AdminAccountEntity, error) {
	var results []*AdminAccountEntity
	for _, entity := range r.store {
		if (from.IsZero() || !entity.CreatedAt.Before(from)) &&
			(to.IsZero() || !entity.CreatedAt.After(to)) {
			results = append(results, entity)
		}
	}
	return results, nil
}

func (r *AdminAccountRepository) Search(query string) ([]*AdminAccountEntity, error) {
	query = strings.ToLower(query)
	var results []*AdminAccountEntity
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

func (r *AdminAccountRepository) UpdateStatus(id int64, status string) error {
	entity, ok := r.store[id]
	if !ok {
		return fmt.Errorf("admin_account_repository: entity %d not found", id)
	}
	entity.Status = status
	entity.IsActive = status == "active"
	entity.UpdatedAt = time.Now()
	return nil
}
