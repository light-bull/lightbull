package mapper

import (
	"github.com/google/uuid"
	"github.com/light-bull/lightbull/shows"
	"github.com/light-bull/lightbull/shows/effects"
	"github.com/light-bull/lightbull/shows/parameters"
)

type GroupJSON struct {
	ID     uuid.UUID  `json:"id"`
	Parts  []string   `json:"parts"`
	Effect EffectJSON `json:"effect"`
}

type EffectJSON struct {
	Type       string                  `json:"type"`
	Parameters []*parameters.Parameter `json:"parameters"`
}

func MapGroup(group *shows.Group) GroupJSON {
	data := GroupJSON{
		ID:     group.ID,
		Parts:  make([]string, len(group.Parts())),
		Effect: MapEffect(&group.Effect),
	}

	copy(data.Parts, group.Parts())

	return data
}

func MapEffect(effect *effects.Effect) EffectJSON {
	return EffectJSON{
		Type:       (*effect).Type(),
		Parameters: (*effect).Parameters(),
	}
}
