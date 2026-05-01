package repository

import (
	"fmt"
	"strings"
	"time"
)

type ShippingRewardEntity struct {
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

type ShippingRewardQueryOptions struct {
	Where     []string
	Args      []interface{}
	OrderBy   string
	Direction string
	Limit     int
	Offset    int
}

type ShippingRewardRepository struct {
	store  map[int64]*ShippingRewardEntity
	nextID int64
}

func NewShippingRewardRepository() *ShippingRewardRepository {
	return &ShippingRewardRepository{
		store:  make(map[int64]*ShippingRewardEntity),
		nextID: 1,
	}
}

func (r *ShippingRewardRepository) FindByID(id int64) (*ShippingRewardEntity, error) {
	entity, ok := r.store[id]
	if !ok {
		return nil, fmt.Errorf("shipping_reward_repository: entity %d not found", id)
	}
	return entity, nil
}

func (r *ShippingRewardRepository) FindAll(opts *ShippingRewardQueryOptions) ([]*ShippingRewardEntity, error) {
	var results []*ShippingRewardEntity
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

func (r *ShippingRewardRepository) FindByStatus(status string) ([]*ShippingRewardEntity, error) {
	var results []*ShippingRewardEntity
	for _, entity := range r.store {
		if entity.Status == status {
			results = append(results, entity)
		}
	}
	return results, nil
}

func (r *ShippingRewardRepository) FindByCode(code string) (*ShippingRewardEntity, error) {
	for _, entity := range r.store {
		if entity.Code == code {
			return entity, nil
		}
	}
	return nil, fmt.Errorf("shipping_reward_repository: entity with code %q not found", code)
}

func (r *ShippingRewardRepository) FindByEmail(email string) (*ShippingRewardEntity, error) {
	email = strings.ToLower(strings.TrimSpace(email))
	for _, entity := range r.store {
		if strings.ToLower(entity.Email) == email {
			return entity, nil
		}
	}
	return nil, fmt.Errorf("shipping_reward_repository: entity with email %q not found", email)
}

func (r *ShippingRewardRepository) Save(entity *ShippingRewardEntity) error {
	if entity.ID == 0 {
		entity.ID = r.nextID
		r.nextID++
		entity.CreatedAt = time.Now()
	}
	entity.UpdatedAt = time.Now()
	r.store[entity.ID] = entity
	return nil
}

func (r *ShippingRewardRepository) Update(entity *ShippingRewardEntity) error {
	if _, ok := r.store[entity.ID]; !ok {
		return fmt.Errorf("shipping_reward_repository: entity %d not found", entity.ID)
	}
	entity.UpdatedAt = time.Now()
	entity.Version++
	r.store[entity.ID] = entity
	return nil
}

func (r *ShippingRewardRepository) Delete(id int64) error {
	if _, ok := r.store[id]; !ok {
		return fmt.Errorf("shipping_reward_repository: entity %d not found", id)
	}
	delete(r.store, id)
	return nil
}

func (r *ShippingRewardRepository) SoftDelete(id int64) error {
	entity, ok := r.store[id]
	if !ok {
		return fmt.Errorf("shipping_reward_repository: entity %d not found", id)
	}
	entity.Status = "deleted"
	entity.IsActive = false
	entity.UpdatedAt = time.Now()
	return nil
}

func (r *ShippingRewardRepository) Count() (int64, error) {
	return int64(len(r.store)), nil
}

func (r *ShippingRewardRepository) CountByStatus(status string) (int64, error) {
	var count int64
	for _, entity := range r.store {
		if entity.Status == status {
			count++
		}
	}
	return count, nil
}

func (r *ShippingRewardRepository) Exists(id int64) (bool, error) {
	_, ok := r.store[id]
	return ok, nil
}

func (r *ShippingRewardRepository) ExistsByCode(code string) (bool, error) {
	for _, entity := range r.store {
		if entity.Code == code {
			return true, nil
		}
	}
	return false, nil
}

func (r *ShippingRewardRepository) BulkInsert(entities []*ShippingRewardEntity) error {
	for _, entity := range entities {
		if err := r.Save(entity); err != nil {
			return fmt.Errorf("shipping_reward_repository: bulk insert failed: %w", err)
		}
	}
	return nil
}

func (r *ShippingRewardRepository) FindByDateRange(from, to time.Time) ([]*ShippingRewardEntity, error) {
	var results []*ShippingRewardEntity
	for _, entity := range r.store {
		if (from.IsZero() || !entity.CreatedAt.Before(from)) &&
			(to.IsZero() || !entity.CreatedAt.After(to)) {
			results = append(results, entity)
		}
	}
	return results, nil
}

func (r *ShippingRewardRepository) Search(query string) ([]*ShippingRewardEntity, error) {
	query = strings.ToLower(query)
	var results []*ShippingRewardEntity
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

func (r *ShippingRewardRepository) UpdateStatus(id int64, status string) error {
	entity, ok := r.store[id]
	if !ok {
		return fmt.Errorf("shipping_reward_repository: entity %d not found", id)
	}
	entity.Status = status
	entity.IsActive = status == "active"
	entity.UpdatedAt = time.Now()
	return nil
}
