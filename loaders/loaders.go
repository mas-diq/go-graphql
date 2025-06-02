package loaders

import (
	"context"
	"mas-diq/go-graphql/models"
	"sync"

	"gorm.io/gorm"
)

type UserLoader struct {
	db    *gorm.DB
	cache map[uint]*models.User
	mutex sync.Mutex
}

func NewUserLoader(db *gorm.DB) *UserLoader {
	return &UserLoader{
		db:    db,
		cache: make(map[uint]*models.User),
	}
}

func (l *UserLoader) Load(ctx context.Context, ids []uint) ([]*models.User, error) {
	l.mutex.Lock()
	defer l.mutex.Unlock()

	// Find uncached IDs
	var uncached []uint
	for _, id := range ids {
		if _, exists := l.cache[id]; !exists {
			uncached = append(uncached, id)
		}
	}

	// Fetch uncached users
	if len(uncached) > 0 {
		var users []*models.User
		if err := l.db.Find(&users, uncached).Error; err != nil {
			return nil, err
		}
		for _, user := range users {
			l.cache[user.ID] = user
		}
	}

	// Return results in requested order
	result := make([]*models.User, len(ids))
	for i, id := range ids {
		result[i] = l.cache[id]
	}
	return result, nil
}
