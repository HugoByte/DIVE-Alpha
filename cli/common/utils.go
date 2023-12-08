package common

import (
	"os/exec"
	"runtime"
	"strings"

	"github.com/kurtosis-tech/kurtosis/api/golang/core/kurtosis_core_rpc_api_bindings"
	"github.com/kurtosis-tech/kurtosis/api/golang/core/lib/starlark_run_config"
	"github.com/kurtosis-tech/stacktrace"
)

func ValidateArgs(args []string) error {
	if len(args) != 0 {

		return ErrInvalidCommand

	}
	return nil
}

func WriteServiceResponseData(serviceName string, data DiveServiceResponse, cliContext *Cli) error {
	var jsonDataFromFile = Services{}
	err := cliContext.FileHandler().ReadJson("services.json", &jsonDataFromFile)

	if err != nil {
		return WrapMessageToErrorf(err, "Failed To Read %s", "services.json")
	}

	_, ok := jsonDataFromFile[serviceName]
	if !ok {
		jsonDataFromFile[serviceName] = &data

	}
	err = cliContext.FileHandler().WriteJson("services.json", jsonDataFromFile)
	if err != nil {
		return WrapMessageToErrorf(err, "Failed To Write %s", "services.json")
	}

	return nil
}

func OpenFile(URL string) error {
	var args []string
	switch runtime.GOOS {
	case linuxOSName:
		args = []string{openFileLinuxCommandName, URL}
	case macOSName:
		args = []string{openFileMacCommandName, URL}
	case windowsOSName:
		args = []string{openFileWindowsCommandName, openFileWindowsCommandFirstArgumentDefault, URL}
	default:
		return WrapMessageToError(ErrUnsupportedOS, stacktrace.NewError("Unsupported operating system").Error())
	}
	command := exec.Command(args[0], args[1:]...)
	if err := command.Start(); err != nil {
		return WrapMessageToError(ErrOpenFile, stacktrace.Propagate(err, "An error occurred while opening '%v'", URL).Error())
	}
	return nil
}

func LoadConfig(cliContext *Cli, config ConfigLoader, filePath string) error {
	if filePath == "" {
		config.LoadDefaultConfig()
	} else {
		err := config.LoadConfigFromFile(cliContext, filePath)
		if err != nil {
			return err
		}
	}
	return nil
}

func GetStarlarkRunConfig(params string, relativePathToMainFile string, mainFunctionName string) *starlark_run_config.StarlarkRunConfig {

	starlarkConfig := &starlark_run_config.StarlarkRunConfig{
		RelativePathToMainFile:   relativePathToMainFile,
		MainFunctionName:         mainFunctionName,
		DryRun:                   DiveDryRun,
		SerializedParams:         params,
		Parallelism:              DiveDefaultParallelism,
		ExperimentalFeatureFlags: []kurtosis_core_rpc_api_bindings.KurtosisFeatureFlag{},
	}
	return starlarkConfig
}

func GetSerializedData(cliContext *Cli, response chan *kurtosis_core_rpc_api_bindings.StarlarkRunResponseLine) (string, map[string]string, map[string]bool, error) {

	var serializedOutputObj string
	services := map[string]string{}

	skippedInstruction := map[string]bool{}
	for executionResponse := range response {

		if strings.Contains(executionResponse.GetInstructionResult().GetSerializedInstructionResult(), "added with service") {
			res1 := strings.Split(executionResponse.GetInstructionResult().GetSerializedInstructionResult(), " ")
			serviceName := res1[1][1 : len(res1[1])-1]
			serviceUUID := res1[len(res1)-1][1 : len(res1[len(res1)-1])-1]
			services[serviceName] = serviceUUID
		}

		cliContext.log.Info(executionResponse.String())

		if executionResponse.GetInstruction().GetIsSkipped() {
			skippedInstruction[executionResponse.GetInstruction().GetExecutableInstruction()] = executionResponse.GetInstruction().GetIsSkipped()
			break
		}

		if executionResponse.GetError() != nil {

			return "", services, nil, WrapMessageToError(ErrStarlarkResponse, executionResponse.GetError().String())

		}

		runFinishedEvent := executionResponse.GetRunFinishedEvent()

		if runFinishedEvent != nil {

			if runFinishedEvent.GetIsRunSuccessful() {
				serializedOutputObj = runFinishedEvent.GetSerializedOutput()

			} else {
				return "", services, nil, WrapMessageToError(ErrStarlarkResponse, executionResponse.GetError().String())
			}

		} else {
			cliContext.spinner.SetColor("blue")
			if executionResponse.GetProgressInfo() != nil {

				cliContext.spinner.SetSuffixMessage(strings.ReplaceAll(executionResponse.GetProgressInfo().String(), "current_step_info:", " "), "fgGreen")

			}
		}

	}

	return serializedOutputObj, services, skippedInstruction, nil
}