package util

import (
	"strconv"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

func ParseIntDefault(v string, def, min, max int) int {
	i, err := strconv.Atoi(v)
	if err != nil || i < min || i > max {
		return def
	}
	return i
}

func Deref[T any](ptr *T) T {
	var zero T // the zero value of T
	if ptr == nil {
		return zero
	}
	return *ptr
}

func DerefOr[T any](ptr *T, fallback T) T {
	if ptr == nil {
		return fallback
	}
	return *ptr
}

func CoalesceOptionalStr(a pgtype.Text, b *string) pgtype.Text {
	if b != nil {
		return pgtype.Text{String: *b, Valid: true}
	}
	return a
}

func Coalesce[T any](a T, b *T) T {
	if b != nil {
		return *b
	}
	return a
}

func Includes[T comparable](slice []T, item T) bool {
	for _, v := range slice {
		if v == item {
			return true
		}
	}
	return false
}

func ParseTime(v string) time.Time {
	if v == "" {
		return time.Time{}
	}
	t, _ := time.Parse(time.RFC3339, v)
	return t
}

func NullableFloat32(t pgtype.Float4) *float32 {
	if t.Valid {
		return &t.Float32
	}
	return nil
}

func NullableStr(t pgtype.Text) *string {
	if t.Valid {
		return &t.String
	}
	return nil
}

func NullableInt64(t pgtype.Int8) *int64 {
	if t.Valid {
		return &t.Int64
	}
	return nil
}

func NullableInt32(t pgtype.Int4) *int32 {
	if t.Valid {
		return &t.Int32
	}
	return nil
}

func ParseInt64(val string) *int64 {
	if val == "" {
		return nil
	}
	res, err := strconv.ParseInt(val, 10, 64)
	if err != nil {
		return nil
	}
	return &res
}
