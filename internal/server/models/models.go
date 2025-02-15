package models

type Metrics struct {
	ID    string   `json:"id" db:"name"`               // имя метрики
	MType string   `json:"type" db:"type"`             // параметр, принимающий значение gauge или counter
	Delta *int64   `json:"delta,omitempty" db:"delta"` // значение метрики в случае передачи counter
	Value *float64 `json:"value,omitempty" db:"value"` // значение метрики в случае передачи gauge
}
