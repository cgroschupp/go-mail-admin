package internal

import (
	"net/http"

	"github.com/go-chi/render"
)

func (m *MailServerConfiguratorInterface) getFeatureToggles(w http.ResponseWriter, r *http.Request) {

	render.JSON(w, r, m.Config.Feature)
}

func getVersion(w http.ResponseWriter, r *http.Request) {

	res := struct {
		Version string `json:"version"`
	}{
		Version: version,
	}
	render.JSON(w, r, res)
}
