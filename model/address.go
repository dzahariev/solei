package model

import (
	"context"
	"fmt"
	"html"
	"strings"

	"github.com/dzahariev/respite/basemodel"
	"github.com/gofrs/uuid/v5"
)

// Address
type Address struct {
	basemodel.Base
	Country string    `json:"country"`
	City    string    `json:"city"`
	Street  string    `json:"street"`
	Phone   string    `json:"phone"`
	UserID  uuid.UUID `json:"user_id"`
	User    basemodel.User
}

func (t *Address) ResourceName() string {
	return "address"
}

func (t *Address) SetUserID(userID uuid.UUID) {
	t.UserID = userID
}

func (t *Address) Preloads() []string {
	return []string{"User"}
}

// Validate checks structure consistency
func (t *Address) Validate(ctx context.Context) error {
	if t.Country == "" {
		return fmt.Errorf("required Country")
	}
	if t.City == "" {
		return fmt.Errorf("required City")
	}
	if t.Street == "" {
		return fmt.Errorf("required Street")
	}
	if t.Phone == "" {
		return fmt.Errorf("required Phone")
	}

	return nil
}

func (t *Address) Prepare(ctx context.Context) error {
	err := t.BasePrepare(ctx)
	if err != nil {
		return err
	}

	t.Country = html.EscapeString(strings.TrimSpace(t.Country))
	t.City = html.EscapeString(strings.TrimSpace(t.City))
	t.Street = html.EscapeString(strings.TrimSpace(t.Street))
	t.Phone = html.EscapeString(strings.TrimSpace(t.Phone))

	return nil
}
