package ante

import (
	"fmt"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	ibcconnectiontypes "github.com/cosmos/ibc-go/v4/modules/core/03-connection/types"
	ibcchanneltypes "github.com/cosmos/ibc-go/v4/modules/core/04-channel/types"
)

type (
	// ConsumerKeeper defines the interface required by a consumer module keeper.
	ConsumerKeeper interface {
		GetProviderChannel(ctx sdk.Context) (string, bool)
	}

	// MsgFilterDecorator defines an AnteHandler decorator that enables message
	// filtering based on certain criteria.
	MsgFilterDecorator struct {
		ConsumerKeeper ConsumerKeeper
	}
)

func NewMsgFilterDecorator(k ConsumerKeeper) MsgFilterDecorator {
	return MsgFilterDecorator{
		ConsumerKeeper: k,
	}
}

func (mfd MsgFilterDecorator) AnteHandle(ctx sdk.Context, tx sdk.Tx, simulate bool, next sdk.AnteHandler) (newCtx sdk.Context, err error) {
	currHeight := ctx.BlockHeight()

	// If the CCV channel has not yet been established, then we must only allow certain
	// message types.
	if _, ok := mfd.ConsumerKeeper.GetProviderChannel(ctx); !ok {
		if !hasValidMsgsPreCCV(tx.GetMsgs()) {
			return ctx, fmt.Errorf("tx contains unsupported message types at height %d", currHeight)
		}
	}

	return next(ctx, tx, simulate)
}

func hasValidMsgsPreCCV(msgs []sdk.Msg) bool {
	for _, msg := range msgs {
		msgType := sdk.MsgTypeURL(msg)

		// We want to make sure that the first connection and channel are the correct ICS opening,
		// so that the IBC denom can be deterministic. 07-tendermint-0 is automatically created
		// for ICS setup, we need to make sure that connection-0 is to this client, and that
		// channel-0 uses the right ports and connection-0. ICS automatically creates the
		// connection-1 for transfer port. The port id, channel id, and denom name are used for
		// hash for IBC denom. Orderings and other fields are checked by preexisting checks in ICS
		switch msg := msg.(type) {
		case *ibcconnectiontypes.MsgConnectionOpenInit:
			if msg.ClientId != "07-tendermint-0" {
				return false
			}
		case *ibcchanneltypes.MsgChannelOpenInit:
			if (msg.PortId != "consumer") || (msg.Channel.Counterparty.PortId != "provider") || (len(msg.Channel.ConnectionHops) != 1) || (msg.Channel.ConnectionHops[0] != "connection-0") {
				return false
			}
		default:
		}

		// Only accept IBC messages prior to the CCV channel being established.
		// Note, rather than listing out all possible IBC message types, we assume
		// all IBC message types have a correct and canonical prefix -- /ibc.*
		if !strings.HasPrefix(msgType, "/ibc.") {
			return false
		}
	}

	return true
}
