package component

import (
	"math/rand"

	"github.com/iancanderson/spookypaths/game/core"
	"github.com/yohamta/donburi"
)

func RandomEnergyType() core.EnergyType {
	return core.EnergyType(rand.Intn(int(core.Psychic) + 1))
}

type EnergyData struct {
	EnergyType core.EnergyType
}

var Energy = donburi.NewComponentType[EnergyData]()
