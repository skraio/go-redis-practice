package data

import "github.com/skraio/mini-godis/internal/validator"

type Record struct {
	Key   string `json:"key"`
	Value int32  `json:"value"`
}

func ValidateRecord(v *validator.Validator, record *Record) {
	v.Check(record.Key != "", "key", "must be provided")

	v.Check(len(record.Key) <= 28, "key", "must not be more than 28 bytes long")
	v.Check(record.Value != 0, "value", "must be provided")
}
