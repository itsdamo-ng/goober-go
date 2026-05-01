package repository

import (
	"fmt"
	"strings"
	"time"
)

type InventoryEventEntity struct {
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

type InventoryEventQueryOptions struct {
	Where     []string
	Args      []interface{}
	OrderBy   string
	Direction string
	Limit     int
	Offset    int
}

type InventoryEventRepository struct {
	store  map[int64]*InventoryEventEntity
	nextID int64
}

func NewInventoryEventRepository() *InventoryEventRepository {
	return &InventoryEventRepository{
		store:  make(map[int64]*InventoryEventEntity),
		nextID: 1,
	}
}

func (r *InventoryEventRepository) FindByID(id int64) (*InventoryEventEntity, error) {
	entity, ok := r.store[id]
	if !ok {
		return nil, fmt.Errorf("inventory_event_repository: entity %d not found", id)
	}
	return entity, nil
}

func (r *InventoryEventRepository) FindAll(opts *InventoryEventQueryOptions) ([]*InventoryEventEntity, error) {
	var results []*InventoryEventEntity
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

func (r *InventoryEventRepository) FindByStatus(status string) ([]*InventoryEventEntity, error) {
	var results []*InventoryEventEntity
	for _, entity := range r.store {
		if entity.Status == status {
			results = append(results, entity)
		}
	}
	return results, nil
}

func (r *InventoryEventRepository) FindByCode(code string) (*InventoryEventEntity, error) {
	for _, entity := range r.store {
		if entity.Code == code {
			return entity, nil
		}
	}
	return nil, fmt.Errorf("inventory_event_repository: entity with code %q not found", code)
}

func (r *InventoryEventRepository) FindByEmail(email string) (*InventoryEventEntity, error) {
	email = strings.ToLower(strings.TrimSpace(email))
	for _, entity := range r.store {
		if strings.ToLower(entity.Email) == email {
			return entity, nil
		}
	}
	return nil, fmt.Errorf("inventory_event_repository: entity with email %q not found", email)
}

func (r *InventoryEventRepository) Save(entity *InventoryEventEntity) error {
	if entity.ID == 0 {
		entity.ID = r.nextID
		r.nextID++
		entity.CreatedAt = time.Now()
	}
	entity.UpdatedAt = time.Now()
	r.store[entity.ID] = entity
	return nil
}

func (r *InventoryEventRepository) Update(entity *InventoryEventEntity) error {
	if _, ok := r.store[entity.ID]; !ok {
		return fmt.Errorf("inventory_event_repository: entity %d not found", entity.ID)
	}
	entity.UpdatedAt = time.Now()
	entity.Version++
	r.store[entity.ID] = entity
	return nil
}

func (r *InventoryEventRepository) Delete(id int64) error {
	if _, ok := r.store[id]; !ok {
		return fmt.Errorf("inventory_event_repository: entity %d not found", id)
	}
	delete(r.store, id)
	return nil
}

func (r *InventoryEventRepository) SoftDelete(id int64) error {
	entity, ok := r.store[id]
	if !ok {
		return fmt.Errorf("inventory_event_repository: entity %d not found", id)
	}
	entity.Status = "deleted"
	entity.IsActive = false
	entity.UpdatedAt = time.Now()
	return nil
}

func (r *InventoryEventRepository) Count() (int64, error) {
	return int64(len(r.store)), nil
}

func (r *InventoryEventRepository) CountByStatus(status string) (int64, error) {
	var count int64
	for _, entity := range r.store {
		if entity.Status == status {
			count++
		}
	}
	return count, nil
}

func (r *InventoryEventRepository) Exists(id int64) (bool, error) {
	_, ok := r.store[id]
	return ok, nil
}

func (r *InventoryEventRepository) ExistsByCode(code string) (bool, error) {
	for _, entity := range r.store {
		if entity.Code == code {
			return true, nil
		}
	}
	return false, nil
}

func (r *InventoryEventRepository) BulkInsert(entities []*InventoryEventEntity) error {
	for _, entity := range entities {
		if err := r.Save(entity); err != nil {
			return fmt.Errorf("inventory_event_repository: bulk insert failed: %w", err)
		}
	}
	return nil
}

func (r *InventoryEventRepository) FindByDateRange(from, to time.Time) ([]*InventoryEventEntity, error) {
	var results []*InventoryEventEntity
	for _, entity := range r.store {
		if (from.IsZero() || !entity.CreatedAt.Before(from)) &&
			(to.IsZero() || !entity.CreatedAt.After(to)) {
			results = append(results, entity)
		}
	}
	return results, nil
}

func (r *InventoryEventRepository) Search(query string) ([]*InventoryEventEntity, error) {
	query = strings.ToLower(query)
	var results []*InventoryEventEntity
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

func (r *InventoryEventRepository) UpdateStatus(id int64, status string) error {
	entity, ok := r.store[id]
	if !ok {
		return fmt.Errorf("inventory_event_repository: entity %d not found", id)
	}
	entity.Status = status
	entity.IsActive = status == "active"
	entity.UpdatedAt = time.Now()
	return nil
}
