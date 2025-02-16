package cache

import (
	"data-transformer-demo/internal/db"
	"log"
	"sync"
	"time"
)

type Cache struct {
	// Optimize memory layout with maps of maps
	mappingsByCountry map[string]map[int]*db.QuestionMapping
	cellsByMappingID  map[int]map[int]*db.QuestionMappingCell
	answersByCellID   map[int]map[int]*db.QuestionMappingCellAnswer
	mu                sync.RWMutex
}

var (
	instance *Cache
	once     sync.Once
)

func GetInstance() *Cache {
	once.Do(func() {
		instance = &Cache{
			mappingsByCountry: make(map[string]map[int]*db.QuestionMapping),
			cellsByMappingID:  make(map[int]map[int]*db.QuestionMappingCell),
			answersByCellID:   make(map[int]map[int]*db.QuestionMappingCellAnswer),
		}
	})
	return instance
}

func (c *Cache) LoadData() error {
	startTime := time.Now()

	// Create new maps to swap with old ones
	newMappingsByCountry := make(map[string]map[int]*db.QuestionMapping)
	newCellsByMappingID := make(map[int]map[int]*db.QuestionMappingCell)
	newAnswersByCellID := make(map[int]map[int]*db.QuestionMappingCellAnswer)

	// Load all mappings
	mappings, err := db.GetAllMappings()
	if err != nil {
		return err
	}

	// Pre-allocate maps with capacity
	for _, mapping := range mappings {
		country := mapping.Country
		if newMappingsByCountry[country] == nil {
			newMappingsByCountry[country] = make(map[int]*db.QuestionMapping)
		}
		m := mapping // Create a new variable to avoid pointing to loop variable
		newMappingsByCountry[country][mapping.ID] = &m

		// Pre-load cells
		cells, err := db.GetCellsByMappingID(mapping.ID)
		if err != nil {
			return err
		}

		newCellsByMappingID[mapping.ID] = make(map[int]*db.QuestionMappingCell, len(cells))
		for _, cell := range cells {
			c := cell // Create a new variable
			newCellsByMappingID[mapping.ID][cell.ID] = &c

			// Pre-load answers
			answers, err := db.GetAnswersByCellID(cell.ID)
			if err != nil {
				return err
			}

			newAnswersByCellID[cell.ID] = make(map[int]*db.QuestionMappingCellAnswer, len(answers))
			for _, answer := range answers {
				a := answer // Create a new variable
				newAnswersByCellID[cell.ID][answer.ID] = &a
			}
		}
	}

	// Atomic swap
	c.mu.Lock()
	c.mappingsByCountry = newMappingsByCountry
	c.cellsByMappingID = newCellsByMappingID
	c.answersByCellID = newAnswersByCellID
	c.mu.Unlock()

	log.Printf("Cache refresh completed in %v", time.Since(startTime))
	return nil
}

func (c *Cache) GetMappings(country string) []*db.QuestionMapping {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if mappings, exists := c.mappingsByCountry[country]; exists {
		result := make([]*db.QuestionMapping, 0, len(mappings))
		for _, mapping := range mappings {
			result = append(result, mapping)
		}
		return result
	}
	return nil
}

func (c *Cache) GetCells(mappingID int) []*db.QuestionMappingCell {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if cells, exists := c.cellsByMappingID[mappingID]; exists {
		result := make([]*db.QuestionMappingCell, 0, len(cells))
		for _, cell := range cells {
			result = append(result, cell)
		}
		return result
	}
	return nil
}

func (c *Cache) GetAnswers(cellID int) []*db.QuestionMappingCellAnswer {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if answers, exists := c.answersByCellID[cellID]; exists {
		result := make([]*db.QuestionMappingCellAnswer, 0, len(answers))
		for _, answer := range answers {
			result = append(result, answer)
		}
		return result
	}
	return nil
}

func (c *Cache) StartRefreshLoop(interval time.Duration) {
	go func() {
		ticker := time.NewTicker(interval)
		for range ticker.C {
			if err := c.LoadData(); err != nil {
				log.Printf("Error refreshing cache: %v", err)
			}
		}
	}()
}
