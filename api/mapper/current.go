package mapper

import (
	"github.com/google/uuid"
	"github.com/light-bull/lightbull/shows"
)

type CurrentShowAndVisualJSON struct {
	ShowId *uuid.UUID `json:"showId"`
	VisualId *uuid.UUID `json:"visualId"`
}

func MapCurrent(show *shows.Show, visual *shows.Visual) CurrentShowAndVisualJSON {
	data := CurrentShowAndVisualJSON{
		ShowId: nil,
		VisualId: nil,
	}

	if show != nil {
		showId := show.ID
		data.ShowId = &showId
	}
	if visual != nil {
		visualId := visual.ID
		data.VisualId = &visualId
	}

	return data
}
