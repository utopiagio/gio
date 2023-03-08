package widget

import (
	"fmt"
	"image"
	"testing"

	"github.com/utopiagio/gio/font/gofont"
	"github.com/utopiagio/gio/io/key"
	"github.com/utopiagio/gio/layout"
	"github.com/utopiagio/gio/op"
	"github.com/utopiagio/gio/text"
	"github.com/utopiagio/gio/unit"
)

func TestSelectableZeroValue(t *testing.T) {
	var s Selectable
	if s.Text() != "" {
		t.Errorf("expected zero value to have no text, got %q", s.Text())
	}
	if start, end := s.Selection(); start != 0 || end != 0 {
		t.Errorf("expected start=0, end=0, got start=%d, end=%d", start, end)
	}
	if selected := s.SelectedText(); selected != "" {
		t.Errorf("expected selected text to be \"\", got %q", selected)
	}
	s.SetCaret(5, 5)
	if start, end := s.Selection(); start != 0 || end != 0 {
		t.Errorf("expected start=0, end=0, got start=%d, end=%d", start, end)
	}
}

// Verify that an existing selection is dismissed when you press arrow keys.
func TestSelectableMove(t *testing.T) {
	gtx := layout.Context{
		Ops:    new(op.Ops),
		Locale: english,
	}
	cache := text.NewShaper(gofont.Collection())
	font := text.Font{}
	fontSize := unit.Sp(10)

	str := `0123456789`

	// Layout once to populate e.lines and get focus.
	gtx.Queue = newQueue(key.FocusEvent{Focus: true})
	s := new(Selectable)

	w := func(layout.Context) layout.Dimensions { return layout.Dimensions{} }
	Label{
		Selectable: s,
	}.LayoutSelectable(gtx, cache, text.Font{}, fontSize, str, w)

	testKey := func(keyName string) {
		// Select 345
		s.SetCaret(3, 6)
		if start, end := s.Selection(); start != 3 || end != 6 {
			t.Errorf("expected start=%d, end=%d, got start=%d, end=%d", 3, 6, start, end)
		}
		if expected, got := "345", s.SelectedText(); expected != got {
			t.Errorf("KeyName %s, expected %q, got %q", keyName, expected, got)
		}

		// Press the key
		gtx.Queue = newQueue(key.Event{State: key.Press, Name: keyName})
		Label{
			Selectable: s,
		}.LayoutSelectable(gtx, cache, font, fontSize, str, w)

		if expected, got := "", s.SelectedText(); expected != got {
			t.Errorf("KeyName %s, expected %q, got %q", keyName, expected, got)
		}
	}

	testKey(key.NameLeftArrow)
	testKey(key.NameRightArrow)
	testKey(key.NameUpArrow)
	testKey(key.NameDownArrow)
}

func TestSelectableConfigurations(t *testing.T) {
	gtx := layout.Context{
		Ops:         new(op.Ops),
		Constraints: layout.Exact(image.Pt(300, 300)),
		Locale:      english,
	}
	cache := text.NewShaper(gofont.Collection())
	fontSize := unit.Sp(10)
	font := text.Font{}
	sentence := "\n\n\n\n\n\n\n\n\n\n\n\nthe quick brown fox jumps over the lazy dog"
	w := func(layout.Context) layout.Dimensions { return layout.Dimensions{} }

	for _, alignment := range []text.Alignment{text.Start, text.Middle, text.End} {
		for _, zeroMin := range []bool{true, false} {
			t.Run(fmt.Sprintf("Alignment: %v ZeroMinConstraint: %v", alignment, zeroMin), func(t *testing.T) {
				defer func() {
					if err := recover(); err != nil {
						t.Error(err)
					}
				}()
				if zeroMin {
					gtx.Constraints.Min = image.Point{}
				} else {
					gtx.Constraints.Min = gtx.Constraints.Max
				}
				s := new(Selectable)
				label := Label{
					Alignment:  alignment,
					Selectable: s,
				}
				interactiveDims := label.LayoutSelectable(gtx, cache, font, fontSize, sentence, w)
				staticDims := label.Layout(gtx, cache, font, fontSize, sentence)

				if interactiveDims != staticDims {
					t.Errorf("expected consistent dimensions, static returned %#+v, interactive returned %#+v", staticDims, interactiveDims)
				}
			})
		}
	}
}
