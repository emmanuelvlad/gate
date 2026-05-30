package proxy

import (
	"strings"
	"testing"

	"go.minekube.com/gate/pkg/edition/java/profile"
	"go.minekube.com/gate/pkg/util/uuid"
)

func TestUnregisterConnection_DoesNotRemoveDifferentPlayerWithSameIdentity(t *testing.T) {
	id := uuid.New()
	oldPlayer := &connectedPlayer{profile: &profile.GameProfile{ID: id, Name: "Notch"}}
	activePlayer := &connectedPlayer{profile: &profile.GameProfile{ID: id, Name: "Notch"}}

	p := &Proxy{
		playerNames: map[string]*connectedPlayer{strings.ToLower(activePlayer.Username()): activePlayer},
		playerIDs:   map[uuid.UUID]*connectedPlayer{activePlayer.ID(): activePlayer},
	}

	found := p.unregisterConnection(oldPlayer)
	if found {
		t.Fatal("expected stale player unregister to return false")
	}
	if got := p.playerIDs[id]; got != activePlayer {
		t.Fatal("active player id mapping was removed or replaced")
	}
	if got := p.playerNames[strings.ToLower(activePlayer.Username())]; got != activePlayer {
		t.Fatal("active player name mapping was removed or replaced")
	}
}

func TestUnregisterConnection_RemovesExactPlayerInstance(t *testing.T) {
	id := uuid.New()
	player := &connectedPlayer{profile: &profile.GameProfile{ID: id, Name: "Notch"}}

	p := &Proxy{
		playerNames: map[string]*connectedPlayer{strings.ToLower(player.Username()): player},
		playerIDs:   map[uuid.UUID]*connectedPlayer{player.ID(): player},
	}

	found := p.unregisterConnection(player)
	if !found {
		t.Fatal("expected unregister of exact player instance to return true")
	}
	if _, ok := p.playerIDs[id]; ok {
		t.Fatal("expected player id mapping to be removed")
	}
	if _, ok := p.playerNames[strings.ToLower(player.Username())]; ok {
		t.Fatal("expected player name mapping to be removed")
	}
}
