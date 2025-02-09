package user

import (
	"time"
)

const APP_USER_TABLE = "app_user"

type AppUser struct {
	ID          int64      `json:"id" db:"id"`
	Name        string     `json:"name" db:"name" validate:"required,max=255"`
	Email       string     `json:"email" db:"email" validate:"required,email,max=255"`
	PhoneNumber *string    `json:"phone_number,omitempty" db:"phone_number" validate:"omitempty,max=15"`
	Webhook     *string    `json:"webhook,omitempty" db:"webhook"`
	Active      bool       `json:"active" db:"active"`
	OptOutDate  *time.Time `json:"opt_out_date,omitempty" db:"opt_out_date"`
	CreatedAt   time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at" db:"updated_at"`
}

// func (u *AppUser) SetPhoneNumber(phone string) {
// 	if phone == "" {
// 		u.PhoneNumber = sql.NullString{}
// 		return
// 	}
// 	u.PhoneNumber = sql.NullString{
// 		String: phone,
// 		Valid:  true,
// 	}
// }

// func (u *AppUser) SetWebhook(webhook string) {
// 	if webhook == "" {
// 		u.Webhook = sql.NullString{}
// 		return
// 	}
// 	u.Webhook = sql.NullString{
// 		String: webhook,
// 		Valid:  true,
// 	}
// }

// func (u *AppUser) SetOptOutDate(date time.Time) {
// 	if date.IsZero() {
// 		u.OptOutDate = sql.NullTime{}
// 		return
// 	}
// 	u.OptOutDate = sql.NullTime{
// 		Time:  date,
// 		Valid: true,
// 	}
// }

// func (u *AppUser) GetPhoneNumber() string {
// 	if !u.PhoneNumber.Valid {
// 		return ""
// 	}
// 	return u.PhoneNumber.String
// }

// func (u *AppUser) GetWebhook() string {
// 	if !u.Webhook.Valid {
// 		return ""
// 	}
// 	return u.Webhook.String
// }

// func (u *AppUser) GetOptOutDate() time.Time {
// 	if !u.OptOutDate.Valid {
// 		return time.Time{}
// 	}
// 	return u.OptOutDate.Time
// }
