package pterm

import ( 
	"github.com/jroimartin/gocui"
)

type KeyDescriptor struct {
	Key     gocui.Key
	Mod     gocui.Modifier
	KeyFunc func(*gocui.Gui, *gocui.View) error
}

type GLayout func(*gocui.Gui) error
type GKeyBindFunc func(*gocui.Gui, *gocui.View) error
type GMapView2KeyDesc map[string]KeyDescriptor
