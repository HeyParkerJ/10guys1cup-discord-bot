package main

import (
	"fmt"
	"testing"

	"github.com/andersfylling/disgord"
)

type mockSession struct{}

func (s *mockSession) handleMsg()

func TestHandleMsg(t *testing.T) {
	t.Run("returns pong when the message is exactly 'ping'", func(t *testing.T) {
		session := "hi"
		data := "hi"
		response := handleMsg(session, data)

	})
}
