package repository

import (
	"fmt"
	"strings"
	"time"
)

type SupportExperimentEntity struct {
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

type SupportExperimentQueryOptions struct {
	Where     []string
	Args      []interface{}
	OrderBy   string
	Direction string
	Limit     int
	Offset    int
}

type SupportExperimentRepository struct {
	store  map[int64]*SupportExperimentEntity
	nextID int64
}

func NewSupportExperimentRepository() *SupportExperimentRepository {
	return &SupportExperimentRepository{
		store:  make(map[int64]*SupportExperimentEntity),
		nextID: 1,
	}
}

func (r *SupportExperimentRepository) FindByID(id int64) (*SupportExperimentEntity, error) {
	entity, ok := r.store[id]
	if !ok {
		return nil, fmt.Errorf("support_experiment_repository: entity %d not found", id)
	}
	return entity, nil
}

func (r *SupportExperimentRepository) FindAll(opts *SupportExperimentQueryOptions) ([]*SupportExperimentEntity, error) {
	var results []*SupportExperimentEntity
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

func (r *SupportExperimentRepository) FindByStatus(status string) ([]*SupportExperimentEntity, error) {
	var results []*SupportExperimentEntity
	for _, entity := range r.store {
		if entity.Status == status {
			results = append(results, entity)
		}
	}
	return results, nil
}

func (r *SupportExperimentRepository) FindByCode(code string) (*SupportExperimentEntity, error) {
	for _, entity := range r.store {
		if entity.Code == code {
			return entity, nil
		}
	}
	return nil, fmt.Errorf("support_experiment_repository: entity with code %q not found", code)
}

func (r *SupportExperimentRepository) FindByEmail(email string) (*SupportExperimentEntity, error) {
	email = strings.ToLower(strings.TrimSpace(email))
	for _, entity := range r.store {
		if strings.ToLower(entity.Email) == email {
			return entity, nil
		}
	}
	return nil, fmt.Errorf("support_experiment_repository: entity with email %q not found", email)
}

func (r *SupportExperimentRepository) Save(entity *SupportExperimentEntity) error {
	if entity.ID == 0 {
		entity.ID = r.nextID
		r.nextID++
		entity.CreatedAt = time.Now()
	}
	entity.UpdatedAt = time.Now()
	r.store[entity.ID] = entity
	return nil
}

func (r *SupportExperimentRepository) Update(entity *SupportExperimentEntity) error {
	if _, ok := r.store[entity.ID]; !ok {
		return fmt.Errorf("support_experiment_repository: entity %d not found", entity.ID)
	}
	entity.UpdatedAt = time.Now()
	entity.Version++
	r.store[entity.ID] = entity
	return nil
}

func (r *SupportExperimentRepository) Delete(id int64) error {
	if _, ok := r.store[id]; !ok {
		return fmt.Errorf("support_experiment_repository: entity %d not found", id)
	}
	delete(r.store, id)
	return nil
}

func (r *SupportExperimentRepository) SoftDelete(id int64) error {
	entity, ok := r.store[id]
	if !ok {
		return fmt.Errorf("support_experiment_repository: entity %d not found", id)
	}
	entity.Status = "deleted"
	entity.IsActive = false
	entity.UpdatedAt = time.Now()
	return nil
}

func (r *SupportExperimentRepository) Count() (int64, error) {
	return int64(len(r.store)), nil
}

func (r *SupportExperimentRepository) CountByStatus(status string) (int64, error) {
	var count int64
	for _, entity := range r.store {
		if entity.Status == status {
			count++
		}
	}
	return count, nil
}

func (r *SupportExperimentRepository) Exists(id int64) (bool, error) {
	_, ok := r.store[id]
	return ok, nil
}

func (r *SupportExperimentRepository) ExistsByCode(code string) (bool, error) {
	for _, entity := range r.store {
		if entity.Code == code {
			return true, nil
		}
	}
	return false, nil
}

func (r *SupportExperimentRepository) BulkInsert(entities []*SupportExperimentEntity) error {
	for _, entity := range entities {
		if err := r.Save(entity); err != nil {
			return fmt.Errorf("support_experiment_repository: bulk insert failed: %w", err)
		}
	}
	return nil
}

func (r *SupportExperimentRepository) FindByDateRange(from, to time.Time) ([]*SupportExperimentEntity, error) {
	var results []*SupportExperimentEntity
	for _, entity := range r.store {
		if (from.IsZero() || !entity.CreatedAt.Before(from)) &&
			(to.IsZero() || !entity.CreatedAt.After(to)) {
			results = append(results, entity)
		}
	}
	return results, nil
}

func (r *SupportExperimentRepository) Search(query string) ([]*SupportExperimentEntity, error) {
	query = strings.ToLower(query)
	var results []*SupportExperimentEntity
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

func (r *SupportExperimentRepository) UpdateStatus(id int64, status string) error {
	entity, ok := r.store[id]
	if !ok {
		return fmt.Errorf("support_experiment_repository: entity %d not found", id)
	}
	entity.Status = status
	entity.IsActive = status == "active"
	entity.UpdatedAt = time.Now()
	return nil
}
