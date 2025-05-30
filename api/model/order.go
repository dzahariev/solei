package model

import (
	"fmt"

	"github.com/gofrs/uuid/v5"
)

// Order
type Order struct {
	Base
	Price  float32   `json:"price"`
	Status string    `json:"status"`
	UserID uuid.UUID `json:"user_id"`
	User   User
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
func (o *Order) Validate() error {
	if o.Price == 0 {
		return fmt.Errorf("required Price")
	}

	return nil
}

func (o *Order) Prepare() error {
	err := o.BasePrepare()
	if err != nil {
		return err
	}
	return nil
}
