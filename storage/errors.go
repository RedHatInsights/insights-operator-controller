package storage

import (
	"fmt"
)

// ItemNotFoundError shows that item with provided ItemID wasn't found in storage
type ItemNotFoundError struct {
	ItemID interface{}
}

func (e *ItemNotFoundError) Error() string {
	return fmt.Sprintf("Item with ID %s was not found in the storage", e.ItemID)
}
