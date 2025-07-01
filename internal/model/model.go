package model

import (
	"encoding/json"
	"errors"
	"time"
)

type Pack struct {
	ID      uint64    `db:"id"`
	Name    string    `db:"name"`
	Created time.Time `db:"created"`
}

type PackEvent struct {
	ID     uint64 `db:"id"`
	Type   string `db:"type"`
	Status bool   `db:"lock"`
	Entity *Pack  `db:"payload"`
}

// func (p *Pack) String() string {
// 	return fmt.Sprintf("%d %s %v", p.ID, p.Name, p.Created)
// }

func (p *Pack) Scan(src any) error {
	var pack Pack

	switch src := src.(type) {
	case []byte:
		if err := json.Unmarshal([]byte(src), &pack); err != nil {
			return err
		}
	default:
		return errors.New("unsupported type for Scan")
	}

	*p = pack
	return nil
}
