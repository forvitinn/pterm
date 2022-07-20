package pterm

import (
	"fmt"
	"strings"

	"atomicgo.dev/cursor"
	"github.com/forvitinn/pterm/internal"
	"github.com/jroimartin/gocui"
)

// DefaultArea is the default area printer.
var DefaultArea = AreaPrinter{}

// AreaPrinter prints an area which can be updated easily.
// use this printer for live output like charts, algorithm visualizations, simulations and even games.
type AreaPrinter struct {
	RemoveWhenDone bool
	Fullscreen     bool
	Center         bool

	content  string
	isActive bool

	area *cursor.Area
}

type KeyDescriptor struct {
	key gocui.Key
	mod gocui.Modifier
	keyFunc func(*gocui.Gui,*gocui.View) error
}
// GetContent returns the current area content.
func (p *AreaPrinter) GetContent() string {
	return p.content
}

// WithRemoveWhenDone removes the AreaPrinter content after it is stopped.
func (p AreaPrinter) WithRemoveWhenDone(b ...bool) *AreaPrinter {
	p.RemoveWhenDone = internal.WithBoolean(b)
	return &p
}

// WithFullscreen sets the AreaPrinter height the same height as the terminal, making it fullscreen.
func (p AreaPrinter) WithFullscreen(b ...bool) *AreaPrinter {
	p.Fullscreen = internal.WithBoolean(b)
	return &p
}

// WithCenter centers the AreaPrinter content to the terminal.
func (p AreaPrinter) WithCenter(b ...bool) *AreaPrinter {
	p.Center = internal.WithBoolean(b)
	return &p
}

// Update overwrites the content of the AreaPrinter.
// Can be used live.
func (p *AreaPrinter) Update(text ...interface{}) {
	if p.area == nil {
		newArea := cursor.NewArea()
		p.area = &newArea
	}
	str := Sprint(text...)
	p.content = str

	if p.Center {
		str = DefaultCenter.Sprint(str)
	}

	if p.Fullscreen {
		str = strings.TrimRight(str, "\n")
		height := GetTerminalHeight()
		contentHeight := strings.Count(str, "\n")

		topPadding := 0
		bottomPadding := height - contentHeight - 2

		if p.Center {
			topPadding = (bottomPadding / 2) + 1
			bottomPadding /= 2
		}

		if height > contentHeight {
			str = strings.Repeat("\n", topPadding) + str
			str += strings.Repeat("\n", bottomPadding)
		}
	}
	p.area.Update(str)
}

// Start the AreaPrinter.
func (p *AreaPrinter) Start(text ...interface{}) (*AreaPrinter, error) {
	p.isActive = true
	str := Sprint(text...)
	newArea := cursor.NewArea()
	p.area = &newArea

	p.Update(str)

	return p, nil
}

// Stop terminates the AreaPrinter immediately.
// The AreaPrinter will not resolve into anything.
func (p *AreaPrinter) Stop() error {
	p.isActive = false
	if p.RemoveWhenDone {
		p.Clear()
	}
	return nil
}

// GenericStart runs Start, but returns a LivePrinter.
// This is used for the interface LivePrinter.
// You most likely want to use Start instead of this in your program.
func (p *AreaPrinter) GenericStart() (*LivePrinter, error) {
	_, _ = p.Start()
	lp := LivePrinter(p)
	return &lp, nil
}

// GenericStop runs Stop, but returns a LivePrinter.
// This is used for the interface LivePrinter.
// You most likely want to use Stop instead of this in your program.
func (p *AreaPrinter) GenericStop() (*LivePrinter, error) {
	_ = p.Stop()
	lp := LivePrinter(p)
	return &lp, nil
}

func (area AreaPrinter) HandleGocui(KeyBinds map[string]KeyDescriptor, managers ...gocui.Manager) (err error) {
	// pause pterm ui
	// used to copy data before hiding isActive 
	area.RemoveWhenDone = true
	area_copy := area
	area.Stop()
	// ...
	// gocui initialize
	g, err := gocui.NewGui(gocui.Output256)
	if err != nil {
		return err
	}
	//defer the showing of pterm ui
	defer func() {
		g.Close()
		area_copy.Start()
	} ()
	g.SetManager(managers...)
	// set KeyBinds
	for view, val := range KeyBinds {
		if err := g.SetKeybinding(view,val.key,val.mod,val.keyFunc); err != nil {
			fmt.Printf("Couldn't establish keybind for %s\n",view)
			return err
		}
	}
	//cleanup gocui
	return err
}
// Wrapper function that clears the content of the Area.
// Moves the cursor to the bottom of the terminal, clears n lines upwards from
// the current position and moves the cursor again.
func (p *AreaPrinter) Clear() {
	p.area.Clear()
}
