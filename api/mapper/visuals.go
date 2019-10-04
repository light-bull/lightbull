package mapper

import (
	"github.com/google/uuid"
	"github.com/light-bull/lightbull/shows"
)

type VisualJSON struct {
	ID     uuid.UUID `json:"id"`
	Name   string    `json:"name"`
	ShowId uuid.UUID `json:"showId"`
}

type VisualWithGroupsJson struct {
	VisualJSON
	Groups []GroupJSON `json:"groups"`
}

type VisualCollectionJSON struct {
	Visuals []VisualJSON `json:"visuals"`
}

func MapVisual(showId uuid.UUID, visual *shows.Visual) VisualJSON {
	return VisualJSON{
		ID:     visual.ID,
		Name:   visual.Name,
		ShowId: showId,
	}
}

func MapVisualWithGroups(showId uuid.UUID, visual *shows.Visual) VisualWithGroupsJson {
	data := VisualWithGroupsJson{
		VisualJSON: MapVisual(showId, visual),
		Groups:     make([]GroupJSON, len(visual.Groups())),
	}

	for i, group := range visual.Groups() {
		data.Groups[i] = MapGroup(visual.ID, group)
	}

	return data
}

func MapAllVisualsFromShows(shows []*shows.Show) VisualCollectionJSON {
	var visuals = make([]VisualJSON, 0)

	for _, show := range shows {
		for _, visual := range show.Visuals() {
			visuals = append(visuals, MapVisual(show.ID, visual))
		}
	}

	return VisualCollectionJSON{
		Visuals: visuals,
	}
}
