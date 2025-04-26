package halper

import (
	"fmt"
	"net/http"
)

func CheckBodyContent(r *http.Request) error {
	contentType := r.Header.Get("Content-Type")
	if contentType != "application/json" {
		return fmt.Errorf("Content-Type must be application/json")
	}
	return nil
}
