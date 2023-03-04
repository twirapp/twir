package obs

import (
	"net/http"
)

func (c *OBS) HandleRequest(w http.ResponseWriter, r *http.Request) {
	c.manager.HandleRequest(w, r)
}
