package common

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"runtime"

	"github.com/kurtosis-tech/kurtosis/api/golang/core/kurtosis_core_rpc_api_bindings"
	"github.com/kurtosis-tech/kurtosis/api/golang/core/lib/enclaves"
	"github.com/kurtosis-tech/kurtosis/api/golang/engine/lib/kurtosis_context"
	"github.com/kurtosis-tech/stacktrace"
	"github.com/sirupsen/logrus"
)

const (
	linuxOSName   = "linux"
	macOSName     = "darwin"
	windowsOSName = "windows"

	openFileLinuxCommandName   = "xdg-open"
	openFileMacCommandName     = "open"
	openFileWindowsCommandName = "rundll32"

	openFileWindowsCommandFirstArgumentDefault = "url.dll,FileProtocolHandler"
)

type DiveserviceResponse struct {
	ServiceName     string `json:"service_name"`
	PublicEndpoint  string `json:"endpoint_public"`
	PrivateEndpoint string `json:"endpoint"`
	KeyPassword     string `json:"keypassword"`
	KeystorePath    string `json:"keystore_path"`
	Network         string `json:"network"`
	NetworkName     string `json:"network_name"`
	NetworkId       string `json:"nid"`
}

func (dive *DiveserviceResponse) Decode(responseData []byte) (*DiveserviceResponse, error) {

	err := json.Unmarshal(responseData, &dive)
	if err != nil {
		return nil, err
	}
	return dive, nil
}
func (dive *DiveserviceResponse) EncodeToString() (string, error) {
	encodedBytes, err := json.Marshal(dive)
	if err != nil {
		return "", nil
	}

	return string(encodedBytes), nil
}

func GetSerializedData(response chan *kurtosis_core_rpc_api_bindings.StarlarkRunResponseLine) string {

	var serializedOutputObj string
	for executionResponseLine := range response {
		fmt.Println(executionResponseLine)
		runFinishedEvent := executionResponseLine.GetRunFinishedEvent()
		if runFinishedEvent == nil {
			logrus.Info("Execution in progress...")
		} else {
			logrus.Info("Execution finished successfully")
			if runFinishedEvent.GetIsRunSuccessful() {
				serializedOutputObj = runFinishedEvent.GetSerializedOutput()
			} else {
				panic("Starlark run failed")
			}
		}
	}

	return serializedOutputObj

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
		return stacktrace.NewError("Unsupported operating system")
	}
	command := exec.Command(args[0], args[1:]...)
	if err := command.Start(); err != nil {
		return stacktrace.Propagate(err, "An error occurred while opening '%v'", URL)
	}
	return nil
}

type DiveContext struct {
	Ctx             context.Context
	KurtosisContext *kurtosis_context.KurtosisContext
}

func NewDiveContext() *DiveContext {

	ctx := context.Background()

	kurtosisContext, err := kurtosis_context.NewKurtosisContextFromLocalEngine()
	if err != nil {
		panic(err)

	}
	return &DiveContext{Ctx: ctx, KurtosisContext: kurtosisContext}
}

func (diveContext *DiveContext) GetEnclaveContext() (*enclaves.EnclaveContext, error) {

	_, err := diveContext.KurtosisContext.GetEnclave(diveContext.Ctx, DiveEnclave)
	if err != nil {
		enclaveCtx, err := diveContext.KurtosisContext.CreateEnclave(diveContext.Ctx, DiveEnclave, false)
		if err != nil {
			return nil, err

		}
		return enclaveCtx, nil
	}
	enclaveCtx, err := diveContext.KurtosisContext.GetEnclaveContext(diveContext.Ctx, DiveEnclave)

	if err != nil {
		return nil, err
	}
	return enclaveCtx, nil
}

func ReadConfigFile(filePath string) ([]byte, error) {

	file, err := os.ReadFile(filePath)

	if err != nil {
		return nil, err
	}

	return file, nil
}
func WriteToFile(data string) {
	file, err := os.Create("bridge.json")
	if err != nil {
		return
	}
	defer file.Close()

	file.WriteString(data)
}
