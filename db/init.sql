USE transformer;

CREATE TABLE IF NOT EXISTS question_mappings (
    id INT AUTO_INCREMENT PRIMARY KEY,
    country VARCHAR(10) NOT NULL
);

CREATE TABLE IF NOT EXISTS question_mapping_cells (
    id INT AUTO_INCREMENT PRIMARY KEY,
    question_mapping_id INT NOT NULL,
    FOREIGN KEY (question_mapping_id) REFERENCES question_mappings(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS question_mapping_cell_answers (
    id INT AUTO_INCREMENT PRIMARY KEY,
    question_mapping_cell_id INT NOT NULL,
    question_id INT NOT NULL,
    answer_id INT NOT NULL,
    FOREIGN KEY (question_mapping_cell_id) REFERENCES question_mapping_cells(id) ON DELETE CASCADE
);
