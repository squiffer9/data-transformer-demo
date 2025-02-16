package db

import (
	//"database/sql"
	"fmt"
	"log"
)

// Data structures
type QuestionMapping struct {
	ID      int
	Country string
}

type QuestionMappingCell struct {
	ID                int
	QuestionMappingID int
}

type QuestionMappingCellAnswer struct {
	ID                    int
	QuestionMappingCellID int
	QuestionID            int
	AnswerID              int
}

// GetAllMappings retrieves all question mappings
func GetAllMappings() ([]QuestionMapping, error) {
	query := `
		SELECT id, country 
		FROM question_mappings
	`
	rows, err := DB.Query(query)
	if err != nil {
		return nil, fmt.Errorf("query error in GetAllMappings: %v", err)
	}
	defer rows.Close()

	var mappings []QuestionMapping
	for rows.Next() {
		var qm QuestionMapping
		if err := rows.Scan(&qm.ID, &qm.Country); err != nil {
			return nil, fmt.Errorf("scan error in GetAllMappings: %v", err)
		}
		mappings = append(mappings, qm)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("iteration error in GetAllMappings: %v", err)
	}

	log.Printf("Retrieved %d mappings", len(mappings))
	return mappings, nil
}

// GetCellsByMappingID retrieves cells for a specific mapping
func GetCellsByMappingID(mappingID int) ([]QuestionMappingCell, error) {
	query := `
		SELECT id, question_mapping_id 
		FROM question_mapping_cells 
		WHERE question_mapping_id = ?
	`
	rows, err := DB.Query(query, mappingID)
	if err != nil {
		return nil, fmt.Errorf("query error in GetCellsByMappingID: %v", err)
	}
	defer rows.Close()

	var cells []QuestionMappingCell
	for rows.Next() {
		var cell QuestionMappingCell
		if err := rows.Scan(&cell.ID, &cell.QuestionMappingID); err != nil {
			return nil, fmt.Errorf("scan error in GetCellsByMappingID: %v", err)
		}
		cells = append(cells, cell)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("iteration error in GetCellsByMappingID: %v", err)
	}

	return cells, nil
}

// GetAnswersByCellID retrieves answers for a specific cell
func GetAnswersByCellID(cellID int) ([]QuestionMappingCellAnswer, error) {
	query := `
		SELECT id, question_mapping_cell_id, question_id, answer_id 
		FROM question_mapping_cell_answers 
		WHERE question_mapping_cell_id = ?
	`
	rows, err := DB.Query(query, cellID)
	if err != nil {
		return nil, fmt.Errorf("query error in GetAnswersByCellID: %v", err)
	}
	defer rows.Close()

	var answers []QuestionMappingCellAnswer
	for rows.Next() {
		var ans QuestionMappingCellAnswer
		if err := rows.Scan(&ans.ID, &ans.QuestionMappingCellID, &ans.QuestionID, &ans.AnswerID); err != nil {
			return nil, fmt.Errorf("scan error in GetAnswersByCellID: %v", err)
		}
		answers = append(answers, ans)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("iteration error in GetAnswersByCellID: %v", err)
	}

	return answers, nil
}

// GetQuestionMapping retrieves question mappings for a specific country
func GetQuestionMapping(country string) ([]QuestionMapping, error) {
	query := `
		SELECT id, country 
		FROM question_mappings 
		WHERE country = ?
	`
	rows, err := DB.Query(query, country)
	if err != nil {
		return nil, fmt.Errorf("query error in GetQuestionMapping: %v", err)
	}
	defer rows.Close()

	var mappings []QuestionMapping
	for rows.Next() {
		var qm QuestionMapping
		if err := rows.Scan(&qm.ID, &qm.Country); err != nil {
			return nil, fmt.Errorf("scan error in GetQuestionMapping: %v", err)
		}
		mappings = append(mappings, qm)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("iteration error in GetQuestionMapping: %v", err)
	}

	return mappings, nil
}
