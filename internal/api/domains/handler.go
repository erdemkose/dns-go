package domains

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"

	"dns/internal/dns"
)

type Controller struct {
	DNS dns.Service
}

func (c *Controller) Show(rw http.ResponseWriter, r *http.Request) {
	domain := chi.URLParam(r, "domain")
	resolver := chi.URLParam(r, "resolver")

	records, err := c.DNS.Resolve(r.Context(), domain, resolver)
	if err != nil {
		rw.Header().Set("Content-Type", "application/json")
		_, _ = rw.Write([]byte(err.Error()))

		return
	}

	out, _ := json.Marshal(records)

	rw.Header().Set("Content-Type", "application/json")
	_, _ = rw.Write(out)
}
