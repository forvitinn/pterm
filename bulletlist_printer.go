package pterm

import (
	"strings"

	"github.com/pterm/pterm/internal"
)

// NewBulletListFromStrings returns a BulletListPrinter with Text using the NewTreeListItemFromString method.
func NewBulletListFromStrings(s []string, padding string) BulletListPrinter {
	var lis []BulletListItem
	for _, line := range s {
		lis = append(lis, NewBulletListItemFromString(line, padding))
	}
	return *DefaultBulletList.WithItems(lis)
}

// NewBulletListItemFromString returns a BulletListItem with a Text. The padding is counted in the Text to define the Level of the ListItem.
func NewBulletListItemFromString(text string, padding string) BulletListItem {
	s, l := internal.RemoveAndCountPrefix(text, padding)
	return BulletListItem{
		Level: l,
		Text:  s,
	}
}

// DefaultListItem contains standards, which can be used to print a ListItem.
var DefaultListItem = BulletListItem{
	Bullet:      "•",
	TextStyle:   &ThemeDefault.BulletListTextStyle,
	BulletStyle: &ThemeDefault.BulletListBulletStyle,
}

// BulletListItem is able to render a ListItem.
type BulletListItem struct {
	Level       int
	Text        string
	TextStyle   *Style
	Bullet      string
	BulletStyle *Style
}

// WithLevel returns a new BulletListItem with a specific Level.
func (p BulletListItem) WithLevel(level int) BulletListItem {
	p.Level = level
	return p
}

// WithText returns a new BulletListItem with a specific Text.
func (p BulletListItem) WithText(text string) BulletListItem {
	p.Text = text
	return p
}

// WithTextStyle returns a new BulletListItem with a specific TextStyle.
func (p BulletListItem) WithTextStyle(style *Style) BulletListItem {
	p.TextStyle = style
	return p
}

// WithBullet returns a new BulletListItem with a specific Prefix.
func (p BulletListItem) WithBullet(bullet string) BulletListItem {
	p.Bullet = bullet
	return p
}

// WithBulletStyle returns a new BulletListItem with a specific BulletStyle.
func (p BulletListItem) WithBulletStyle(style *Style) BulletListItem {
	p.BulletStyle = style
	return p
}

// Render renders the BulletListItem as a string.
func (p BulletListItem) Render() {
	Println(p.Srender())
}

// Srender renders the BulletListItem as a string.
func (p BulletListItem) Srender() string {
	if p.TextStyle == nil {
		p.TextStyle = NewStyle()
	}
	if p.BulletStyle == nil {
		p.BulletStyle = &ThemeDefault.BulletListBulletStyle
	}
	return strings.Repeat(" ", p.Level) + p.BulletStyle.Sprint(p.Bullet) + " " + p.TextStyle.Sprint(p.Text)
}

// NewBulletListFromString returns a BulletListPrinter with Text using the NewTreeListItemFromString method, splitting after return (\n).
func NewBulletListFromString(s string, padding string) BulletListPrinter {
	return NewBulletListFromStrings(strings.Split(s, "\n"), padding)
}

// DefaultBulletList contains standards, which can be used to print a BulletListPrinter.
var DefaultBulletList = BulletListPrinter{
	Bullet:      "•",
	TextStyle:   &ThemeDefault.BulletListTextStyle,
	BulletStyle: &ThemeDefault.BulletListBulletStyle,
}

// BulletListPrinter is able to render a list.
type BulletListPrinter struct {
	Items       []BulletListItem
	TextStyle   *Style
	Bullet      string
	BulletStyle *Style
}

// WithItems returns a new list with specific Items.
func (l BulletListPrinter) WithItems(items []BulletListItem) *BulletListPrinter {
	l.Items = append(l.Items, items...)
	return &l
}

// WithTextStyle returns a new list with a specific text style.
func (l BulletListPrinter) WithTextStyle(style *Style) *BulletListPrinter {
	l.TextStyle = style
	return &l
}

// WithBullet returns a new list with a specific bullet.
func (l BulletListPrinter) WithBullet(bullet string) *BulletListPrinter {
	l.Bullet = bullet
	return &l
}

// WithBulletStyle returns a new list with a specific bullet style.
func (l BulletListPrinter) WithBulletStyle(style *Style) *BulletListPrinter {
	l.BulletStyle = style
	return &l
}

// Render prints the list to the terminal.
func (l BulletListPrinter) Render() error {
	s, _ := l.Srender()
	Println(s)

	return nil
}

// Srender renders the list as a string.
func (l BulletListPrinter) Srender() (string, error) {
	var ret string
	for _, item := range l.Items {
		if item.Bullet == "" {
			ret += item.WithBullet(l.Bullet).Srender() + "\n"
		} else {
			ret += item.WithBullet(item.Bullet).Srender() + "\n"
		}
	}
	return ret, nil
}
