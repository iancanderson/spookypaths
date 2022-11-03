package system

import (
	"bytes"
	"image"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/iancanderson/gandermerge/assets/images"
	"github.com/iancanderson/gandermerge/component"
	"github.com/iancanderson/gandermerge/layers"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
)

const columns = 8
const rows = 8
const columnWidth = 48
const rowHeight = 48

type OrbSpawner struct {
}

func NewOrbSpawner() *OrbSpawner {
	return &OrbSpawner{}
}

func (s *OrbSpawner) Update(ecs *ecs.ECS) {
}

func (s *OrbSpawner) Startup(ecs *ecs.ECS) {
	orbs := ecs.CreateMany(
		layers.LayerOrbs,
		rows*columns,
		component.Position,
		component.Sprite,
		component.Energy,
	)

	images := loadEnergyTypeImages()

	for row := 0; row < rows; row++ {
		for col := 0; col < columns; col++ {
			entry := ecs.World.Entry(orbs[row*columns+col])

			donburi.SetValue(entry, component.Position,
				component.PositionData{
					X: float64(col) * columnWidth,
					Y: float64(row) * rowHeight,
				})

			energyType := component.RandomEnergyType()
			donburi.SetValue(entry, component.Energy,
				component.EnergyData{
					EnergyType: energyType,
				})

			donburi.SetValue(entry, component.Sprite,
				component.SpriteData{Image: images[energyType]})
		}
	}
}

func loadEnergyTypeImages() map[component.EnergyType]*ebiten.Image {
	return map[component.EnergyType]*ebiten.Image{
		component.Electric: loadImage(images.Electric_png),
		component.Fire:     loadImage(images.Fire_png),
		component.Ghost:    loadImage(images.Ghost_png),
		component.Poison:   loadImage(images.Poison_png),
		component.Psychic:  loadImage(images.Psychic_png),
	}
}

func loadImage(data []byte) *ebiten.Image {
	img, _, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		log.Fatal(err)
	}
	return ebiten.NewImageFromImage(img)
}
