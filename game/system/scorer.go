package system

import (
	"github.com/iancanderson/spookypaths/game/component"
	"github.com/iancanderson/spookypaths/game/config"
	"github.com/iancanderson/spookypaths/game/layers"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
	"github.com/yohamta/donburi/filter"
	"github.com/yohamta/donburi/query"
)

type scorer struct {
	enemiesQuery *query.Query
	scoreQuery   *query.Query
}

var Scorer = &scorer{
	scoreQuery: ecs.NewQuery(
		layers.LayerUI,
		filter.Contains(
			component.Score,
		)),
	enemiesQuery: ecs.NewQuery(
		layers.LayerEnemy,
		filter.Contains(
			component.Hitpoints,
		)),
}

func (s *scorer) Startup(ecs *ecs.ECS) {
	entity := ecs.Create(layers.LayerUI, component.Score)
	entry := ecs.World.Entry(entity)

	component.Score.Set(entry, &component.ScoreData{
		MovesRemaining: config.MovesAllowed,
		EnemiesAreDead: false,
	})
}

func (s *scorer) Update(ecs *ecs.ECS) {
	if s.allEnemiesAreDead(ecs) {
		score, ok := s.scoreQuery.FirstEntity(ecs.World)
		if !ok {
			return
		}
		component.Score.Get(score).EnemiesAreDead = true
	}
}

func (s *scorer) allEnemiesAreDead(ecs *ecs.ECS) bool {
	result := true
	s.enemiesQuery.EachEntity(ecs.World, func(entry *donburi.Entry) {
		hp := component.Hitpoints.Get(entry)
		if !hp.IsDead() {
			result = false
		}
	})

	return result
}
