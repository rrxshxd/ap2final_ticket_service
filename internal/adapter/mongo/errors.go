package mongo

import "fmt"

func mongoError(op string, err error) error {
	return fmt.Errorf("mongo: %s: %w", op, err)
}
