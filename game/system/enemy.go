package system

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/iancanderson/gandermerge/game/component"
	"github.com/iancanderson/gandermerge/game/config"
	"github.com/iancanderson/gandermerge/game/layers"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
	"github.com/yohamta/donburi/filter"
	"github.com/yohamta/donburi/query"
)

type enemy struct {
	images map[component.EnergyType]*ebiten.Image
	query  *query.Query
}

var Enemy = &enemy{
	query: ecs.NewQuery(
		layers.LayerEnemy,
		filter.Contains(
			component.Sprite,
		)),
}

const enemyWidth = 309

func (e *enemy) Startup(ecs *ecs.ECS) {
	//TODO: share these with orb_spawner?
	e.images = loadEnergyTypeImages()

	entity := ecs.Create(
		layers.LayerEnemy,
		component.Energy,
		component.Sprite,
	)
	entry := ecs.World.Entry(entity)

	energyType := component.RandomEnergyType()
	donburi.SetValue(entry, component.Energy,
		component.EnergyData{
			EnergyType: energyType,
		})

	donburi.SetValue(entry, component.Sprite,
		component.SpriteData{
			Image:     e.images[energyType],
			X:         config.WindowWidth/2 - enemyWidth/2,
			Y:         100,
			Scale:     0.5,
			GreenTint: energyType == component.Poison,
			RedTint:   energyType == component.Fire,
		})
}

func (e *enemy) Update(ecs *ecs.ECS) {
}

func (e *enemy) Draw(ecs *ecs.ECS, screen *ebiten.Image) {
	// TODO: consolidate with render.go
	e.query.EachEntity(ecs.World, func(entry *donburi.Entry) {
		sprite := component.GetSprite(entry)

		op := &ebiten.DrawImageOptions{}
		op.GeoM.Scale(sprite.Scale, sprite.Scale)
		op.GeoM.Translate(sprite.X, sprite.Y)
		op.Filter = ebiten.FilterLinear
		if sprite.GreenTint {
			op.ColorM.Scale(0.5, 1.0, 0.5, 1)
		} else if sprite.RedTint {
			op.ColorM.Scale(1.0, 0.5, 0.5, 1)
		}

		screen.DrawImage(sprite.Image, op)
	})
}
