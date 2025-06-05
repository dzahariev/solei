package model

import (
	"context"
	"fmt"

	"github.com/dzahariev/respite/basemodel"
	"github.com/gofrs/uuid/v5"
)

// Order
type Order struct {
	basemodel.Base
	Price  float32   `json:"price"`
	Status string    `json:"status"`
	UserID uuid.UUID `json:"user_id"`
	User   basemodel.User
}

func (o *Order) ResourceName() string {
	return "order"
}

func (o *Order) SetUserID(userID uuid.UUID) {
	o.UserID = userID
}

func (o *Order) Preloads() []string {
	return []string{"User"}
}

// Validate checks structure consistency
func (o *Order) Validate(ctx context.Context) error {
	if o.Price == 0 {
		return fmt.Errorf("required Price")
	}

	return nil
}

func (o *Order) Prepare(ctx context.Context) error {
	err := o.BasePrepare(ctx)
	if err != nil {
		return err
	}
	return nil
}
