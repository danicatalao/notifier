package user

import (
	"time"
)

const APP_USER_TABLE = "app_user"

type AppUser struct {
	Id          int64      `json:"id" db:"id"`
	Name        string     `json:"name" db:"name" validate:"required,max=255,alpha"`
	Email       string     `json:"email" db:"email" validate:"required,email,max=255"`
	PhoneNumber *string    `json:"phone_number,omitempty" db:"phone_number" validate:"omitempty,max=15,phone"`
	Webhook     *string    `json:"webhook,omitempty" db:"webhook" validate:"omitempty,url"`
	Active      bool       `json:"active" db:"active" validate:"required"`
	OptOutDate  *time.Time `json:"opt_out_date,omitempty" db:"opt_out_date"`
	CreatedAt   time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at" db:"updated_at"`
}
