package radio

import (
	"dubmer-bono/app/types/entity"
	"dubmer-bono/app/utility"
	"io"
	"testing"
	"time"
)

func newTestService() *Service {
	return &Service{
		getDriverPressure: func() int { return 0 },
		prio_sorted_vc:    make(map[int]*utility.ExpiryQueue[entity.RadioMessage]),
	}
}

func pushMessage(t *testing.T, s *Service, priority int, message entity.RadioPayload) {
	t.Helper()
	q, ok := s.prio_sorted_vc[priority]
	if !ok {
		queue := utility.NewQueue[entity.RadioMessage]()
		q = &queue
		s.prio_sorted_vc[priority] = q
	}
	q.Push(entity.RadioMessage{
		Priority: priority,
		Message:  message,
		Type:     entity.DIRECT,
		Expiry:   time.Now().Add(time.Minute),
	})
}

func TestGetMessageByMinPriorityEmpty(t *testing.T) {
	s := newTestService()
	if _, ok := s.GetMessageByMinPriority(); ok {
		t.Fatal("expected no message from empty queue")
	}
}

func TestGetMessageByMinPriorityDirect(t *testing.T) {
	s := newTestService()
	pushMessage(t, s, 3, entity.DirectMessage{Text: "hello"})

	reader, ok := s.GetMessageByMinPriority()
	if !ok {
		t.Fatal("expected a message")
	}
	data, err := io.ReadAll(reader)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if string(data) != "hello" {
		t.Fatalf("got %q, want %q", string(data), "hello")
	}
}

func TestGetMessageByMinPriorityReturnsHighestPriority(t *testing.T) {
	s := newTestService()
	pushMessage(t, s, 1, entity.DirectMessage{Text: "low"})
	pushMessage(t, s, 5, entity.DirectMessage{Text: "high"})

	reader, ok := s.GetMessageByMinPriority()
	if !ok {
		t.Fatal("expected a message")
	}
	data, _ := io.ReadAll(reader)
	if string(data) != "high" {
		t.Fatalf("got %q, want %q", string(data), "high")
	}
}

func TestGetMessageByMinPriorityFunctionMessage(t *testing.T) {
	s := newTestService()
	pushMessage(t, s, 2, entity.FunctionMessage{Fn: func() string { return "function-output" }})

	reader, ok := s.GetMessageByMinPriority()
	if !ok {
		t.Fatal("expected a message")
	}
	data, _ := io.ReadAll(reader)
	if string(data) != "function-output" {
		t.Fatalf("got %q, want %q", string(data), "function-output")
	}
}

func TestGetMessageByMinPriorityRespectsPressure(t *testing.T) {
	s := newTestService()
	s.getDriverPressure = func() int { return 3 }
	pushMessage(t, s, 2, entity.DirectMessage{Text: "low"})
	pushMessage(t, s, 3, entity.DirectMessage{Text: "ok"})

	reader, ok := s.GetMessageByMinPriority()
	if !ok {
		t.Fatal("expected a message")
	}
	data, _ := io.ReadAll(reader)
	if string(data) != "ok" {
		t.Fatalf("got %q, want %q (only messages at or above pressure should surface)", string(data), "ok")
	}
}

func TestGetMessageByMinPriorityMuted(t *testing.T) {
	s := newTestService()
	s.muted = true
	pushMessage(t, s, 4, entity.DirectMessage{Text: "quiet"})
	pushMessage(t, s, 5, entity.DirectMessage{Text: "top"})

	reader, ok := s.GetMessageByMinPriority()
	if !ok {
		t.Fatal("expected a message")
	}
	data, _ := io.ReadAll(reader)
	if string(data) != "top" {
		t.Fatalf("got %q, want %q (muted should only surface max-priority)", string(data), "top")
	}
}

func TestMessageChanListnerPopulatesQueue(t *testing.T) {
	s := newTestService()
	s.msg_chan = make(chan entity.RadioMessage)
	done := make(chan struct{})
	go func() {
		s.MessageChanListner()
		close(done)
	}()

	s.msg_chan <- entity.RadioMessage{Priority: 4, Message: entity.DirectMessage{Text: "x"}, Expiry: time.Now().Add(time.Minute)}
	close(s.msg_chan)
	<-done

	queue, ok := s.prio_sorted_vc[4]
	if !ok {
		t.Fatal("expected queue for priority 4")
	}
	if _, ok := queue.Pop(); !ok {
		t.Fatal("expected message in queue")
	}
}
