package types

import (
	"context"
	"fmt"

	"github.com/hugobyte/dive/common"
	"github.com/kurtosis-tech/kurtosis/api/golang/core/kurtosis_core_rpc_api_bindings"
	"github.com/kurtosis-tech/kurtosis/api/golang/core/lib/enclaves"
	"github.com/spf13/cobra"
)

func NewHardhatCmd(ctx context.Context, kurtosisEnclaveContext *enclaves.EnclaveContext) *cobra.Command {

	var ethCmd = &cobra.Command{
		Use:   "hardhat",
		Short: "Runs Hardhat Node",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {

			fmt.Println(RunHardhatNode(ctx, kurtosisEnclaveContext).EncodeToString())
		},
	}

	return ethCmd

}

func RunHardhatNode(ctx context.Context, kurtosisEnclaveContext *enclaves.EnclaveContext) *common.DiveserviceResponse {

	data, _, err := kurtosisEnclaveContext.RunStarlarkPackage(ctx, "../", "services/evm/eth/eth.star", "start_eth_node_serivce", `{"args":{},"node_type":"hardhat"}`, false, 4, []kurtosis_core_rpc_api_bindings.KurtosisFeatureFlag{})

	if err != nil {
		fmt.Println(err)
	}

	responseData := common.GetSerializedData(data)

	ethResponseData := &common.DiveserviceResponse{}

	result, err := ethResponseData.Decode([]byte(responseData))

	if err != nil {
		fmt.Println(err)
	}

	return result

}
