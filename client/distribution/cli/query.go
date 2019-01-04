package cli

import (
	"fmt"

	"github.com/irisnet/irishub/codec"
	sdk "github.com/irisnet/irishub/types"
	"github.com/irisnet/irishub/modules/distribution"
	"github.com/irisnet/irishub/modules/distribution/types"
	"github.com/irisnet/irishub/client/context"
	distributionclient "github.com/irisnet/irishub/client/distribution"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const NULL = "null"

// GetDelegationDistInfo returns the delegation distribution information of a given delegation
func GetDelegationDistInfo(storeName string, cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "delegation-distr-info",
		Short:   "Query delegation distribution information",
		Example: "iriscli distribution delegation-distr-info --address-delegator=<delegator address> --address-validator=<validator address>",
		RunE: func(cmd *cobra.Command, args []string) error {

			valAddr, err := sdk.ValAddressFromBech32(viper.GetString(FlagAddressValidator))
			if err != nil {
				return err
			}
			delAddr, err := sdk.AccAddressFromBech32(viper.GetString(FlagAddressDelegator))
			if err != nil {
				return err
			}

			key := distribution.GetDelegationDistInfoKey(delAddr, valAddr)

			cliCtx := context.NewCLIContext().WithCodec(cdc)

			res, err := cliCtx.QueryStore(key, storeName)
			if err != nil {
				return err
			}
			if len(res) == 0 {
				fmt.Println(NULL)
				return nil
			}
			var ddi types.DelegationDistInfo
			err = cdc.UnmarshalBinaryLengthPrefixed(res, &ddi)
			if err != nil {
				return err
			}

			output, err := codec.MarshalJSONIndent(cdc, ddi)
			if err != nil {
				return err
			}

			fmt.Println(string(output))
			return nil
		},
	}
	cmd.Flags().String(FlagAddressDelegator, "", "bech address of the delegator")
	cmd.Flags().String(FlagAddressValidator, "", "bech address of the validator")
	cmd.MarkFlagRequired(FlagAddressDelegator)
	cmd.MarkFlagRequired(FlagAddressValidator)
	return cmd
}

// GetAllDelegationDistInfo returns all delegation distribution information of a given delegator
func GetAllDelegationDistInfo(storeName string, cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "delegator-distr-info",
		Short:   "Query delegator distribution information",
		Example: "iriscli distribution delegator-distr-info <delegator address>",
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			addrString := args[0]

			delAddr, err := sdk.AccAddressFromBech32(addrString)
			if err != nil {
				return err
			}

			key := distribution.GetDelegationDistInfosKey(delAddr)

			cliCtx := context.NewCLIContext().WithCodec(cdc)

			resKVs, err := cliCtx.QuerySubspace(key, storeName)
			if err != nil {
				return err
			}
			if len(resKVs) == 0 {
				fmt.Println(NULL)
				return nil
			}
			var ddiList []types.DelegationDistInfo
			for _, kv := range resKVs {
				var ddi types.DelegationDistInfo
				err = cdc.UnmarshalBinaryLengthPrefixed(kv.Value, &ddi)
				if err != nil {
					return err
				}
				ddiList = append(ddiList, ddi)
			}

			output, err := codec.MarshalJSONIndent(cdc, ddiList)
			if err != nil {
				return err
			}

			fmt.Println(string(output))
			return nil
		},
	}
	return cmd
}

// GetValidatorDistInfo returns the validator distribution information of a given validator
func GetValidatorDistInfo(storeName string, cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "validator-distr-info",
		Short:   "Query validator distribution information",
		Example: "iriscli distribution validator-distr-info <validator address>",
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {

			addrString := args[0]

			valAddr, err := sdk.ValAddressFromBech32(addrString)
			if err != nil {
				return err
			}

			key := distribution.GetValidatorDistInfoKey(valAddr)

			cliCtx := context.NewCLIContext().WithCodec(cdc)

			res, err := cliCtx.QueryStore(key, storeName)
			if err != nil {
				return err
			}
			if len(res) == 0 {
				fmt.Println(NULL)
				return nil
			}

			var vdi types.ValidatorDistInfo
			err = cdc.UnmarshalBinaryLengthPrefixed(res, &vdi)
			if err != nil {
				return err
			}

			vdiOutput := distributionclient.ConvertToValidatorDistInfoOutput(cliCtx, vdi)

			output, err := codec.MarshalJSONIndent(cdc, vdiOutput)
			if err != nil {
				return err
			}

			fmt.Println(string(output))
			return nil
		},
	}
	return cmd
}
