package mapper

import (
	"github.com/google/uuid"
	"github.com/light-bull/lightbull/shows"
)

type VisualJSON struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

type VisualWithShowIdJSON struct {
	VisualJSON
	ShowId uuid.UUID `json:"showId"`
}

type VisualWithGroupsJson struct {
	VisualJSON
	Groups []GroupJSON `json:"groups"`
}

type VisualCollectionJSON struct {
	Visuals []VisualWithShowIdJSON `json:"visuals"`
}

func MapVisual(visual *shows.Visual) VisualJSON {
	return VisualJSON{
		ID:   visual.ID,
		Name: visual.Name,
	}
}

func MapVisualWithShowId(showId uuid.UUID, visual *shows.Visual) VisualWithShowIdJSON {
	return VisualWithShowIdJSON{
		VisualJSON: MapVisual(visual),
		ShowId:     showId,
	}
}

func MapVisualWithGroups(visual *shows.Visual) VisualWithGroupsJson {
	data := VisualWithGroupsJson{
		VisualJSON: MapVisual(visual),
		Groups:     make([]GroupJSON, len(visual.Groups())),
	}

	for i, group := range visual.Groups() {
		data.Groups[i] = MapGroup(group)
	}

	return data
}

func MapAllVisualsFromShows(shows []*shows.Show) VisualCollectionJSON {
	var visuals = make([]VisualWithShowIdJSON, 0)

	for _, show := range shows {
		for _, visual := range show.Visuals() {
			visuals = append(visuals, MapVisualWithShowId(show.ID, visual))
		}
	}

	return VisualCollectionJSON{
		Visuals: visuals,
	}
}
