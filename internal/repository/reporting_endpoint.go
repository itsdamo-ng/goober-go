package repository

import (
	"fmt"
	"strings"
	"time"
)

type ReportingEndpointEntity struct {
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

type ReportingEndpointQueryOptions struct {
	Where     []string
	Args      []interface{}
	OrderBy   string
	Direction string
	Limit     int
	Offset    int
}

type ReportingEndpointRepository struct {
	store  map[int64]*ReportingEndpointEntity
	nextID int64
}

func NewReportingEndpointRepository() *ReportingEndpointRepository {
	return &ReportingEndpointRepository{
		store:  make(map[int64]*ReportingEndpointEntity),
		nextID: 1,
	}
}

func (r *ReportingEndpointRepository) FindByID(id int64) (*ReportingEndpointEntity, error) {
	entity, ok := r.store[id]
	if !ok {
		return nil, fmt.Errorf("reporting_endpoint_repository: entity %d not found", id)
	}
	return entity, nil
}

func (r *ReportingEndpointRepository) FindAll(opts *ReportingEndpointQueryOptions) ([]*ReportingEndpointEntity, error) {
	var results []*ReportingEndpointEntity
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

func (r *ReportingEndpointRepository) FindByStatus(status string) ([]*ReportingEndpointEntity, error) {
	var results []*ReportingEndpointEntity
	for _, entity := range r.store {
		if entity.Status == status {
			results = append(results, entity)
		}
	}
	return results, nil
}

func (r *ReportingEndpointRepository) FindByCode(code string) (*ReportingEndpointEntity, error) {
	for _, entity := range r.store {
		if entity.Code == code {
			return entity, nil
		}
	}
	return nil, fmt.Errorf("reporting_endpoint_repository: entity with code %q not found", code)
}

func (r *ReportingEndpointRepository) FindByEmail(email string) (*ReportingEndpointEntity, error) {
	email = strings.ToLower(strings.TrimSpace(email))
	for _, entity := range r.store {
		if strings.ToLower(entity.Email) == email {
			return entity, nil
		}
	}
	return nil, fmt.Errorf("reporting_endpoint_repository: entity with email %q not found", email)
}

func (r *ReportingEndpointRepository) Save(entity *ReportingEndpointEntity) error {
	if entity.ID == 0 {
		entity.ID = r.nextID
		r.nextID++
		entity.CreatedAt = time.Now()
	}
	entity.UpdatedAt = time.Now()
	r.store[entity.ID] = entity
	return nil
}

func (r *ReportingEndpointRepository) Update(entity *ReportingEndpointEntity) error {
	if _, ok := r.store[entity.ID]; !ok {
		return fmt.Errorf("reporting_endpoint_repository: entity %d not found", entity.ID)
	}
	entity.UpdatedAt = time.Now()
	entity.Version++
	r.store[entity.ID] = entity
	return nil
}

func (r *ReportingEndpointRepository) Delete(id int64) error {
	if _, ok := r.store[id]; !ok {
		return fmt.Errorf("reporting_endpoint_repository: entity %d not found", id)
	}
	delete(r.store, id)
	return nil
}

func (r *ReportingEndpointRepository) SoftDelete(id int64) error {
	entity, ok := r.store[id]
	if !ok {
		return fmt.Errorf("reporting_endpoint_repository: entity %d not found", id)
	}
	entity.Status = "deleted"
	entity.IsActive = false
	entity.UpdatedAt = time.Now()
	return nil
}

func (r *ReportingEndpointRepository) Count() (int64, error) {
	return int64(len(r.store)), nil
}

func (r *ReportingEndpointRepository) CountByStatus(status string) (int64, error) {
	var count int64
	for _, entity := range r.store {
		if entity.Status == status {
			count++
		}
	}
	return count, nil
}

func (r *ReportingEndpointRepository) Exists(id int64) (bool, error) {
	_, ok := r.store[id]
	return ok, nil
}

func (r *ReportingEndpointRepository) ExistsByCode(code string) (bool, error) {
	for _, entity := range r.store {
		if entity.Code == code {
			return true, nil
		}
	}
	return false, nil
}

func (r *ReportingEndpointRepository) BulkInsert(entities []*ReportingEndpointEntity) error {
	for _, entity := range entities {
		if err := r.Save(entity); err != nil {
			return fmt.Errorf("reporting_endpoint_repository: bulk insert failed: %w", err)
		}
	}
	return nil
}

func (r *ReportingEndpointRepository) FindByDateRange(from, to time.Time) ([]*ReportingEndpointEntity, error) {
	var results []*ReportingEndpointEntity
	for _, entity := range r.store {
		if (from.IsZero() || !entity.CreatedAt.Before(from)) &&
			(to.IsZero() || !entity.CreatedAt.After(to)) {
			results = append(results, entity)
		}
	}
	return results, nil
}

func (r *ReportingEndpointRepository) Search(query string) ([]*ReportingEndpointEntity, error) {
	query = strings.ToLower(query)
	var results []*ReportingEndpointEntity
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

func (r *ReportingEndpointRepository) UpdateStatus(id int64, status string) error {
	entity, ok := r.store[id]
	if !ok {
		return fmt.Errorf("reporting_endpoint_repository: entity %d not found", id)
	}
	entity.Status = status
	entity.IsActive = status == "active"
	entity.UpdatedAt = time.Now()
	return nil
}
