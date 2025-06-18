package openapiauth

import "net/http"

// Bind implements render.Binder.
func (l *LoginParameter) Bind(r *http.Request) error {
	return nil
}
