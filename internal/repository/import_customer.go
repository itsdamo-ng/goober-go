package repository

import (
	"fmt"
	"strings"
	"time"
)

type ImportCustomerEntity struct {
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

type ImportCustomerQueryOptions struct {
	Where     []string
	Args      []interface{}
	OrderBy   string
	Direction string
	Limit     int
	Offset    int
}

type ImportCustomerRepository struct {
	store  map[int64]*ImportCustomerEntity
	nextID int64
}

func NewImportCustomerRepository() *ImportCustomerRepository {
	return &ImportCustomerRepository{
		store:  make(map[int64]*ImportCustomerEntity),
		nextID: 1,
	}
}

func (r *ImportCustomerRepository) FindByID(id int64) (*ImportCustomerEntity, error) {
	entity, ok := r.store[id]
	if !ok {
		return nil, fmt.Errorf("import_customer_repository: entity %d not found", id)
	}
	return entity, nil
}

func (r *ImportCustomerRepository) FindAll(opts *ImportCustomerQueryOptions) ([]*ImportCustomerEntity, error) {
	var results []*ImportCustomerEntity
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

func (r *ImportCustomerRepository) FindByStatus(status string) ([]*ImportCustomerEntity, error) {
	var results []*ImportCustomerEntity
	for _, entity := range r.store {
		if entity.Status == status {
			results = append(results, entity)
		}
	}
	return results, nil
}

func (r *ImportCustomerRepository) FindByCode(code string) (*ImportCustomerEntity, error) {
	for _, entity := range r.store {
		if entity.Code == code {
			return entity, nil
		}
	}
	return nil, fmt.Errorf("import_customer_repository: entity with code %q not found", code)
}

func (r *ImportCustomerRepository) FindByEmail(email string) (*ImportCustomerEntity, error) {
	email = strings.ToLower(strings.TrimSpace(email))
	for _, entity := range r.store {
		if strings.ToLower(entity.Email) == email {
			return entity, nil
		}
	}
	return nil, fmt.Errorf("import_customer_repository: entity with email %q not found", email)
}

func (r *ImportCustomerRepository) Save(entity *ImportCustomerEntity) error {
	if entity.ID == 0 {
		entity.ID = r.nextID
		r.nextID++
		entity.CreatedAt = time.Now()
	}
	entity.UpdatedAt = time.Now()
	r.store[entity.ID] = entity
	return nil
}

func (r *ImportCustomerRepository) Update(entity *ImportCustomerEntity) error {
	if _, ok := r.store[entity.ID]; !ok {
		return fmt.Errorf("import_customer_repository: entity %d not found", entity.ID)
	}
	entity.UpdatedAt = time.Now()
	entity.Version++
	r.store[entity.ID] = entity
	return nil
}

func (r *ImportCustomerRepository) Delete(id int64) error {
	if _, ok := r.store[id]; !ok {
		return fmt.Errorf("import_customer_repository: entity %d not found", id)
	}
	delete(r.store, id)
	return nil
}

func (r *ImportCustomerRepository) SoftDelete(id int64) error {
	entity, ok := r.store[id]
	if !ok {
		return fmt.Errorf("import_customer_repository: entity %d not found", id)
	}
	entity.Status = "deleted"
	entity.IsActive = false
	entity.UpdatedAt = time.Now()
	return nil
}

func (r *ImportCustomerRepository) Count() (int64, error) {
	return int64(len(r.store)), nil
}

func (r *ImportCustomerRepository) CountByStatus(status string) (int64, error) {
	var count int64
	for _, entity := range r.store {
		if entity.Status == status {
			count++
		}
	}
	return count, nil
}

func (r *ImportCustomerRepository) Exists(id int64) (bool, error) {
	_, ok := r.store[id]
	return ok, nil
}

func (r *ImportCustomerRepository) ExistsByCode(code string) (bool, error) {
	for _, entity := range r.store {
		if entity.Code == code {
			return true, nil
		}
	}
	return false, nil
}

func (r *ImportCustomerRepository) BulkInsert(entities []*ImportCustomerEntity) error {
	for _, entity := range entities {
		if err := r.Save(entity); err != nil {
			return fmt.Errorf("import_customer_repository: bulk insert failed: %w", err)
		}
	}
	return nil
}

func (r *ImportCustomerRepository) FindByDateRange(from, to time.Time) ([]*ImportCustomerEntity, error) {
	var results []*ImportCustomerEntity
	for _, entity := range r.store {
		if (from.IsZero() || !entity.CreatedAt.Before(from)) &&
			(to.IsZero() || !entity.CreatedAt.After(to)) {
			results = append(results, entity)
		}
	}
	return results, nil
}

func (r *ImportCustomerRepository) Search(query string) ([]*ImportCustomerEntity, error) {
	query = strings.ToLower(query)
	var results []*ImportCustomerEntity
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

func (r *ImportCustomerRepository) UpdateStatus(id int64, status string) error {
	entity, ok := r.store[id]
	if !ok {
		return fmt.Errorf("import_customer_repository: entity %d not found", id)
	}
	entity.Status = status
	entity.IsActive = status == "active"
	entity.UpdatedAt = time.Now()
	return nil
}
