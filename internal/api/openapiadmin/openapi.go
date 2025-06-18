package openapiadmin

import "net/http"

// Render implements render.Renderer.
func (e Error) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

// Bind implements render.Binder.
func (d Domain) Bind(r *http.Request) error {
	return nil
}

// Bind implements render.Binder.
func (a AliasList) Bind(r *http.Request) error {
	return nil
}

// Bind implements render.Binder.
func (a *Alias) Bind(r *http.Request) error {
	return nil
}

// Bind implements render.Binder.
func (t *TLSPolicy) Bind(r *http.Request) error {
	return nil
}

// Bind implements render.Binder.
func (a *AccountCreate) Bind(r *http.Request) error {
	return nil
}

// Bind implements render.Binder.
func (c *ChangePasswordRequest) Bind(r *http.Request) error {
	return nil
}
