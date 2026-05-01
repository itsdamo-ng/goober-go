package repository

import (
	"fmt"
	"strings"
	"time"
)

type ShippingShipmentEntity struct {
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

type ShippingShipmentQueryOptions struct {
	Where     []string
	Args      []interface{}
	OrderBy   string
	Direction string
	Limit     int
	Offset    int
}

type ShippingShipmentRepository struct {
	store  map[int64]*ShippingShipmentEntity
	nextID int64
}

func NewShippingShipmentRepository() *ShippingShipmentRepository {
	return &ShippingShipmentRepository{
		store:  make(map[int64]*ShippingShipmentEntity),
		nextID: 1,
	}
}

func (r *ShippingShipmentRepository) FindByID(id int64) (*ShippingShipmentEntity, error) {
	entity, ok := r.store[id]
	if !ok {
		return nil, fmt.Errorf("shipping_shipment_repository: entity %d not found", id)
	}
	return entity, nil
}

func (r *ShippingShipmentRepository) FindAll(opts *ShippingShipmentQueryOptions) ([]*ShippingShipmentEntity, error) {
	var results []*ShippingShipmentEntity
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

func (r *ShippingShipmentRepository) FindByStatus(status string) ([]*ShippingShipmentEntity, error) {
	var results []*ShippingShipmentEntity
	for _, entity := range r.store {
		if entity.Status == status {
			results = append(results, entity)
		}
	}
	return results, nil
}

func (r *ShippingShipmentRepository) FindByCode(code string) (*ShippingShipmentEntity, error) {
	for _, entity := range r.store {
		if entity.Code == code {
			return entity, nil
		}
	}
	return nil, fmt.Errorf("shipping_shipment_repository: entity with code %q not found", code)
}

func (r *ShippingShipmentRepository) FindByEmail(email string) (*ShippingShipmentEntity, error) {
	email = strings.ToLower(strings.TrimSpace(email))
	for _, entity := range r.store {
		if strings.ToLower(entity.Email) == email {
			return entity, nil
		}
	}
	return nil, fmt.Errorf("shipping_shipment_repository: entity with email %q not found", email)
}

func (r *ShippingShipmentRepository) Save(entity *ShippingShipmentEntity) error {
	if entity.ID == 0 {
		entity.ID = r.nextID
		r.nextID++
		entity.CreatedAt = time.Now()
	}
	entity.UpdatedAt = time.Now()
	r.store[entity.ID] = entity
	return nil
}

func (r *ShippingShipmentRepository) Update(entity *ShippingShipmentEntity) error {
	if _, ok := r.store[entity.ID]; !ok {
		return fmt.Errorf("shipping_shipment_repository: entity %d not found", entity.ID)
	}
	entity.UpdatedAt = time.Now()
	entity.Version++
	r.store[entity.ID] = entity
	return nil
}

func (r *ShippingShipmentRepository) Delete(id int64) error {
	if _, ok := r.store[id]; !ok {
		return fmt.Errorf("shipping_shipment_repository: entity %d not found", id)
	}
	delete(r.store, id)
	return nil
}

func (r *ShippingShipmentRepository) SoftDelete(id int64) error {
	entity, ok := r.store[id]
	if !ok {
		return fmt.Errorf("shipping_shipment_repository: entity %d not found", id)
	}
	entity.Status = "deleted"
	entity.IsActive = false
	entity.UpdatedAt = time.Now()
	return nil
}

func (r *ShippingShipmentRepository) Count() (int64, error) {
	return int64(len(r.store)), nil
}

func (r *ShippingShipmentRepository) CountByStatus(status string) (int64, error) {
	var count int64
	for _, entity := range r.store {
		if entity.Status == status {
			count++
		}
	}
	return count, nil
}

func (r *ShippingShipmentRepository) Exists(id int64) (bool, error) {
	_, ok := r.store[id]
	return ok, nil
}

func (r *ShippingShipmentRepository) ExistsByCode(code string) (bool, error) {
	for _, entity := range r.store {
		if entity.Code == code {
			return true, nil
		}
	}
	return false, nil
}

func (r *ShippingShipmentRepository) BulkInsert(entities []*ShippingShipmentEntity) error {
	for _, entity := range entities {
		if err := r.Save(entity); err != nil {
			return fmt.Errorf("shipping_shipment_repository: bulk insert failed: %w", err)
		}
	}
	return nil
}

func (r *ShippingShipmentRepository) FindByDateRange(from, to time.Time) ([]*ShippingShipmentEntity, error) {
	var results []*ShippingShipmentEntity
	for _, entity := range r.store {
		if (from.IsZero() || !entity.CreatedAt.Before(from)) &&
			(to.IsZero() || !entity.CreatedAt.After(to)) {
			results = append(results, entity)
		}
	}
	return results, nil
}

func (r *ShippingShipmentRepository) Search(query string) ([]*ShippingShipmentEntity, error) {
	query = strings.ToLower(query)
	var results []*ShippingShipmentEntity
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

func (r *ShippingShipmentRepository) UpdateStatus(id int64, status string) error {
	entity, ok := r.store[id]
	if !ok {
		return fmt.Errorf("shipping_shipment_repository: entity %d not found", id)
	}
	entity.Status = status
	entity.IsActive = status == "active"
	entity.UpdatedAt = time.Now()
	return nil
}
