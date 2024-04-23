package widget

import (
	"image"
	"testing"

	"github.com/utopiagio/gio/f32"
	"github.com/utopiagio/gio/io/event"
	"github.com/utopiagio/gio/io/input"
	"github.com/utopiagio/gio/io/pointer"
	"github.com/utopiagio/gio/io/transfer"
	"github.com/utopiagio/gio/layout"
	"github.com/utopiagio/gio/op"
	"github.com/utopiagio/gio/op/clip"
)

func TestDraggable(t *testing.T) {
	var r input.Router
	gtx := layout.Context{
		Constraints: layout.Exact(image.Pt(100, 100)),
		Source:      r.Source(),
		Ops:         new(op.Ops),
	}

	drag := &Draggable{
		Type: "file",
	}
	tgt := new(int)
	defer pointer.PassOp{}.Push(gtx.Ops).Pop()
	dims := drag.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return layout.Dimensions{Size: gtx.Constraints.Min}
	}, nil)
	stack := clip.Rect{Max: dims.Size}.Push(gtx.Ops)
	event.Op(gtx.Ops, tgt)
	stack.Pop()

	drag.Update(gtx)
	r.Event(transfer.TargetFilter{Target: tgt, Type: drag.Type})
	r.Frame(gtx.Ops)
	r.Queue(
		pointer.Event{
			Position: f32.Pt(10, 10),
			Kind:     pointer.Press,
		},
		pointer.Event{
			Position: f32.Pt(20, 10),
			Kind:     pointer.Move,
		},
		pointer.Event{
			Position: f32.Pt(20, 10),
			Kind:     pointer.Release,
		},
	)
	ofr := &offer{data: "hello"}
	drag.Update(gtx)
	r.Event(transfer.TargetFilter{Target: tgt, Type: drag.Type})
	drag.Offer(gtx, "file", ofr)

	e, ok := r.Event(transfer.TargetFilter{Target: tgt, Type: drag.Type})
	if !ok {
		t.Fatalf("expected event")
	}
	ev := e.(transfer.DataEvent)
	if got, want := ev.Type, "file"; got != want {
		t.Errorf("expected %v; got %v", got, want)
	}
	if ofr.closed {
		t.Error("offer closed prematurely")
	}
	e, ok = r.Event(transfer.TargetFilter{Target: tgt, Type: drag.Type})
	if !ok {
		t.Fatalf("expected event")
	}
	if _, ok := e.(transfer.CancelEvent); !ok {
		t.Fatalf("expected transfer.CancelEvent event")
	}
	r.Frame(gtx.Ops)
	if !ofr.closed {
		t.Error("offer was not closed")
	}
}

// offer satisfies io.ReadCloser for use in data transfers.
type offer struct {
	data   string
	closed bool
}

func (*offer) Read([]byte) (int, error) { return 0, nil }

func (o *offer) Close() error {
	o.closed = true
	return nil
}
