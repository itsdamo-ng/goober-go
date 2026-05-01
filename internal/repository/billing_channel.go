package repository

import (
	"fmt"
	"strings"
	"time"
)

type BillingChannelEntity struct {
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

type BillingChannelQueryOptions struct {
	Where     []string
	Args      []interface{}
	OrderBy   string
	Direction string
	Limit     int
	Offset    int
}

type BillingChannelRepository struct {
	store  map[int64]*BillingChannelEntity
	nextID int64
}

func NewBillingChannelRepository() *BillingChannelRepository {
	return &BillingChannelRepository{
		store:  make(map[int64]*BillingChannelEntity),
		nextID: 1,
	}
}

func (r *BillingChannelRepository) FindByID(id int64) (*BillingChannelEntity, error) {
	entity, ok := r.store[id]
	if !ok {
		return nil, fmt.Errorf("billing_channel_repository: entity %d not found", id)
	}
	return entity, nil
}

func (r *BillingChannelRepository) FindAll(opts *BillingChannelQueryOptions) ([]*BillingChannelEntity, error) {
	var results []*BillingChannelEntity
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

func (r *BillingChannelRepository) FindByStatus(status string) ([]*BillingChannelEntity, error) {
	var results []*BillingChannelEntity
	for _, entity := range r.store {
		if entity.Status == status {
			results = append(results, entity)
		}
	}
	return results, nil
}

func (r *BillingChannelRepository) FindByCode(code string) (*BillingChannelEntity, error) {
	for _, entity := range r.store {
		if entity.Code == code {
			return entity, nil
		}
	}
	return nil, fmt.Errorf("billing_channel_repository: entity with code %q not found", code)
}

func (r *BillingChannelRepository) FindByEmail(email string) (*BillingChannelEntity, error) {
	email = strings.ToLower(strings.TrimSpace(email))
	for _, entity := range r.store {
		if strings.ToLower(entity.Email) == email {
			return entity, nil
		}
	}
	return nil, fmt.Errorf("billing_channel_repository: entity with email %q not found", email)
}

func (r *BillingChannelRepository) Save(entity *BillingChannelEntity) error {
	if entity.ID == 0 {
		entity.ID = r.nextID
		r.nextID++
		entity.CreatedAt = time.Now()
	}
	entity.UpdatedAt = time.Now()
	r.store[entity.ID] = entity
	return nil
}

func (r *BillingChannelRepository) Update(entity *BillingChannelEntity) error {
	if _, ok := r.store[entity.ID]; !ok {
		return fmt.Errorf("billing_channel_repository: entity %d not found", entity.ID)
	}
	entity.UpdatedAt = time.Now()
	entity.Version++
	r.store[entity.ID] = entity
	return nil
}

func (r *BillingChannelRepository) Delete(id int64) error {
	if _, ok := r.store[id]; !ok {
		return fmt.Errorf("billing_channel_repository: entity %d not found", id)
	}
	delete(r.store, id)
	return nil
}

func (r *BillingChannelRepository) SoftDelete(id int64) error {
	entity, ok := r.store[id]
	if !ok {
		return fmt.Errorf("billing_channel_repository: entity %d not found", id)
	}
	entity.Status = "deleted"
	entity.IsActive = false
	entity.UpdatedAt = time.Now()
	return nil
}

func (r *BillingChannelRepository) Count() (int64, error) {
	return int64(len(r.store)), nil
}

func (r *BillingChannelRepository) CountByStatus(status string) (int64, error) {
	var count int64
	for _, entity := range r.store {
		if entity.Status == status {
			count++
		}
	}
	return count, nil
}

func (r *BillingChannelRepository) Exists(id int64) (bool, error) {
	_, ok := r.store[id]
	return ok, nil
}

func (r *BillingChannelRepository) ExistsByCode(code string) (bool, error) {
	for _, entity := range r.store {
		if entity.Code == code {
			return true, nil
		}
	}
	return false, nil
}

func (r *BillingChannelRepository) BulkInsert(entities []*BillingChannelEntity) error {
	for _, entity := range entities {
		if err := r.Save(entity); err != nil {
			return fmt.Errorf("billing_channel_repository: bulk insert failed: %w", err)
		}
	}
	return nil
}

func (r *BillingChannelRepository) FindByDateRange(from, to time.Time) ([]*BillingChannelEntity, error) {
	var results []*BillingChannelEntity
	for _, entity := range r.store {
		if (from.IsZero() || !entity.CreatedAt.Before(from)) &&
			(to.IsZero() || !entity.CreatedAt.After(to)) {
			results = append(results, entity)
		}
	}
	return results, nil
}

func (r *BillingChannelRepository) Search(query string) ([]*BillingChannelEntity, error) {
	query = strings.ToLower(query)
	var results []*BillingChannelEntity
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

func (r *BillingChannelRepository) UpdateStatus(id int64, status string) error {
	entity, ok := r.store[id]
	if !ok {
		return fmt.Errorf("billing_channel_repository: entity %d not found", id)
	}
	entity.Status = status
	entity.IsActive = status == "active"
	entity.UpdatedAt = time.Now()
	return nil
}
