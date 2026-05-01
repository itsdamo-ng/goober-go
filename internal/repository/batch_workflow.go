package repository

import (
	"fmt"
	"strings"
	"time"
)

type BatchWorkflowEntity struct {
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

type BatchWorkflowQueryOptions struct {
	Where     []string
	Args      []interface{}
	OrderBy   string
	Direction string
	Limit     int
	Offset    int
}

type BatchWorkflowRepository struct {
	store  map[int64]*BatchWorkflowEntity
	nextID int64
}

func NewBatchWorkflowRepository() *BatchWorkflowRepository {
	return &BatchWorkflowRepository{
		store:  make(map[int64]*BatchWorkflowEntity),
		nextID: 1,
	}
}

func (r *BatchWorkflowRepository) FindByID(id int64) (*BatchWorkflowEntity, error) {
	entity, ok := r.store[id]
	if !ok {
		return nil, fmt.Errorf("batch_workflow_repository: entity %d not found", id)
	}
	return entity, nil
}

func (r *BatchWorkflowRepository) FindAll(opts *BatchWorkflowQueryOptions) ([]*BatchWorkflowEntity, error) {
	var results []*BatchWorkflowEntity
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

func (r *BatchWorkflowRepository) FindByStatus(status string) ([]*BatchWorkflowEntity, error) {
	var results []*BatchWorkflowEntity
	for _, entity := range r.store {
		if entity.Status == status {
			results = append(results, entity)
		}
	}
	return results, nil
}

func (r *BatchWorkflowRepository) FindByCode(code string) (*BatchWorkflowEntity, error) {
	for _, entity := range r.store {
		if entity.Code == code {
			return entity, nil
		}
	}
	return nil, fmt.Errorf("batch_workflow_repository: entity with code %q not found", code)
}

func (r *BatchWorkflowRepository) FindByEmail(email string) (*BatchWorkflowEntity, error) {
	email = strings.ToLower(strings.TrimSpace(email))
	for _, entity := range r.store {
		if strings.ToLower(entity.Email) == email {
			return entity, nil
		}
	}
	return nil, fmt.Errorf("batch_workflow_repository: entity with email %q not found", email)
}

func (r *BatchWorkflowRepository) Save(entity *BatchWorkflowEntity) error {
	if entity.ID == 0 {
		entity.ID = r.nextID
		r.nextID++
		entity.CreatedAt = time.Now()
	}
	entity.UpdatedAt = time.Now()
	r.store[entity.ID] = entity
	return nil
}

func (r *BatchWorkflowRepository) Update(entity *BatchWorkflowEntity) error {
	if _, ok := r.store[entity.ID]; !ok {
		return fmt.Errorf("batch_workflow_repository: entity %d not found", entity.ID)
	}
	entity.UpdatedAt = time.Now()
	entity.Version++
	r.store[entity.ID] = entity
	return nil
}

func (r *BatchWorkflowRepository) Delete(id int64) error {
	if _, ok := r.store[id]; !ok {
		return fmt.Errorf("batch_workflow_repository: entity %d not found", id)
	}
	delete(r.store, id)
	return nil
}

func (r *BatchWorkflowRepository) SoftDelete(id int64) error {
	entity, ok := r.store[id]
	if !ok {
		return fmt.Errorf("batch_workflow_repository: entity %d not found", id)
	}
	entity.Status = "deleted"
	entity.IsActive = false
	entity.UpdatedAt = time.Now()
	return nil
}

func (r *BatchWorkflowRepository) Count() (int64, error) {
	return int64(len(r.store)), nil
}

func (r *BatchWorkflowRepository) CountByStatus(status string) (int64, error) {
	var count int64
	for _, entity := range r.store {
		if entity.Status == status {
			count++
		}
	}
	return count, nil
}

func (r *BatchWorkflowRepository) Exists(id int64) (bool, error) {
	_, ok := r.store[id]
	return ok, nil
}

func (r *BatchWorkflowRepository) ExistsByCode(code string) (bool, error) {
	for _, entity := range r.store {
		if entity.Code == code {
			return true, nil
		}
	}
	return false, nil
}

func (r *BatchWorkflowRepository) BulkInsert(entities []*BatchWorkflowEntity) error {
	for _, entity := range entities {
		if err := r.Save(entity); err != nil {
			return fmt.Errorf("batch_workflow_repository: bulk insert failed: %w", err)
		}
	}
	return nil
}

func (r *BatchWorkflowRepository) FindByDateRange(from, to time.Time) ([]*BatchWorkflowEntity, error) {
	var results []*BatchWorkflowEntity
	for _, entity := range r.store {
		if (from.IsZero() || !entity.CreatedAt.Before(from)) &&
			(to.IsZero() || !entity.CreatedAt.After(to)) {
			results = append(results, entity)
		}
	}
	return results, nil
}

func (r *BatchWorkflowRepository) Search(query string) ([]*BatchWorkflowEntity, error) {
	query = strings.ToLower(query)
	var results []*BatchWorkflowEntity
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

func (r *BatchWorkflowRepository) UpdateStatus(id int64, status string) error {
	entity, ok := r.store[id]
	if !ok {
		return fmt.Errorf("batch_workflow_repository: entity %d not found", id)
	}
	entity.Status = status
	entity.IsActive = status == "active"
	entity.UpdatedAt = time.Now()
	return nil
}
