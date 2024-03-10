package ante_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	ibcclienttypes "github.com/cosmos/ibc-go/v4/modules/core/02-client/types"
	ibcconnectiontypes "github.com/cosmos/ibc-go/v4/modules/core/03-connection/types"
	ibcchanneltypes "github.com/cosmos/ibc-go/v4/modules/core/04-channel/types"
	appconsumer "github.com/onomyprotocol/multiverse/app/consumer-democracy"
	ante "github.com/onomyprotocol/multiverse/app/consumer-democracy/consumer-ante"
	"github.com/stretchr/testify/require"
	"github.com/tendermint/spm/cosmoscmd"
)

type consumerKeeper struct {
	channelExists bool
}

func (k consumerKeeper) GetProviderChannel(_ sdk.Context) (string, bool) {
	return "", k.channelExists
}

func noOpAnteDecorator() sdk.AnteHandler {
	return func(ctx sdk.Context, _ sdk.Tx, _ bool) (sdk.Context, error) {
		return ctx, nil
	}
}

func TestMsgFilterDecorator(t *testing.T) {
	txCfg := cosmoscmd.MakeEncodingConfig(appconsumer.ModuleBasics).TxConfig

	testCases := []struct {
		name           string
		ctx            sdk.Context
		consumerKeeper ante.ConsumerKeeper
		msgs           []sdk.Msg
		expectErr      bool
	}{
		{
			name:           "valid tx pre-CCV",
			ctx:            sdk.Context{},
			consumerKeeper: consumerKeeper{channelExists: false},
			msgs: []sdk.Msg{
				&ibcclienttypes.MsgUpdateClient{},
			},
			expectErr: false,
		},
		{
			name:           "invalid tx pre-CCV",
			ctx:            sdk.Context{},
			consumerKeeper: consumerKeeper{channelExists: false},
			msgs: []sdk.Msg{
				&banktypes.MsgSend{},
			},
			expectErr: true,
		},
		{
			name:           "valid tx post-CCV",
			ctx:            sdk.Context{},
			consumerKeeper: consumerKeeper{channelExists: true},
			msgs: []sdk.Msg{
				&banktypes.MsgSend{},
			},
			expectErr: false,
		},
		{
			name:           "invalid pre-CCV MsgConnectionOpenInit",
			ctx:            sdk.Context{},
			consumerKeeper: consumerKeeper{channelExists: false},
			msgs: []sdk.Msg{
				&ibcconnectiontypes.MsgConnectionOpenInit{ClientId: "07-tendermint-1"},
			},
			expectErr: true,
		},
		{
			name:           "valid pre-CCV MsgConnectionOpenInit",
			ctx:            sdk.Context{},
			consumerKeeper: consumerKeeper{channelExists: false},
			msgs: []sdk.Msg{
				&ibcconnectiontypes.MsgConnectionOpenInit{ClientId: "07-tendermint-0"},
			},
			expectErr: false,
		},
		{
			name:           "invalid pre-CCV MsgChannelOpenInit",
			ctx:            sdk.Context{},
			consumerKeeper: consumerKeeper{channelExists: false},
			msgs: []sdk.Msg{
				&ibcchanneltypes.MsgChannelOpenInit{PortId: "transfer", Channel: ibcchanneltypes.Channel{Counterparty: ibcchanneltypes.Counterparty{PortId: "transfer"}, ConnectionHops: []string{"connection-0"}}},
			},
			expectErr: true,
		},
		{
			name:           "invalid pre-CCV MsgChannelOpenInit",
			ctx:            sdk.Context{},
			consumerKeeper: consumerKeeper{channelExists: false},
			msgs: []sdk.Msg{
				&ibcchanneltypes.MsgChannelOpenInit{PortId: "consumer", Channel: ibcchanneltypes.Channel{Counterparty: ibcchanneltypes.Counterparty{PortId: "provider"}, ConnectionHops: []string{"connection-1"}}},
			},
			expectErr: true,
		},
		{
			name:           "valid pre-CCV MsgChannelOpenInit",
			ctx:            sdk.Context{},
			consumerKeeper: consumerKeeper{channelExists: false},
			msgs: []sdk.Msg{
				&ibcchanneltypes.MsgChannelOpenInit{PortId: "consumer", Channel: ibcchanneltypes.Channel{Counterparty: ibcchanneltypes.Counterparty{PortId: "provider"}, ConnectionHops: []string{"connection-0"}}},
			},
			expectErr: false,
		},
	}

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			handler := ante.NewMsgFilterDecorator(tc.consumerKeeper)

			txBuilder := txCfg.NewTxBuilder()
			require.NoError(t, txBuilder.SetMsgs(tc.msgs...))

			_, err := handler.AnteHandle(tc.ctx, txBuilder.GetTx(), false, noOpAnteDecorator())
			if tc.expectErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
