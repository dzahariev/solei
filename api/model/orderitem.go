package model

import (
	"fmt"
	"html"
	"strings"

	"github.com/gofrs/uuid/v5"
)

// OrderItem
type OrderItem struct {
	Base
	Amount  int       `json:"amount"`
	Comment string    `json:"comment"`
	MealID  uuid.UUID `json:"meal_id"`
	Meal    Meal
	OrderID uuid.UUID `json:"order_id"`
	Order   Order
	UserID  uuid.UUID `json:"user_id"`
	User    User
}

func (o *OrderItem) ResourceName() string {
	return "orderitem"
}

func (o *OrderItem) SetUserID(userID uuid.UUID) {
	o.UserID = userID
}

func (o *OrderItem) Preloads() []string {
	return []string{"Meal", "Order", "User"}
}

// Validate checks structure consistency
func (o *OrderItem) Validate() error {
	if o.Amount == 0 {
		return fmt.Errorf("required Amount")
	}

	return nil
}

func (o *OrderItem) Prepare() error {
	err := o.BasePrepare()
	if err != nil {
		return err
	}

	o.Comment = html.EscapeString(strings.TrimSpace(o.Comment))
	return nil
}
