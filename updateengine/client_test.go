package updateengine

import (
	"reflect"
	"testing"
	"time"

	"github.com/coreos/locksmith/Godeps/_workspace/src/github.com/godbus/dbus"
)

func makeSig(curOp string) *dbus.Signal {
	return &dbus.Signal{
		Body: []interface{}{
			int64(0),
			0.0,
			curOp,
			"newVer",
			int64(1024),
		},
	}
}

func makeStat(curOp string) Status {
	return Status{
		0,
		0.0,
		curOp,
		"newVer",
		1024,
	}
}

func TestRebootNeededSignal(t *testing.T) {
	c := &Client{
		ch: make(chan *dbus.Signal, signalBuffer),
	}
	r := make(chan Status)
	s := make(chan struct{})
	var done bool
	go func() {
		c.RebootNeededSignal(r, s)
		done = true
	}()

	if done {
		t.Fatal("RebootNeededSignal stopped prematurely")
	}
	c.ch <- makeSig(UpdateStatusUpdatedNeedReboot)
	if done {
		t.Fatal("RebootNeededSignal stopped prematurely")
	}

	time.Sleep(10 * time.Millisecond)

	select {
	case stat := <-r:
		if !reflect.DeepEqual(stat, makeStat(UpdateStatusUpdatedNeedReboot)) {
			t.Fatalf("bad status received: %#v", stat)
		}
	default:
		t.Fatal("RebootNeededSignal did not send expected Status update")
	}

	if done {
		t.Fatal("RebootNeededSignal stopped prematurely")
	}

	c.ch <- makeSig("some other ignored signal")

	time.Sleep(10 * time.Millisecond)

	select {
	case stat := <-r:
		t.Fatalf("unexpected status on unknown signal: %#v", stat)
	default:
	}

	if done {
		t.Fatal("RebootNeededSignal stopped prematurely")
	}

	close(s)

	time.Sleep(10 * time.Millisecond)

	if !done {
		t.Fatal("RebootNeededSignal did not stop as expected")
	}
}
