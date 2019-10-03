package mapper

import (
	"github.com/google/uuid"
	"github.com/light-bull/lightbull/shows"
)

type ShowCollectionJSON struct {
	Shows []ShowWithVisualIdsJSON `json:"shows"`
}

type ShowJSON struct {
	ID       uuid.UUID `json:"id"`
	Name     string    `json:"name"`
	Favorite bool      `json:"favorite"`
}

type ShowWithVisualIdsJSON struct {
	ShowJSON
	VisualIds []uuid.UUID `json:"visualIds"`
}

type ShowWithVisualsJSON struct {
	ShowJSON
	Visuals []VisualJSON `json:"visuals"`
}

func MapShow(show *shows.Show) ShowJSON {
	return ShowJSON{
		ID:       show.ID,
		Name:     show.Name,
		Favorite: show.Favorite,
	}
}

func MapShowWithVisualIds(show *shows.Show) ShowWithVisualIdsJSON {
	data := ShowWithVisualIdsJSON{
		ShowJSON:  MapShow(show),
		VisualIds: make([]uuid.UUID, len(show.Visuals())),
	}

	for i, visual := range show.Visuals() {
		data.VisualIds[i] = visual.ID
	}

	return data
}

func MapShowWithVisuals(show *shows.Show) ShowWithVisualsJSON {
	data := ShowWithVisualsJSON{
		ShowJSON: MapShow(show),
		Visuals:  make([]VisualJSON, len(show.Visuals())),
	}

	for i, visual := range show.Visuals() {
		data.Visuals[i] = MapVisual(visual)
	}

	return data
}

func MapShows(shows []*shows.Show) ShowCollectionJSON {
	data := ShowCollectionJSON{
		Shows: make([]ShowWithVisualIdsJSON, len(shows)),
	}

	for i, show := range shows {
		data.Shows[i] = MapShowWithVisualIds(show)
	}

	return data
}
