package model

import (
	"fmt"
	"html"
	"strings"

	"github.com/gofrs/uuid/v5"
)

// Meal
type Meal struct {
	Base
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Cost        float32   `json:"cost"`
	CategoryID  uuid.UUID `json:"category_id"`
	Category    Category
}

func (t *Meal) ResourceName() string {
	return "meal"
}

func (t *Meal) Preloads() []string {
	return []string{"Category"}
}

// Validate checks structure consistency
func (t *Meal) Validate() error {
	if t.Name == "" {
		return fmt.Errorf("required Name")
	}
	if t.Description == "" {
		return fmt.Errorf("required Description")
	}
	if t.Cost == 0 {
		return fmt.Errorf("required Cost")
	}

	return nil
}

func (t *Meal) Prepare() error {
	err := t.BasePrepare()
	if err != nil {
		return err
	}

	t.Name = html.EscapeString(strings.TrimSpace(t.Name))
	t.Description = html.EscapeString(strings.TrimSpace(t.Description))

	return nil
}
