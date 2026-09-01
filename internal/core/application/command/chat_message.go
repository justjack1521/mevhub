package command

import (
	"errors"
	"strings"
	"unicode/utf8"
)

// MaxChatMessageLength caps a chat message in runes, not bytes, so the limit
// reads the same to a client counting characters regardless of script.
const MaxChatMessageLength = 256

var (
	ErrChatMessageEmpty   = errors.New("chat message is empty")
	ErrChatMessageTooLong = errors.New("chat message exceeds maximum length")
)

// sanitiseChatMessage is the single validation point for both lobby and game
// chat: the wire contract carries the message verbatim, so anything accepted
// here reaches every client in the group.
func sanitiseChatMessage(message string) (string, error) {
	var trimmed = strings.TrimSpace(message)
	if trimmed == "" {
		return "", ErrChatMessageEmpty
	}
	if utf8.RuneCountInString(trimmed) > MaxChatMessageLength {
		return "", ErrChatMessageTooLong
	}
	return trimmed, nil
}
