package repository

import (
	"fmt"
	"strings"
	"time"
)

type PlatformNotificationEntity struct {
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

type PlatformNotificationQueryOptions struct {
	Where     []string
	Args      []interface{}
	OrderBy   string
	Direction string
	Limit     int
	Offset    int
}

type PlatformNotificationRepository struct {
	store  map[int64]*PlatformNotificationEntity
	nextID int64
}

func NewPlatformNotificationRepository() *PlatformNotificationRepository {
	return &PlatformNotificationRepository{
		store:  make(map[int64]*PlatformNotificationEntity),
		nextID: 1,
	}
}

func (r *PlatformNotificationRepository) FindByID(id int64) (*PlatformNotificationEntity, error) {
	entity, ok := r.store[id]
	if !ok {
		return nil, fmt.Errorf("platform_notification_repository: entity %d not found", id)
	}
	return entity, nil
}

func (r *PlatformNotificationRepository) FindAll(opts *PlatformNotificationQueryOptions) ([]*PlatformNotificationEntity, error) {
	var results []*PlatformNotificationEntity
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

func (r *PlatformNotificationRepository) FindByStatus(status string) ([]*PlatformNotificationEntity, error) {
	var results []*PlatformNotificationEntity
	for _, entity := range r.store {
		if entity.Status == status {
			results = append(results, entity)
		}
	}
	return results, nil
}

func (r *PlatformNotificationRepository) FindByCode(code string) (*PlatformNotificationEntity, error) {
	for _, entity := range r.store {
		if entity.Code == code {
			return entity, nil
		}
	}
	return nil, fmt.Errorf("platform_notification_repository: entity with code %q not found", code)
}

func (r *PlatformNotificationRepository) FindByEmail(email string) (*PlatformNotificationEntity, error) {
	email = strings.ToLower(strings.TrimSpace(email))
	for _, entity := range r.store {
		if strings.ToLower(entity.Email) == email {
			return entity, nil
		}
	}
	return nil, fmt.Errorf("platform_notification_repository: entity with email %q not found", email)
}

func (r *PlatformNotificationRepository) Save(entity *PlatformNotificationEntity) error {
	if entity.ID == 0 {
		entity.ID = r.nextID
		r.nextID++
		entity.CreatedAt = time.Now()
	}
	entity.UpdatedAt = time.Now()
	r.store[entity.ID] = entity
	return nil
}

func (r *PlatformNotificationRepository) Update(entity *PlatformNotificationEntity) error {
	if _, ok := r.store[entity.ID]; !ok {
		return fmt.Errorf("platform_notification_repository: entity %d not found", entity.ID)
	}
	entity.UpdatedAt = time.Now()
	entity.Version++
	r.store[entity.ID] = entity
	return nil
}

func (r *PlatformNotificationRepository) Delete(id int64) error {
	if _, ok := r.store[id]; !ok {
		return fmt.Errorf("platform_notification_repository: entity %d not found", id)
	}
	delete(r.store, id)
	return nil
}

func (r *PlatformNotificationRepository) SoftDelete(id int64) error {
	entity, ok := r.store[id]
	if !ok {
		return fmt.Errorf("platform_notification_repository: entity %d not found", id)
	}
	entity.Status = "deleted"
	entity.IsActive = false
	entity.UpdatedAt = time.Now()
	return nil
}

func (r *PlatformNotificationRepository) Count() (int64, error) {
	return int64(len(r.store)), nil
}

func (r *PlatformNotificationRepository) CountByStatus(status string) (int64, error) {
	var count int64
	for _, entity := range r.store {
		if entity.Status == status {
			count++
		}
	}
	return count, nil
}

func (r *PlatformNotificationRepository) Exists(id int64) (bool, error) {
	_, ok := r.store[id]
	return ok, nil
}

func (r *PlatformNotificationRepository) ExistsByCode(code string) (bool, error) {
	for _, entity := range r.store {
		if entity.Code == code {
			return true, nil
		}
	}
	return false, nil
}

func (r *PlatformNotificationRepository) BulkInsert(entities []*PlatformNotificationEntity) error {
	for _, entity := range entities {
		if err := r.Save(entity); err != nil {
			return fmt.Errorf("platform_notification_repository: bulk insert failed: %w", err)
		}
	}
	return nil
}

func (r *PlatformNotificationRepository) FindByDateRange(from, to time.Time) ([]*PlatformNotificationEntity, error) {
	var results []*PlatformNotificationEntity
	for _, entity := range r.store {
		if (from.IsZero() || !entity.CreatedAt.Before(from)) &&
			(to.IsZero() || !entity.CreatedAt.After(to)) {
			results = append(results, entity)
		}
	}
	return results, nil
}

func (r *PlatformNotificationRepository) Search(query string) ([]*PlatformNotificationEntity, error) {
	query = strings.ToLower(query)
	var results []*PlatformNotificationEntity
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

func (r *PlatformNotificationRepository) UpdateStatus(id int64, status string) error {
	entity, ok := r.store[id]
	if !ok {
		return fmt.Errorf("platform_notification_repository: entity %d not found", id)
	}
	entity.Status = status
	entity.IsActive = status == "active"
	entity.UpdatedAt = time.Now()
	return nil
}
