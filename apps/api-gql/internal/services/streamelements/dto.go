package streamelements

import (
	"github.com/twirapp/twir/libs/integrations/streamelements"
)

type Data struct {
	Commands []streamelements.Command
	Timers   []streamelements.Timer
}
