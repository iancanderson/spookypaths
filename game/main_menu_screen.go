package game

import (
	"image"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"github.com/iancanderson/spookypaths/game/assets"
	"github.com/iancanderson/spookypaths/game/config"
	"github.com/iancanderson/spookypaths/game/uicomponent"
	"github.com/yohamta/furex/v2"
	"golang.org/x/image/font"
)

type MainMenuScreen struct {
	gameUI           *furex.View
	startDailyLevel  StartDailyLevel
	startRandomLevel StartRandomLevel
}

type StartDailyLevel func()
type StartRandomLevel func()

func NewMainMenuScreen(startDailyLevel StartDailyLevel, startRandomLevel StartRandomLevel) *MainMenuScreen {
	g := &MainMenuScreen{
		startDailyLevel:  startDailyLevel,
		startRandomLevel: startRandomLevel,
	}
	g.setupMenuUI()
	return g
}

func (g *MainMenuScreen) Draw(screen *ebiten.Image) {
	g.gameUI.Draw(screen)
}

func (g *MainMenuScreen) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return config.WindowWidth, config.WindowHeight
}

func (g *MainMenuScreen) Update() error {
	g.gameUI.Update()
	return nil
}

type TextComponent struct {
	Str      string
	fontface font.Face
}

func (textComponent *TextComponent) HandleDraw(screen *ebiten.Image, frame image.Rectangle) {
	fontface := textComponent.fontface
	textBounds := text.BoundString(fontface, textComponent.Str)
	xOffset := frame.Dx()/2 - textBounds.Dx()/2
	yOffset := frame.Dy()/2 + textBounds.Dy()/2

	// Not sure why we need to do this
	yOffset -= 8

	text.Draw(screen, textComponent.Str, fontface, frame.Min.X+xOffset, frame.Min.Y+yOffset, color.White)
}

func (g *MainMenuScreen) setupMenuUI() {
	g.gameUI = &furex.View{
		Width:        config.WindowWidth,
		Height:       config.WindowHeight,
		Direction:    furex.Column,
		Justify:      furex.JustifyCenter,
		AlignItems:   furex.AlignItemCenter,
		AlignContent: furex.AlignContentCenter,
		Wrap:         furex.Wrap,
	}

	g.gameUI.AddChild(&furex.View{
		Width:   config.WindowWidth - 100,
		Height:  120,
		Handler: &TextComponent{fontface: assets.FontManager.Creepster160, Str: "Spooky"},
	})
	g.gameUI.AddChild(&furex.View{
		Width:        config.WindowWidth - 100,
		Height:       140,
		Handler:      &TextComponent{fontface: assets.FontManager.Creepster200, Str: "Paths"},
		MarginBottom: 200,
	})

	g.gameUI.AddChild(&furex.View{
		MarginBottom: 20,
		Width:        600,
		Height:       100,
		Handler: &uicomponent.Button{
			Text:    "Play random game",
			OnClick: g.startRandomLevel,
		},
	})

	g.gameUI.AddChild(&furex.View{
		MarginBottom: 20,
		Width:        600,
		Height:       100,
		Handler: &uicomponent.Button{
			Text:    "Play daily game",
			OnClick: g.startDailyLevel,
		},
	})
}
