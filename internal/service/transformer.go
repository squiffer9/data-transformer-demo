package service

import (
	"data-transformer-demo/internal/cache"
	"sync"
)

type TransformRequest struct {
	Country string    `json:"country"`
	Data    []QAEntry `json:"data"`
}

type QAEntry struct {
	QuestionID int `json:"question_id"`
	AnswerID   int `json:"answer_id"`
}

type TransformResponse struct {
	Data []QAEntry `json:"data"`
}

// Pool for frequently used maps
var qaMapPool = sync.Pool{
	New: func() interface{} {
		return make(map[string]struct{}, 1000)
	},
}

var answerCountPool = sync.Pool{
	New: func() interface{} {
		return make(map[int]int, 1000)
	},
}

func Transform(req TransformRequest) TransformResponse {
	c := cache.GetInstance()

	// Get a map from the pool
	existingQAMap := qaMapPool.Get().(map[string]struct{})
	answerCount := answerCountPool.Get().(map[int]int)
	defer func() {
		// Clear the map and return it to the pool
		for k := range existingQAMap {
			delete(existingQAMap, k)
		}
		for k := range answerCount {
			delete(answerCount, k)
		}
		qaMapPool.Put(existingQAMap)
		answerCountPool.Put(answerCount)
	}()

	// Copy the data to a new slice
	resultData := make([]QAEntry, 0, len(req.Data)*2)
	resultData = append(resultData, req.Data...)

	// Add existing QA pairs to the map
	for _, qa := range req.Data {
		key := string(qa.QuestionID) + ":" + string(qa.AnswerID)
		existingQAMap[key] = struct{}{}
	}

	// Get mappings for the country
	mappings := c.GetMappings(req.Country)
	if len(mappings) == 0 {
		return TransformResponse{Data: req.Data}
	}

	// Count the number of answers
	for _, mapping := range mappings {
		cells := c.GetCells(mapping.ID)
		for _, cell := range cells {
			answers := c.GetAnswers(cell.ID)
			for _, ans := range answers {
				answerCount[ans.AnswerID]++
			}
		}
	}

	// Add valid QA pairs to the result
	for _, mapping := range mappings {
		cells := c.GetCells(mapping.ID)
		for _, cell := range cells {
			answers := c.GetAnswers(cell.ID)

			// Count the number of questions
			questionCount := make(map[int]int, len(answers))
			for _, ans := range answers {
				questionCount[ans.QuestionID]++
			}

			// Check for duplicates
			hasDuplicates := false
			for _, count := range questionCount {
				if count > 1 {
					hasDuplicates = true
					break
				}
			}
			if hasDuplicates {
				continue
			}

			// Add valid QA pairs to the result
			for _, ans := range answers {
				if answerCount[ans.AnswerID] > 1 {
					continue
				}

				key := string(ans.QuestionID) + ":" + string(ans.AnswerID)
				if _, exists := existingQAMap[key]; !exists {
					existingQAMap[key] = struct{}{}
					resultData = append(resultData, QAEntry{
						QuestionID: ans.QuestionID,
						AnswerID:   ans.AnswerID,
					})
				}
			}
		}
	}

	return TransformResponse{Data: resultData}
}
