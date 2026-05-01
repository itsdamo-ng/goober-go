package repository

import (
	"fmt"
	"strings"
	"time"
)

type AuditSegmentEntity struct {
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

type AuditSegmentQueryOptions struct {
	Where     []string
	Args      []interface{}
	OrderBy   string
	Direction string
	Limit     int
	Offset    int
}

type AuditSegmentRepository struct {
	store  map[int64]*AuditSegmentEntity
	nextID int64
}

func NewAuditSegmentRepository() *AuditSegmentRepository {
	return &AuditSegmentRepository{
		store:  make(map[int64]*AuditSegmentEntity),
		nextID: 1,
	}
}

func (r *AuditSegmentRepository) FindByID(id int64) (*AuditSegmentEntity, error) {
	entity, ok := r.store[id]
	if !ok {
		return nil, fmt.Errorf("audit_segment_repository: entity %d not found", id)
	}
	return entity, nil
}

func (r *AuditSegmentRepository) FindAll(opts *AuditSegmentQueryOptions) ([]*AuditSegmentEntity, error) {
	var results []*AuditSegmentEntity
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

func (r *AuditSegmentRepository) FindByStatus(status string) ([]*AuditSegmentEntity, error) {
	var results []*AuditSegmentEntity
	for _, entity := range r.store {
		if entity.Status == status {
			results = append(results, entity)
		}
	}
	return results, nil
}

func (r *AuditSegmentRepository) FindByCode(code string) (*AuditSegmentEntity, error) {
	for _, entity := range r.store {
		if entity.Code == code {
			return entity, nil
		}
	}
	return nil, fmt.Errorf("audit_segment_repository: entity with code %q not found", code)
}

func (r *AuditSegmentRepository) FindByEmail(email string) (*AuditSegmentEntity, error) {
	email = strings.ToLower(strings.TrimSpace(email))
	for _, entity := range r.store {
		if strings.ToLower(entity.Email) == email {
			return entity, nil
		}
	}
	return nil, fmt.Errorf("audit_segment_repository: entity with email %q not found", email)
}

func (r *AuditSegmentRepository) Save(entity *AuditSegmentEntity) error {
	if entity.ID == 0 {
		entity.ID = r.nextID
		r.nextID++
		entity.CreatedAt = time.Now()
	}
	entity.UpdatedAt = time.Now()
	r.store[entity.ID] = entity
	return nil
}

func (r *AuditSegmentRepository) Update(entity *AuditSegmentEntity) error {
	if _, ok := r.store[entity.ID]; !ok {
		return fmt.Errorf("audit_segment_repository: entity %d not found", entity.ID)
	}
	entity.UpdatedAt = time.Now()
	entity.Version++
	r.store[entity.ID] = entity
	return nil
}

func (r *AuditSegmentRepository) Delete(id int64) error {
	if _, ok := r.store[id]; !ok {
		return fmt.Errorf("audit_segment_repository: entity %d not found", id)
	}
	delete(r.store, id)
	return nil
}

func (r *AuditSegmentRepository) SoftDelete(id int64) error {
	entity, ok := r.store[id]
	if !ok {
		return fmt.Errorf("audit_segment_repository: entity %d not found", id)
	}
	entity.Status = "deleted"
	entity.IsActive = false
	entity.UpdatedAt = time.Now()
	return nil
}

func (r *AuditSegmentRepository) Count() (int64, error) {
	return int64(len(r.store)), nil
}

func (r *AuditSegmentRepository) CountByStatus(status string) (int64, error) {
	var count int64
	for _, entity := range r.store {
		if entity.Status == status {
			count++
		}
	}
	return count, nil
}

func (r *AuditSegmentRepository) Exists(id int64) (bool, error) {
	_, ok := r.store[id]
	return ok, nil
}

func (r *AuditSegmentRepository) ExistsByCode(code string) (bool, error) {
	for _, entity := range r.store {
		if entity.Code == code {
			return true, nil
		}
	}
	return false, nil
}

func (r *AuditSegmentRepository) BulkInsert(entities []*AuditSegmentEntity) error {
	for _, entity := range entities {
		if err := r.Save(entity); err != nil {
			return fmt.Errorf("audit_segment_repository: bulk insert failed: %w", err)
		}
	}
	return nil
}

func (r *AuditSegmentRepository) FindByDateRange(from, to time.Time) ([]*AuditSegmentEntity, error) {
	var results []*AuditSegmentEntity
	for _, entity := range r.store {
		if (from.IsZero() || !entity.CreatedAt.Before(from)) &&
			(to.IsZero() || !entity.CreatedAt.After(to)) {
			results = append(results, entity)
		}
	}
	return results, nil
}

func (r *AuditSegmentRepository) Search(query string) ([]*AuditSegmentEntity, error) {
	query = strings.ToLower(query)
	var results []*AuditSegmentEntity
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

func (r *AuditSegmentRepository) UpdateStatus(id int64, status string) error {
	entity, ok := r.store[id]
	if !ok {
		return fmt.Errorf("audit_segment_repository: entity %d not found", id)
	}
	entity.Status = status
	entity.IsActive = status == "active"
	entity.UpdatedAt = time.Now()
	return nil
}
