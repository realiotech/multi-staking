package cli

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	multistakingtypes "github.com/realio-tech/multi-staking-module/x/multi-staking/types"
	"github.com/spf13/cobra"
	flag "github.com/spf13/pflag"

	"cosmossdk.io/core/address"
	errorsmod "cosmossdk.io/errors"
	"cosmossdk.io/math"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/version"
	"github.com/cosmos/cosmos-sdk/x/staking/client/cli"
	"github.com/cosmos/cosmos-sdk/x/staking/types"
)

// NewCreateValidatorCmd returns a CLI command handler for creating a MsgCreateValidator transaction.
func NewCreateEVMValidatorCmd(ac address.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create-evm-validator [path/to/validator.json]",
		Short: "create new validator using ERC20 token initialized with a self-delegation to it",
		Args:  cobra.ExactArgs(1),
		Long:  `Create a new validator using ERC20 token initialized with a self-delegation by submitting a JSON file with the new validator details.`,
		Example: strings.TrimSpace(
			fmt.Sprintf(`
$ %s tx staking create-evm-validator path/to/validator.json --from keyname

Where validator.json contains:

{
	"pubkey": {"@type":"/cosmos.crypto.ed25519.PubKey","key":"oWg2ISpLF405Jcm2vXV+2v4fnjodh6aafuIdeoW+rUw="},
	"contract-address": "0x1234567890123456789012345678901234567890",
	"amount": "1000000",
	"moniker": "myvalidator",
	"identity": "optional identity signature (ex. UPort or Keybase)",
	"website": "validator's (optional) website",
	"security": "validator's (optional) security contact email",
	"details": "validator's (optional) details",
	"commission-rate": "0.1",
	"commission-max-rate": "0.2",
	"commission-max-change-rate": "0.01",
	"min-self-delegation": "1"
}

where we can get the pubkey using "%s tendermint show-validator"
`, version.AppName, version.AppName)),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			txf, err := tx.NewFactoryCLI(clientCtx, cmd.Flags())
			if err != nil {
				return err
			}

			validator, err := parseAndValidateEVMValidatorJSON(clientCtx.Codec, args[0])
			if err != nil {
				return err
			}

			txf, msg, err := newBuildCreateEVMValidatorMsg(clientCtx, txf, cmd.Flags(), validator, ac)
			if err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxWithFactory(clientCtx, txf, msg)
		},
	}

	cmd.Flags().String(cli.FlagIP, "", fmt.Sprintf("The node's public IP. It takes effect only when used in combination with --%s", flags.FlagGenerateOnly))
	cmd.Flags().String(cli.FlagNodeID, "", "The node's ID")
	flags.AddTxFlagsToCmd(cmd)

	_ = cmd.MarkFlagRequired(flags.FlagFrom)

	return cmd
}

// NewDelegateEVMCmd returns a CLI command handler for creating a MsgDelegate transaction.
func NewDelegateEVMCmd(valAddrCodec, ac address.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delegate-evm [validator-addr] [contract-addr] [amount]",
		Args:  cobra.ExactArgs(3),
		Short: "Delegate liquid tokens to a validator",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Delegate an amount of ERC20 token to a validator from your wallet.

Example:
$ %s tx staking delegate-evm cosmosvalopers1l2rsakp388kuv9k8qzq6lrm9taddae7fpx59wm 0x1234567890123456789012345678901234567890 1000 --from mykey
`,
				version.AppName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			contractAddr := args[1]
			if !common.IsHexAddress(contractAddr) {
				return fmt.Errorf("invalid contract address")
			}

			amount, ok := math.NewIntFromString(args[2])
			if !ok {
				return fmt.Errorf("invalid amount")
			}

			delAddr, err := ac.BytesToString(clientCtx.GetFromAddress())
			if err != nil {
				return err
			}

			_, err = valAddrCodec.StringToBytes(args[0])
			if err != nil {
				return err
			}

			msg := &multistakingtypes.MsgDelegateEVM{
				DelegatorAddress: delAddr,
				ValidatorAddress: args[0],
				ContractAddress:  contractAddr,
				Amount:           amount,
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

// NewRedelegateEVMCmd returns a CLI command handler for creating a MsgBeginRedelegateEVM transaction.
func NewRedelegateEVMCmd(valAddrCodec, ac address.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "redelegate-evm [src-validator-addr] [dst-validator-addr] [contract-addr] [amount]",
		Short: "Redelegate illiquid tokens from one validator to another",
		Args:  cobra.ExactArgs(4),
		Long: strings.TrimSpace(
			fmt.Sprintf(`Redelegate an amount of illiquid staking tokens from one validator to another.

Example:
$ %s tx staking redelegate cosmosvalopers1gghjut3ccd8ay0zduzj64hwre2fxs9ldmqhffj cosmosvalopers1l2rsakp388kuv9k8qzq6lrm9taddae7fpx59wm 0x1234567890123456789012345678901234567890 100 --from mykey
`,
				version.AppName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}
			delAddr, err := ac.BytesToString(clientCtx.GetFromAddress())
			if err != nil {
				return err
			}

			_, err = valAddrCodec.StringToBytes(args[0])
			if err != nil {
				return err
			}

			_, err = valAddrCodec.StringToBytes(args[1])
			if err != nil {
				return err
			}

			contractAddr := args[2]
			if !common.IsHexAddress(contractAddr) {
				return fmt.Errorf("invalid contract address")
			}

			amount, ok := math.NewIntFromString(args[3])
			if !ok {
				return fmt.Errorf("invalid amount")
			}

			msg := &multistakingtypes.MsgBeginRedelegateEVM{
				DelegatorAddress:    delAddr,
				ValidatorSrcAddress: args[0],
				ValidatorDstAddress: args[1],
				ContractAddress:     contractAddr,
				Amount:              amount,
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

// NewUnbondEVMCmd returns a CLI command handler for creating a MsgUndelegate transaction.
func NewUnbondEVMCmd(valAddrCodec, ac address.Codec) *cobra.Command {
	bech32PrefixValAddr := sdk.GetConfig().GetBech32ValidatorAddrPrefix()

	cmd := &cobra.Command{
		Use:   "unbond-evm [validator-addr] [contract-addr] [amount]",
		Short: "Unbond erc20 shares from a validator",
		Args:  cobra.ExactArgs(3),
		Long: strings.TrimSpace(
			fmt.Sprintf(`Unbond an amount of bonded shares from a validator.

Example:
$ %s tx staking unbond %s1gghjut3ccd8ay0zduzj64hwre2fxs9ldmqhffj 0x1234567890123456789012345678901234567890 100 --from mykey
`,
				version.AppName, bech32PrefixValAddr,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			delAddr, err := ac.BytesToString(clientCtx.GetFromAddress())
			if err != nil {
				return err
			}
			_, err = valAddrCodec.StringToBytes(args[0])
			if err != nil {
				return err
			}

			contractAddr := args[1]
			if !common.IsHexAddress(contractAddr) {
				return fmt.Errorf("invalid contract address")
			}

			amount, ok := math.NewIntFromString(args[2])
			if !ok {
				return fmt.Errorf("invalid amount")
			}

			msg := &multistakingtypes.MsgUndelegateEVM{
				DelegatorAddress: delAddr,
				ValidatorAddress: args[0],
				ContractAddress:  contractAddr,
				Amount:           amount,
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

// NewCancelUnbondingEVMDelegation returns a CLI command handler for creating a MsgCancelUnbondingEVMDelegation transaction.
func NewCancelUnbondingEVMDelegation(valAddrCodec, ac address.Codec) *cobra.Command {
	bech32PrefixValAddr := sdk.GetConfig().GetBech32ValidatorAddrPrefix()

	cmd := &cobra.Command{
		Use:   "cancel-unbond-evm [validator-addr] [contract-addr] [amount] [creation-height]",
		Short: "Cancel unbonding delegation and delegate back to the validator",
		Args:  cobra.ExactArgs(3),
		Long: strings.TrimSpace(
			fmt.Sprintf(`Cancel Unbonding Delegation and delegate back to the validator.

Example:
$ %s tx staking cancel-unbond %s1gghjut3ccd8ay0zduzj64hwre2fxs9ldmqhffj 0x1234567890123456789012345678901234567890 100 2 --from mykey
`,
				version.AppName, bech32PrefixValAddr,
			),
		),
		Example: fmt.Sprintf(`$ %s tx staking cancel-unbond %s1gghjut3ccd8ay0zduzj64hwre2fxs9ldmqhffj 100stake 2 --from mykey`,
			version.AppName, bech32PrefixValAddr),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}
			delAddr, err := ac.BytesToString(clientCtx.GetFromAddress())
			if err != nil {
				return err
			}

			_, err = valAddrCodec.StringToBytes(args[0])
			if err != nil {
				return err
			}

			contractAddr := args[1]
			if !common.IsHexAddress(contractAddr) {
				return fmt.Errorf("invalid contract address")
			}

			amount, ok := math.NewIntFromString(args[2])
			if !ok {
				return fmt.Errorf("invalid amount")
			}

			creationHeight, err := strconv.ParseInt(args[2], 10, 64)
			if err != nil {
				return errorsmod.Wrap(fmt.Errorf("invalid height: %d", creationHeight), "invalid height")
			}

			msg := &multistakingtypes.MsgCancelUnbondingEVMDelegation{
				DelegatorAddress: delAddr,
				ValidatorAddress: args[0],
				ContractAddress:  contractAddr,
				Amount:           amount,
				CreationHeight:   creationHeight,
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

func newBuildCreateEVMValidatorMsg(clientCtx client.Context, txf tx.Factory, fs *flag.FlagSet, val evmValidator, valAc address.Codec) (tx.Factory, *multistakingtypes.MsgCreateEVMValidator, error) {
	valAddr := clientCtx.GetFromAddress()

	description := types.NewDescription(
		val.Moniker,
		val.Identity,
		val.Website,
		val.Security,
		val.Details,
	)

	valStr, err := valAc.BytesToString(sdk.ValAddress(valAddr))
	if err != nil {
		return txf, nil, err
	}
	msg, err := multistakingtypes.NewMsgCreateEVMValidator(
		valStr, val.PubKey, val.ContractAddress, val.Amount, description, val.CommissionRates, val.MinSelfDelegation,
	)
	if err != nil {
		return txf, nil, err
	}
	if err := msg.Validate(valAc); err != nil {
		return txf, nil, err
	}

	genOnly, _ := fs.GetBool(flags.FlagGenerateOnly)
	if genOnly {
		ip, _ := fs.GetString(cli.FlagIP)
		p2pPort, _ := fs.GetUint(cli.FlagP2PPort)
		nodeID, _ := fs.GetString(cli.FlagNodeID)

		if nodeID != "" && ip != "" && p2pPort > 0 {
			txf = txf.WithMemo(fmt.Sprintf("%s@%s:%d", nodeID, ip, p2pPort))
		}
	}

	return txf, msg, nil
}
