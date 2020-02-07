package storage

import (
	"fmt"
)

type ItemNotFoundError struct {
	ItemID interface{}
}

func (e *ItemNotFoundError) Error() string {
	return fmt.Sprintf("Item with ID %s was not found in the storage", e.ItemID)
}
