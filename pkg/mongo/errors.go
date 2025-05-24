package mongo

import "fmt"

func connectionError(err error) error {
	return fmt.Errorf("mongo: connection error: %v", err)
}
