package address

import (
	"time"

	db "github.com/aria/app/db/sqlc"
)

/* -------- Requests -------- */

type createAddressRequest struct {
	ClientID int64    `json:"client_id" validate:"required"`
	Title    string   `json:"title"     validate:"required,min=2,max=255"`
	City     string   `json:"city"      validate:"required"`
	Street   *string  `json:"street,omitempty"`
	Phone    *string  `json:"phone,omitempty"`
	Zip      *string  `json:"zip,omitempty"`
	Lat      *float32 `json:"lat,omitempty"`
	Long     *float32 `json:"long,omitempty"`
}

type updateAddressRequest struct {
	Title  string   `json:"title"     validate:"required,min=2,max=255"`
	City   string   `json:"city"      validate:"required"`
	Street *string  `json:"street,omitempty"`
	Phone  *string  `json:"phone,omitempty"`
	Zip    *string  `json:"zip,omitempty"`
	Lat    *float32 `json:"lat,omitempty"`
	Long   *float32 `json:"long,omitempty"`
}

/* -------- Responses -------- */

type addressResponse struct {
	ID        int64     `json:"id"`
	ClientID  int64     `json:"client_id"`
	Title     string    `json:"title"`
	City      string    `json:"city"`
	Street    *string   `json:"street,omitempty"`
	Phone     *string   `json:"phone,omitempty"`
	Zip       *string   `json:"zip,omitempty"`
	Lat       *float32  `json:"lat,omitempty"`
	Long      *float32  `json:"long,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func newAddressResponse(a db.Address) addressResponse {
	resp := addressResponse{
		ID:        a.ID,
		ClientID:  a.ClientID,
		Title:     a.Title,
		City:      a.City,
		CreatedAt: a.CreatedAt,
		UpdatedAt: a.UpdatedAt,
	}
	if a.Street.Valid {
		resp.Street = &a.Street.String
	}
	if a.Phone.Valid {
		resp.Phone = &a.Phone.String
	}
	if a.Zip.Valid {
		resp.Zip = &a.Zip.String
	}
	if a.Lat.Valid {
		v := float32(a.Lat.Float32)
		resp.Lat = &v
	}
	if a.Long.Valid {
		v := float32(a.Long.Float32)
		resp.Long = &v
	}
	return resp
}
