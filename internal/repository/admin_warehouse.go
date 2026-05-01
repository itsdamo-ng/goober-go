package repository

import (
	"fmt"
	"strings"
	"time"
)

type AdminWarehouseEntity struct {
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

type AdminWarehouseQueryOptions struct {
	Where     []string
	Args      []interface{}
	OrderBy   string
	Direction string
	Limit     int
	Offset    int
}

type AdminWarehouseRepository struct {
	store  map[int64]*AdminWarehouseEntity
	nextID int64
}

func NewAdminWarehouseRepository() *AdminWarehouseRepository {
	return &AdminWarehouseRepository{
		store:  make(map[int64]*AdminWarehouseEntity),
		nextID: 1,
	}
}

func (r *AdminWarehouseRepository) FindByID(id int64) (*AdminWarehouseEntity, error) {
	entity, ok := r.store[id]
	if !ok {
		return nil, fmt.Errorf("admin_warehouse_repository: entity %d not found", id)
	}
	return entity, nil
}

func (r *AdminWarehouseRepository) FindAll(opts *AdminWarehouseQueryOptions) ([]*AdminWarehouseEntity, error) {
	var results []*AdminWarehouseEntity
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

func (r *AdminWarehouseRepository) FindByStatus(status string) ([]*AdminWarehouseEntity, error) {
	var results []*AdminWarehouseEntity
	for _, entity := range r.store {
		if entity.Status == status {
			results = append(results, entity)
		}
	}
	return results, nil
}

func (r *AdminWarehouseRepository) FindByCode(code string) (*AdminWarehouseEntity, error) {
	for _, entity := range r.store {
		if entity.Code == code {
			return entity, nil
		}
	}
	return nil, fmt.Errorf("admin_warehouse_repository: entity with code %q not found", code)
}

func (r *AdminWarehouseRepository) FindByEmail(email string) (*AdminWarehouseEntity, error) {
	email = strings.ToLower(strings.TrimSpace(email))
	for _, entity := range r.store {
		if strings.ToLower(entity.Email) == email {
			return entity, nil
		}
	}
	return nil, fmt.Errorf("admin_warehouse_repository: entity with email %q not found", email)
}

func (r *AdminWarehouseRepository) Save(entity *AdminWarehouseEntity) error {
	if entity.ID == 0 {
		entity.ID = r.nextID
		r.nextID++
		entity.CreatedAt = time.Now()
	}
	entity.UpdatedAt = time.Now()
	r.store[entity.ID] = entity
	return nil
}

func (r *AdminWarehouseRepository) Update(entity *AdminWarehouseEntity) error {
	if _, ok := r.store[entity.ID]; !ok {
		return fmt.Errorf("admin_warehouse_repository: entity %d not found", entity.ID)
	}
	entity.UpdatedAt = time.Now()
	entity.Version++
	r.store[entity.ID] = entity
	return nil
}

func (r *AdminWarehouseRepository) Delete(id int64) error {
	if _, ok := r.store[id]; !ok {
		return fmt.Errorf("admin_warehouse_repository: entity %d not found", id)
	}
	delete(r.store, id)
	return nil
}

func (r *AdminWarehouseRepository) SoftDelete(id int64) error {
	entity, ok := r.store[id]
	if !ok {
		return fmt.Errorf("admin_warehouse_repository: entity %d not found", id)
	}
	entity.Status = "deleted"
	entity.IsActive = false
	entity.UpdatedAt = time.Now()
	return nil
}

func (r *AdminWarehouseRepository) Count() (int64, error) {
	return int64(len(r.store)), nil
}

func (r *AdminWarehouseRepository) CountByStatus(status string) (int64, error) {
	var count int64
	for _, entity := range r.store {
		if entity.Status == status {
			count++
		}
	}
	return count, nil
}

func (r *AdminWarehouseRepository) Exists(id int64) (bool, error) {
	_, ok := r.store[id]
	return ok, nil
}

func (r *AdminWarehouseRepository) ExistsByCode(code string) (bool, error) {
	for _, entity := range r.store {
		if entity.Code == code {
			return true, nil
		}
	}
	return false, nil
}

func (r *AdminWarehouseRepository) BulkInsert(entities []*AdminWarehouseEntity) error {
	for _, entity := range entities {
		if err := r.Save(entity); err != nil {
			return fmt.Errorf("admin_warehouse_repository: bulk insert failed: %w", err)
		}
	}
	return nil
}

func (r *AdminWarehouseRepository) FindByDateRange(from, to time.Time) ([]*AdminWarehouseEntity, error) {
	var results []*AdminWarehouseEntity
	for _, entity := range r.store {
		if (from.IsZero() || !entity.CreatedAt.Before(from)) &&
			(to.IsZero() || !entity.CreatedAt.After(to)) {
			results = append(results, entity)
		}
	}
	return results, nil
}

func (r *AdminWarehouseRepository) Search(query string) ([]*AdminWarehouseEntity, error) {
	query = strings.ToLower(query)
	var results []*AdminWarehouseEntity
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

func (r *AdminWarehouseRepository) UpdateStatus(id int64, status string) error {
	entity, ok := r.store[id]
	if !ok {
		return fmt.Errorf("admin_warehouse_repository: entity %d not found", id)
	}
	entity.Status = status
	entity.IsActive = status == "active"
	entity.UpdatedAt = time.Now()
	return nil
}
