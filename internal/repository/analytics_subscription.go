package repository

import (
	"fmt"
	"strings"
	"time"
)

type AnalyticsSubscriptionEntity struct {
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

type AnalyticsSubscriptionQueryOptions struct {
	Where     []string
	Args      []interface{}
	OrderBy   string
	Direction string
	Limit     int
	Offset    int
}

type AnalyticsSubscriptionRepository struct {
	store  map[int64]*AnalyticsSubscriptionEntity
	nextID int64
}

func NewAnalyticsSubscriptionRepository() *AnalyticsSubscriptionRepository {
	return &AnalyticsSubscriptionRepository{
		store:  make(map[int64]*AnalyticsSubscriptionEntity),
		nextID: 1,
	}
}

func (r *AnalyticsSubscriptionRepository) FindByID(id int64) (*AnalyticsSubscriptionEntity, error) {
	entity, ok := r.store[id]
	if !ok {
		return nil, fmt.Errorf("analytics_subscription_repository: entity %d not found", id)
	}
	return entity, nil
}

func (r *AnalyticsSubscriptionRepository) FindAll(opts *AnalyticsSubscriptionQueryOptions) ([]*AnalyticsSubscriptionEntity, error) {
	var results []*AnalyticsSubscriptionEntity
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

func (r *AnalyticsSubscriptionRepository) FindByStatus(status string) ([]*AnalyticsSubscriptionEntity, error) {
	var results []*AnalyticsSubscriptionEntity
	for _, entity := range r.store {
		if entity.Status == status {
			results = append(results, entity)
		}
	}
	return results, nil
}

func (r *AnalyticsSubscriptionRepository) FindByCode(code string) (*AnalyticsSubscriptionEntity, error) {
	for _, entity := range r.store {
		if entity.Code == code {
			return entity, nil
		}
	}
	return nil, fmt.Errorf("analytics_subscription_repository: entity with code %q not found", code)
}

func (r *AnalyticsSubscriptionRepository) FindByEmail(email string) (*AnalyticsSubscriptionEntity, error) {
	email = strings.ToLower(strings.TrimSpace(email))
	for _, entity := range r.store {
		if strings.ToLower(entity.Email) == email {
			return entity, nil
		}
	}
	return nil, fmt.Errorf("analytics_subscription_repository: entity with email %q not found", email)
}

func (r *AnalyticsSubscriptionRepository) Save(entity *AnalyticsSubscriptionEntity) error {
	if entity.ID == 0 {
		entity.ID = r.nextID
		r.nextID++
		entity.CreatedAt = time.Now()
	}
	entity.UpdatedAt = time.Now()
	r.store[entity.ID] = entity
	return nil
}

func (r *AnalyticsSubscriptionRepository) Update(entity *AnalyticsSubscriptionEntity) error {
	if _, ok := r.store[entity.ID]; !ok {
		return fmt.Errorf("analytics_subscription_repository: entity %d not found", entity.ID)
	}
	entity.UpdatedAt = time.Now()
	entity.Version++
	r.store[entity.ID] = entity
	return nil
}

func (r *AnalyticsSubscriptionRepository) Delete(id int64) error {
	if _, ok := r.store[id]; !ok {
		return fmt.Errorf("analytics_subscription_repository: entity %d not found", id)
	}
	delete(r.store, id)
	return nil
}

func (r *AnalyticsSubscriptionRepository) SoftDelete(id int64) error {
	entity, ok := r.store[id]
	if !ok {
		return fmt.Errorf("analytics_subscription_repository: entity %d not found", id)
	}
	entity.Status = "deleted"
	entity.IsActive = false
	entity.UpdatedAt = time.Now()
	return nil
}

func (r *AnalyticsSubscriptionRepository) Count() (int64, error) {
	return int64(len(r.store)), nil
}

func (r *AnalyticsSubscriptionRepository) CountByStatus(status string) (int64, error) {
	var count int64
	for _, entity := range r.store {
		if entity.Status == status {
			count++
		}
	}
	return count, nil
}

func (r *AnalyticsSubscriptionRepository) Exists(id int64) (bool, error) {
	_, ok := r.store[id]
	return ok, nil
}

func (r *AnalyticsSubscriptionRepository) ExistsByCode(code string) (bool, error) {
	for _, entity := range r.store {
		if entity.Code == code {
			return true, nil
		}
	}
	return false, nil
}

func (r *AnalyticsSubscriptionRepository) BulkInsert(entities []*AnalyticsSubscriptionEntity) error {
	for _, entity := range entities {
		if err := r.Save(entity); err != nil {
			return fmt.Errorf("analytics_subscription_repository: bulk insert failed: %w", err)
		}
	}
	return nil
}

func (r *AnalyticsSubscriptionRepository) FindByDateRange(from, to time.Time) ([]*AnalyticsSubscriptionEntity, error) {
	var results []*AnalyticsSubscriptionEntity
	for _, entity := range r.store {
		if (from.IsZero() || !entity.CreatedAt.Before(from)) &&
			(to.IsZero() || !entity.CreatedAt.After(to)) {
			results = append(results, entity)
		}
	}
	return results, nil
}

func (r *AnalyticsSubscriptionRepository) Search(query string) ([]*AnalyticsSubscriptionEntity, error) {
	query = strings.ToLower(query)
	var results []*AnalyticsSubscriptionEntity
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

func (r *AnalyticsSubscriptionRepository) UpdateStatus(id int64, status string) error {
	entity, ok := r.store[id]
	if !ok {
		return fmt.Errorf("analytics_subscription_repository: entity %d not found", id)
	}
	entity.Status = status
	entity.IsActive = status == "active"
	entity.UpdatedAt = time.Now()
	return nil
}
