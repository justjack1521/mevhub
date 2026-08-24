package rpc

import (
	"context"
	"github.com/justjack1521/mevium/pkg/genproto/protomulti"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// GameCatchUp is retired. Clients no longer repair a lost notification by asking
// for it back — the game restates its whole state on every transition, on
// reconnect, and on a heartbeat, so a client that missed something converges on
// the next restatement without asking. See action.GameStateSyncChange.
//
// The method survives only because it is on the generated
// MeviusMultiServiceServer interface, which MultiGrpcServer satisfies without
// embedding the Unimplemented server. Removing it would not compile. It can go
// for good once the RPC is dropped from the mevium proto.
func (g MultiGrpcServer) GameCatchUp(ctx context.Context, request *protomulti.GameCatchUpRequest) (*protomulti.GameCatchUpResponse, error) {
	return nil, status.Error(codes.Unimplemented, "game catch up is retired; the game restates its state periodically")
}
