package repository

import (
	"fmt"
	"strings"
	"time"
)

type StreamingTransactionEntity struct {
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

type StreamingTransactionQueryOptions struct {
	Where     []string
	Args      []interface{}
	OrderBy   string
	Direction string
	Limit     int
	Offset    int
}

type StreamingTransactionRepository struct {
	store  map[int64]*StreamingTransactionEntity
	nextID int64
}

func NewStreamingTransactionRepository() *StreamingTransactionRepository {
	return &StreamingTransactionRepository{
		store:  make(map[int64]*StreamingTransactionEntity),
		nextID: 1,
	}
}

func (r *StreamingTransactionRepository) FindByID(id int64) (*StreamingTransactionEntity, error) {
	entity, ok := r.store[id]
	if !ok {
		return nil, fmt.Errorf("streaming_transaction_repository: entity %d not found", id)
	}
	return entity, nil
}

func (r *StreamingTransactionRepository) FindAll(opts *StreamingTransactionQueryOptions) ([]*StreamingTransactionEntity, error) {
	var results []*StreamingTransactionEntity
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

func (r *StreamingTransactionRepository) FindByStatus(status string) ([]*StreamingTransactionEntity, error) {
	var results []*StreamingTransactionEntity
	for _, entity := range r.store {
		if entity.Status == status {
			results = append(results, entity)
		}
	}
	return results, nil
}

func (r *StreamingTransactionRepository) FindByCode(code string) (*StreamingTransactionEntity, error) {
	for _, entity := range r.store {
		if entity.Code == code {
			return entity, nil
		}
	}
	return nil, fmt.Errorf("streaming_transaction_repository: entity with code %q not found", code)
}

func (r *StreamingTransactionRepository) FindByEmail(email string) (*StreamingTransactionEntity, error) {
	email = strings.ToLower(strings.TrimSpace(email))
	for _, entity := range r.store {
		if strings.ToLower(entity.Email) == email {
			return entity, nil
		}
	}
	return nil, fmt.Errorf("streaming_transaction_repository: entity with email %q not found", email)
}

func (r *StreamingTransactionRepository) Save(entity *StreamingTransactionEntity) error {
	if entity.ID == 0 {
		entity.ID = r.nextID
		r.nextID++
		entity.CreatedAt = time.Now()
	}
	entity.UpdatedAt = time.Now()
	r.store[entity.ID] = entity
	return nil
}

func (r *StreamingTransactionRepository) Update(entity *StreamingTransactionEntity) error {
	if _, ok := r.store[entity.ID]; !ok {
		return fmt.Errorf("streaming_transaction_repository: entity %d not found", entity.ID)
	}
	entity.UpdatedAt = time.Now()
	entity.Version++
	r.store[entity.ID] = entity
	return nil
}

func (r *StreamingTransactionRepository) Delete(id int64) error {
	if _, ok := r.store[id]; !ok {
		return fmt.Errorf("streaming_transaction_repository: entity %d not found", id)
	}
	delete(r.store, id)
	return nil
}

func (r *StreamingTransactionRepository) SoftDelete(id int64) error {
	entity, ok := r.store[id]
	if !ok {
		return fmt.Errorf("streaming_transaction_repository: entity %d not found", id)
	}
	entity.Status = "deleted"
	entity.IsActive = false
	entity.UpdatedAt = time.Now()
	return nil
}

func (r *StreamingTransactionRepository) Count() (int64, error) {
	return int64(len(r.store)), nil
}

func (r *StreamingTransactionRepository) CountByStatus(status string) (int64, error) {
	var count int64
	for _, entity := range r.store {
		if entity.Status == status {
			count++
		}
	}
	return count, nil
}

func (r *StreamingTransactionRepository) Exists(id int64) (bool, error) {
	_, ok := r.store[id]
	return ok, nil
}

func (r *StreamingTransactionRepository) ExistsByCode(code string) (bool, error) {
	for _, entity := range r.store {
		if entity.Code == code {
			return true, nil
		}
	}
	return false, nil
}

func (r *StreamingTransactionRepository) BulkInsert(entities []*StreamingTransactionEntity) error {
	for _, entity := range entities {
		if err := r.Save(entity); err != nil {
			return fmt.Errorf("streaming_transaction_repository: bulk insert failed: %w", err)
		}
	}
	return nil
}

func (r *StreamingTransactionRepository) FindByDateRange(from, to time.Time) ([]*StreamingTransactionEntity, error) {
	var results []*StreamingTransactionEntity
	for _, entity := range r.store {
		if (from.IsZero() || !entity.CreatedAt.Before(from)) &&
			(to.IsZero() || !entity.CreatedAt.After(to)) {
			results = append(results, entity)
		}
	}
	return results, nil
}

func (r *StreamingTransactionRepository) Search(query string) ([]*StreamingTransactionEntity, error) {
	query = strings.ToLower(query)
	var results []*StreamingTransactionEntity
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

func (r *StreamingTransactionRepository) UpdateStatus(id int64, status string) error {
	entity, ok := r.store[id]
	if !ok {
		return fmt.Errorf("streaming_transaction_repository: entity %d not found", id)
	}
	entity.Status = status
	entity.IsActive = status == "active"
	entity.UpdatedAt = time.Now()
	return nil
}
