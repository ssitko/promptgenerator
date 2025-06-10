package database

import "gorm.io/gorm"

type Prompt struct {
	ID            uint   `gorm:"primaryKey"`
	Prompt        string `gorm:"not null"`
	Content       string `gorm:"not null"`
	Actor         string `gorm:"not null"`
	Comments      bool   `gorm:"not null"`
	Documentation bool   `gorm:"not null"`
	Explanations  bool   `gorm:"not null"`
}

// PromptRepository handles database operations
type PromptRepository struct {
	db *gorm.DB
}

// NewPromptRepository creates a new repository instance
func NewPromptRepository(db *gorm.DB) *PromptRepository {
	return &PromptRepository{db: db}
}

// CreatePrompt inserts a new prompt into the database
func (r *PromptRepository) CreatePrompt(prompt Prompt) error {
	return r.db.Create(&prompt).Error
}

// GetAllPrompts retrieves all prompts from the database
func (r *PromptRepository) GetAllPrompts() ([]Prompt, error) {
	var prompts []Prompt
	err := r.db.Find(&prompts).Error
	return prompts, err
}
