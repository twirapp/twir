package youtube

import (
	"net/http"
)

func (c *YouTube) HandleRequest(w http.ResponseWriter, r *http.Request) {
	c.manager.HandleRequest(w, r)
}
