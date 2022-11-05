package system

import (
	"github.com/iancanderson/gandermerge/game/component"
	"github.com/iancanderson/gandermerge/game/layers"
	"github.com/iancanderson/gandermerge/game/util"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
	"github.com/yohamta/donburi/filter"
	"github.com/yohamta/donburi/query"
)

type orbChain struct {
	orbs       []*donburi.Entry
	energyType component.EnergyType
}

func (o *orbChain) IsSecondToLast(entry *donburi.Entry) bool {
	return len(o.orbs) > 1 && o.orbs[len(o.orbs)-2] == entry
}

func (o *orbChain) Add(entry *donburi.Entry) bool {
	if o.contains(entry) {
		return false
	}
	if o.energyType != component.GetEnergy(entry).EnergyType {
		return false
	}
	if len(o.orbs) > 0 {
		// Only allow adjacent orbs to be added to the chain
		lastOrb := o.orbs[len(o.orbs)-1]
		if !component.GetGridPosition(lastOrb).IsAdjacent(component.GetGridPosition(entry)) {
			return false
		}
	}

	o.orbs = append(o.orbs, entry)
	return true
}

func (o *orbChain) CanBeMerged() bool {
	return len(o.orbs) >= 3
}

func (o *orbChain) Pop() *donburi.Entry {
	// Return the last orb in the chain and remove it from the chain
	i := len(o.orbs) - 1
	orb := o.orbs[i]
	o.orbs = o.orbs[:i]
	return orb
}

func (o *orbChain) Len() int {
	return len(o.orbs)
}

func (o *orbChain) contains(entry *donburi.Entry) bool {
	for _, orb := range o.orbs {
		if orb == entry {
			return true
		}
	}
	return false
}

type input struct {
	chain           *orbChain
	inputSource     util.InputSource
	scoreQuery      *query.Query
	selectableQuery *query.Query
}

var Input = &input{
	selectableQuery: ecs.NewQuery(
		layers.LayerOrbs,
		filter.Contains(
			component.Energy,
			component.Selectable,
			component.Sprite,
		)),
	scoreQuery: ecs.NewQuery(
		layers.LayerScoreboard,
		filter.Contains(
			component.Score,
		)),
}

func (r *input) Update(ecs *ecs.ECS) {
	// Check if game is over
	score, ok := r.scoreQuery.FirstEntity(ecs.World)
	if ok && component.GetScore(score).IsGameOver() {
		return
	}

	if r.chain == nil {
		r.inputSource = util.JustPressedInputSource()
		if r.inputSource != nil {
			r.handleInputPressed(ecs)
		}
	}

	if r.inputSource != nil && r.inputSource.JustReleased() {
		r.clearOrbChain(ecs.World)
	}

	if r.chain != nil {
		r.selectableQuery.EachEntity(ecs.World, func(entry *donburi.Entry) {
			sprite := component.GetSprite(entry)
			inputX, inputY := r.inputSource.Position()
			if sprite.In(inputX, inputY) {
				if r.chain.Add(entry) {
					selectable := component.GetSelectable(entry)
					selectable.Selected = true
				} else if r.chain.IsSecondToLast(entry) {
					popped := r.chain.Pop()
					selectable := component.GetSelectable(popped)
					selectable.Selected = false
				}
			}
		})
	}
}

func (r *input) handleInputPressed(ecs *ecs.ECS) {
	r.selectableQuery.EachEntity(ecs.World, func(entry *donburi.Entry) {
		sprite := component.GetSprite(entry)
		inputX, inputY := r.inputSource.Position()
		if sprite.In(inputX, inputY) {
			r.chain = r.createOrbChain(entry)
			selectable := component.GetSelectable(entry)
			selectable.Selected = true
		}
	})
}

func (r *input) createOrbChain(entry *donburi.Entry) *orbChain {
	chain := &orbChain{
		orbs:       []*donburi.Entry{entry},
		energyType: component.GetEnergy(entry).EnergyType,
	}
	return chain
}

func (r *input) clearOrbChain(world donburi.World) {
	if r.chain == nil {
		return
	}

	if r.chain.CanBeMerged() {
		energyEmitted := r.chain.Len()
		entry, ok := r.scoreQuery.FirstEntity(world)
		if ok {
			score := component.GetScore(entry)
			score.EnergyToWin -= energyEmitted
			score.MovesRemaining--
		}
	}

	for _, orb := range r.chain.orbs {
		if r.chain.CanBeMerged() {
			world.Remove(orb.Entity())
		} else {
			selectable := component.GetSelectable(orb)
			selectable.Selected = false
		}
	}
	r.chain = nil
}