package model

import (
	"context"
	"fmt"
	"html"
	"strings"

	"github.com/dzahariev/respite/basemodel"
)

// Category
type Category struct {
	basemodel.Base
	Name string `json:"name"`
}

func (t *Category) ResourceName() string {
	return "category"
}

func (t *Category) IsGlobal() bool {
	return true
}

// Validate checks structure consistency
func (t *Category) Validate(ctx context.Context) error {
	if t.Name == "" {
		return fmt.Errorf("required Name")
	}

	return nil
}

func (t *Category) Prepare(ctx context.Context) error {
	err := t.BasePrepare(ctx)
	if err != nil {
		return err
	}

	t.Name = html.EscapeString(strings.TrimSpace(t.Name))

	return nil
}
