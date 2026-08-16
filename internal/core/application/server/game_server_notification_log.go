package server

import (
	"github.com/justjack1521/mevium/pkg/genproto/protomulti"
	"sync"
)

// notificationLogHardCap is a backstop, not the working bound. Turn-boundary
// compaction is what actually keeps the log proportional to turns played; this
// only catches a game that somehow never reaches a boundary. Entries are
// evicted oldest first, and a client that asked for something evicted finds out
// from the from_sequence its catch-up reports back.
const notificationLogHardCap = 4096

// notificationEntry is one broadcast notification, kept exactly as it went out:
// the same marshalled payload, kind and ordinal, so a replay is byte-identical
// to the original and indistinguishable from live traffic on the client.
type notificationEntry struct {
	sequence uint64
	kind     protomulti.MultiGameNotificationType
	payload  []byte
}

// notificationLog assigns each broadcast its ordinal and retains it for the
// life of the game, so a client that missed some can be sent them again. There
// is no TTL and no reaper: the log is a field on GameServer and dies with it
// when the host unregisters the game.
//
// The ordinal is assigned here rather than in the domain deliberately. Not
// every change is broadcast — party and player adds are registry bookkeeping,
// and the catch-up request is itself a change — so numbering the change stream
// would leave the client staring at gaps it could never fill. Numbering at the
// point of publication makes the sequence gapless in exactly the stream the
// client sees, which is what makes "seq > last + 1" a reliable signal.
//
// It carries its own mutex rather than reusing GameServer.mu, which guards the
// client map and is deliberately dropped before any network I/O.
type notificationLog struct {
	mu      sync.Mutex
	entries []notificationEntry
	// lastConfirm is the sequence of the most recent queue confirm, and so the
	// floor of the turn currently in progress.
	lastConfirm uint64
	// latest is the highest sequence assigned, retained or not, so an empty
	// replay can still report where the client stands.
	latest uint64
}

// Append records a broadcast notification and returns the ordinal it was given.
// It is called before the fan-out, not after, which is what makes a publish
// that errors recoverable rather than lost.
//
// Sequences start at 1; zero is reserved on the wire for "unsequenced", which
// is what everything outside a live game — lobby traffic sharing the same
// envelope — carries.
func (l *notificationLog) Append(kind protomulti.MultiGameNotificationType, payload []byte) uint64 {

	l.mu.Lock()
	defer l.mu.Unlock()

	l.latest++
	var sequence = l.latest

	l.entries = append(l.entries, notificationEntry{
		sequence: sequence,
		kind:     kind,
		payload:  payload,
	})

	// A queue confirm carries every party's full committed action queue, which
	// makes it a natural checkpoint: everything the turn's fiddling produced is
	// now restated in one message.
	if kind == protomulti.MultiGameNotificationType_GAME_NOTIFY_QUEUE_CONFIRM {
		l.compact(sequence)
	}

	l.trim()

	return sequence

}

// Replay returns every retained entry at or after from, along with the range it
// actually covers. When nothing matches, the reported range is deliberately
// inverted (from > to) so a caller can tell "nothing to send" apart from "here
// is the start of the game".
func (l *notificationLog) Replay(from uint64) ([]notificationEntry, uint64, uint64) {

	l.mu.Lock()
	defer l.mu.Unlock()

	var out []notificationEntry
	for _, entry := range l.entries {
		if entry.sequence >= from {
			out = append(out, entry)
		}
	}

	if len(out) == 0 {
		return nil, l.latest + 1, l.latest
	}

	return out, out[0].sequence, out[len(out)-1].sequence

}

// compact drops the entries the queue confirm at sequence supersedes: the
// enqueue, dequeue and lock traffic of the turn it closes. The next player turn
// resets action and lock state outright, so none of it is needed to reconstruct
// the present — only the confirm is. Callers hold l.mu.
func (l *notificationLog) compact(confirm uint64) {

	var kept = l.entries[:0]
	for _, entry := range l.entries {
		if entry.sequence > l.lastConfirm && entry.sequence < confirm && supersededByQueueConfirm(entry.kind) {
			continue
		}
		kept = append(kept, entry)
	}

	l.entries = kept
	l.lastConfirm = confirm

}

// trim enforces the hard cap by dropping the oldest entries. Callers hold l.mu.
func (l *notificationLog) trim() {
	if len(l.entries) <= notificationLogHardCap {
		return
	}
	var excess = len(l.entries) - notificationLogHardCap
	l.entries = append(l.entries[:0], l.entries[excess:]...)
}

func supersededByQueueConfirm(kind protomulti.MultiGameNotificationType) bool {
	switch kind {
	case protomulti.MultiGameNotificationType_GAME_NOTIFY_ENQUEUE_ACTION,
		protomulti.MultiGameNotificationType_GAME_NOTIFY_DEQUEUE_ACTION,
		protomulti.MultiGameNotificationType_GAME_NOTIFY_LOCK_ACTION:
		return true
	}
	return false
}
