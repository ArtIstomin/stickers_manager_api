package grifts

import (
	"github.com/artisomin/stickers_manager_api/actions"
	"github.com/gobuffalo/buffalo"
)

func init() {
	buffalo.Grifts(actions.App())
}
