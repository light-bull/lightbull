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

type VisualWithGroupsJSON struct {
	VisualJSON
	Groups []GroupJSON `json:"groups"`
}

type VisualWithGroupIdsJSON struct {
	VisualJSON
	GroupIds []uuid.UUID `json:"groupIds"`
}

type VisualCollectionJSON struct {
	Visuals []VisualWithGroupIdsJSON `json:"visuals"`
}

func MapVisual(showId uuid.UUID, visual *shows.Visual) VisualJSON {
	return VisualJSON{
		ID:     visual.ID,
		Name:   visual.Name,
		ShowId: showId,
	}
}

func MapVisualWithGroups(showId uuid.UUID, visual *shows.Visual) VisualWithGroupsJSON {
	data := VisualWithGroupsJSON{
		VisualJSON: MapVisual(showId, visual),
		Groups:     make([]GroupJSON, len(visual.Groups())),
	}

	for i, group := range visual.Groups() {
		data.Groups[i] = MapGroup(visual.ID, group)
	}

	return data
}

func MapVisualWithGroupIds(showId uuid.UUID, visual *shows.Visual) VisualWithGroupIdsJSON {
	data := VisualWithGroupIdsJSON{
		VisualJSON: MapVisual(showId, visual),
		GroupIds:   make([]uuid.UUID, len(visual.Groups())),
	}

	for i, group := range visual.Groups() {
		data.GroupIds[i] = group.ID
	}

	return data;
}

func MapAllVisualsFromShows(shows []*shows.Show) VisualCollectionJSON {
	var visuals = make([]VisualWithGroupIdsJSON, 0)

	for _, show := range shows {
		for _, visual := range show.Visuals() {
			visuals = append(visuals, MapVisualWithGroupIds(show.ID, visual))
		}
	}

	return VisualCollectionJSON{
		Visuals: visuals,
	}
}
