package model

import (
	"fmt"
	"html"
	"strings"
)

// Category
type Category struct {
	Base
	Name string `json:"name"`
}

func (t *Category) ResourceName() string {
	return "category"
}

func (t *Category) IsGlobal() bool {
	return true
}

// Validate checks structure consistency
func (t *Category) Validate() error {
	if t.Name == "" {
		return fmt.Errorf("required Name")
	}

	return nil
}

func (t *Category) Prepare() error {
	err := t.BasePrepare()
	if err != nil {
		return err
	}

	t.Name = html.EscapeString(strings.TrimSpace(t.Name))

	return nil
}
