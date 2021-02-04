package genose

import (
	core "github.com/genose-projects/genose-core"
)

func init() {
	/* Default Component Processors */
	core.Register(
		newControllerComponentProcessor,
	)
	/* Application Run Listeners */
	core.Register(NewEventPublishRunListener)
}
