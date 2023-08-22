package top

import (
	"fmt"
	"testing"
	"time"
)

func TestNewMultiChanObserver(t *testing.T) {
	o := NewMultiChanObserver()

	if len(o.chMap) != 0 {
		t.Errorf("Expected 0 channels, got %d", len(o.chMap))
	}

	if cap(o.events) != 10 {
		t.Errorf("Expected events channel capacity of 10, got %d", cap(o.events))
	}
}

func TestAddChannel(t *testing.T) {
	o := NewMultiChanObserver()
	ch := o.AddChannel("test")

	if len(o.chMap) != 1 {
		t.Errorf("Expected 1 channel, got %d", len(o.chMap))
	}

	if cap(*ch) != 5 {
		t.Errorf("Expected channel capacity of 5, got %d", cap(*ch))
	}
}

func TestRemoveChannel(t *testing.T) {
	o := NewMultiChanObserver()
	o.AddChannel("test")
	o.RemoveChannel("test")

	if len(o.chMap) != 0 {
		t.Errorf("Expected 0 channels, got %d", len(o.chMap))
	}
}

func TestNotify(t *testing.T) {
	o := NewMultiChanObserver()
	pp := &Program{}

	o.Notify(pp)
	o.Notify(nil)
	if len(o.events) != 2 {
		t.Errorf("Expected 2 event in channel, got %d", len(o.events))
	}
}
func TestNotifyZero(t *testing.T) {
	o := NewMultiChanObserver()
	pp := &Program{}

	o.Notify(pp)
	o.Notify(nil)
	if len(o.events) != 2 {
		t.Errorf("Expected 2 event in channel, got %d", len(o.events))
	}
	o.AddChannel("test")
	time.Sleep(100 * time.Millisecond)
	if len(o.events) != 0 {
		t.Errorf("Expected 0 event in channel, got %d", len(o.events))
	}

}

func TestNotifyPoll(t *testing.T) {
	o := NewMultiChanObserver()
	go o.onPoolChanged()
	o.AddChannel("test")
	pp := &Program{}

	o.Notify(pp)
	o.Notify(nil)
	time.Sleep(100 * time.Millisecond)

	if len(o.events) != 0 {
		t.Errorf("Expected 0 event in channel, got %d", len(o.events))
	}
}

func TestNotifyPoll2(t *testing.T) {
	o := NewMultiChanObserver()
	go o.onPoolChanged()
	ch1 := o.AddChannel("test")
	ch2 := o.AddChannel("test2")

	pp := &Program{filepath: "test"}

	o.Notify(pp)
	o.Notify(nil)
	o.Notify(nil)

	time.Sleep(100 * time.Millisecond)

	if len(o.events) != 0 {
		t.Errorf("Expected 0 event in channel, got %d", len(o.events))
	}
	if len(*ch1) != 3 {
		t.Errorf("Expected 3 event in channel, got %d", len(*ch1))
	}
	if len(*ch2) != 3 {
		t.Errorf("Expected 3 event in channel, got %d", len(*ch2))
	}

	go func() {
		fmt.Println("ch1:", <-*ch1)
		fmt.Println("ch2:", <-*ch2)

		fmt.Println("ch1:", <-*ch1)
		fmt.Println("ch2:", <-*ch2)

	}()
	time.Sleep(100 * time.Millisecond) // заменили time.After на time.Sleep
	if len(o.events) != 0 {
		t.Errorf("Expected 0 event in channel, got %d", len(o.events))
	}
	if len(*ch1) != 1 {
		t.Errorf("Expected 1 event in channel, got %d", len(*ch1))
	}
	if len(*ch2) != 1 {
		t.Errorf("Expected 1 event in channel, got %d", len(*ch2))
	}
}
