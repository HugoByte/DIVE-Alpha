package dive_test

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"testing"

	dive "github.com/HugoByte/DIVE/test/functional"
	"github.com/hugobyte/dive/cli/common"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
)

// To Print cli output to console
type testWriter struct {
	buffer bytes.Buffer
}

func TestCLIApp(t *testing.T) {
	gomega.RegisterFailHandler(ginkgo.Fail)
	ginkgo.RunSpecs(t, "DIVE CLI App Suite")
}

func (w *testWriter) Write(p []byte) (n int, err error) {
	w.buffer.Write(p)
	os.Stdout.Write(p)
	return len(p), nil
}

var _ = ginkgo.Describe("DIVE CLI App", func() {
	var cmd *exec.Cmd
	var stdout bytes.Buffer

	// run clean before each test
	ginkgo.BeforeEach(func() {
		cmd = dive.GetBinaryCommand()
		cmd.Stdout = &testWriter{}
		cmd.Stderr = &testWriter{}
	})

	ginkgo.AfterEach(func() {
		dive.Clean()
	})

	ginkgo.Describe("Smoke Tests", func() {
		ginkgo.It("should display the correct version", func() {
			cmd.Args = append(cmd.Args, "version")
			cmd.Stdout = &stdout
			err := cmd.Run()
			fmt.Println(stdout.String())
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
			latestVersion := common.GetLatestVersion()
			gomega.Expect(stdout.String()).To(gomega.ContainSubstring(latestVersion))
		})

		ginkgo.It("should open twitter page on browser", func() {
			cmd.Args = append(cmd.Args, "twitter")
			err := cmd.Run()
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
		})

		ginkgo.It("should open Youtube Tutorial Channel on browser", func() {
			cmd.Args = append(cmd.Args, "tutorial")
			err := cmd.Run()
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
		})

		ginkgo.It("should open Discord channel on browser", func() {
			cmd.Args = append(cmd.Args, "discord")
			err := cmd.Run()
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
		})

		ginkgo.It("should start bridge between icon and eth", func() {
			dive.Clean()
			cmd.Args = append(cmd.Args, "bridge", "btp", "--chainA", "icon", "--chainB", "eth")
			err := cmd.Run()
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
		})

		ginkgo.It("should start bridge between icon and hardhat but with icon bridge set to true", func() {
			cmd.Args = append(cmd.Args, "bridge", "btp", "--chainA", "icon", "--chainB", "hardhat", "--bmvbridge")
			err := cmd.Run()
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
		})

		ginkgo.It("should start bridge between icon and icon", func() {
			cmd.Args = append(cmd.Args, "bridge", "btp", "--chainA", "icon", "--chainB", "icon")
			err := cmd.Run()
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
		})

		ginkgo.It("should start bridge between archway and archway using ibc", func() {
			cmd.Args = append(cmd.Args, "bridge", "ibc", "--chainA", "archway", "--chainB", "archway")
			err := cmd.Run()
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
		})

	})

	ginkgo.Describe("Bridge command Test", func() {
		ginkgo.It("should start bridge between icon and eth but with icon bridge set to true", func() {
			cmd.Args = append(cmd.Args, "bridge", "btp", "--chainA", "icon", "--chainB", "eth", "--bmvbridge")
			err := cmd.Run()
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
		})

		ginkgo.It("should start bridge between icon and eth with verbose flag enabled", func() {
			cmd.Args = append(cmd.Args, "bridge", "btp", "--chainA", "icon", "--chainB", "eth", "--verbose")
			err := cmd.Run()
			gomega.Expect(err).NotTo(gomega.HaveOccurred())

		})

		ginkgo.It("should start bridge between icon and eth but with icon bridge set to true with verbose flag enabled", func() {
			cmd.Args = append(cmd.Args, "bridge", "btp", "--chainA", "icon", "--chainB", "eth", "--bmvbridge", "--verbose")
			err := cmd.Run()
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
		})

		ginkgo.It("should start bridge between icon and hardhat", func() {
			cmd.Args = append(cmd.Args, "bridge", "btp", "--chainA", "icon", "--chainB", "hardhat")
			err := cmd.Run()
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
		})

		ginkgo.It("should start bridge between icon and hardhat with verbose flag enabled", func() {
			cmd.Args = append(cmd.Args, "bridge", "btp", "--chainA", "icon", "--chainB", "hardhat", "--verbose")
			err := cmd.Run()
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
		})

		ginkgo.It("should start bridge between icon and hardhat but with icon bridge set to true with verbose flag enabled", func() {
			cmd.Args = append(cmd.Args, "bridge", "btp", "--chainA", "icon", "--chainB", "hardhat", "--bmvbridge", "--verbose")
			err := cmd.Run()
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
		})

		ginkgo.It("should start bridge between icon and icon with verbose flag enabled", func() {
			cmd.Args = append(cmd.Args, "bridge", "btp", "--chainA", "icon", "--chainB", "icon", "--verbose")
			err := cmd.Run()
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
		})

		ginkgo.It("should start bridge between icon and eth by running each chain individually ", func() {
			dive.RunDecentralizedIconNode()
			dive.RunEthNode()
			cmd.Args = append(cmd.Args, "bridge", "btp", "--chainA", "icon", "--chainB", "eth", "--chainAServiceName", "icon-node-0xacbc4e", "--chainBServiceName", "el-1-geth-lighthouse")
			err := cmd.Run()
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
		})

		ginkgo.It("should start bridge between icon and hardhat by running each chain individually ", func() {
			dive.RunDecentralizedIconNode()
			dive.RunHardhatNode()
			cmd.Args = append(cmd.Args, "bridge", "btp", "--chainA", "icon", "--chainB", "hardhat", "--chainAServiceName", "icon-node-0xacbc4e", "--chainBServiceName", "hardhat-node")
			err := cmd.Run()
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
		})

		ginkgo.It("should start bridge between icon and eth by running icon node first and then decentralising it", func() {
			dive.RunIconNode()
			dive.DecentralizeIconNode()
			dive.RunEthNode()
			cmd.Args = append(cmd.Args, "bridge", "btp", "--chainA", "icon", "--chainB", "eth", "--chainAServiceName", "icon-node-0xacbc4e", "--chainBServiceName", "el-1-geth-lighthouse")
			err := cmd.Run()
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
		})

		ginkgo.It("should start bridge between icon and hardhat by running icon node first and then decentralising it", func() {
			dive.RunIconNode()
			dive.DecentralizeIconNode()
			dive.RunHardhatNode()
			cmd.Args = append(cmd.Args, "bridge", "btp", "--chainA", "icon", "--chainB", "hardhat", "--chainAServiceName", "icon-node-0xacbc4e", "--chainBServiceName", "hardhat-node")
			err := cmd.Run()
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
		})

		ginkgo.It("should start bridge between icon and icon by running one custom icon chain", func() {
			dive.RunDecentralizedIconNode()
			dive.RunDecentralizedCustomIconNode1()
			cmd.Args = append(cmd.Args, "bridge", "btp", "--chainA", "icon", "--chainB", "icon", "--chainAServiceName", "icon-node-0xacbc4e", "--chainBServiceName", "icon-node-0x42f1f3")
			err := cmd.Run()
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
		})

		ginkgo.It("should start bridge between icon and running custom icon later decentralising it", func() {
			dive.RunDecentralizedIconNode()
			dive.RunCustomIconNode()
			dive.DecentralizeCustomIconNode()
			cmd.Args = append(cmd.Args, "bridge", "btp", "--chainA", "icon", "--chainB", "icon", "--chainAServiceName", "icon-node-0xacbc4e", "--chainBServiceName", "icon-node-0x42f1f3")
			err := cmd.Run()
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
		})

		ginkgo.It("should start bridge between icon and icon by running one icon chain and later decentralsing it. Running another custom icon chain and then decentralising it", func() {
			dive.RunIconNode()
			dive.DecentralizeIconNode()
			dive.RunCustomIconNode()
			dive.DecentralizeCustomIconNode()
			cmd.Args = append(cmd.Args, "bridge", "btp", "--chainA", "icon", "--chainB", "icon", "--chainAServiceName", "icon-node-0xacbc4e", "--chainBServiceName", "icon-node-0x42f1f3")
			err := cmd.Run()
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
		})

		ginkgo.It("should start bridge between 2 custom icon chains", func() {
			dive.RunDecentralizedCustomIconNode0()
			dive.RunDecentralizedCustomIconNode1()
			cmd.Args = append(cmd.Args, "bridge", "btp", "--chainA", "icon", "--chainB", "icon", "--chainAServiceName", "icon-node-0xacbc4e", "--chainBServiceName", "icon-node-0x42f1f3")
			err := cmd.Run()
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
		})

		ginkgo.It("should start bridge between 2 custom icon chains by running them first and then decentralising it later", func() {
			dive.RunCustomIconNode_0()
			dive.DecentralizeCustomIconNode_0()
			dive.RunCustomIconNode()
			dive.DecentralizeCustomIconNode()
			cmd.Args = append(cmd.Args, "bridge", "btp", "--chainA", "icon", "--chainB", "icon", "--chainAServiceName", "icon-node-0xacbc4e", "--chainBServiceName", "icon-node-0x42f1f3")
			err := cmd.Run()
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
		})

		ginkgo.It("should start bridge between 2 chains when all nodes are running", func() {
			dive.RunDecentralizedIconNode()
			dive.RunEthNode()
			dive.RunHardhatNode()
			cmd.Args = append(cmd.Args, "bridge", "btp", "--chainA", "icon", "--chainB", "eth", "--chainAServiceName", "icon-node-0xacbc4e", "--chainBServiceName", "el-1-geth-lighthouse")
			err := cmd.Run()
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
		})

		ginkgo.It("should handle invalid input for bridge command", func() {
			cmd.Args = append(cmd.Args, "bridge", "btp", "--chainA", "icon", "--chainB", "invalid_input")
			err := cmd.Run()
			gomega.Expect(err).To(gomega.HaveOccurred())
		})

		ginkgo.It("should handle invalid input for bridge command", func() {
			cmd.Args = append(cmd.Args, "bridge", "btp", "--chainA", "invalid_input", "--chainB", "eth")
			err := cmd.Run()
			gomega.Expect(err).To(gomega.HaveOccurred())
		})

		ginkgo.It("should handle invalid input ibc bridge command", func() {
			cmd.Args = append(cmd.Args, "bridge", "ibc", "--chainA", "archway", "--chainB", "invalid")
			err := cmd.Run()
			gomega.Expect(err).To(gomega.HaveOccurred())
		})

		ginkgo.It("should start bridge between archway and archway by running one custom archway chain", func() {
			dive.RunArchwayNode()
			dive.RunCustomArchwayNode1()
			cmd.Args = append(cmd.Args, "bridge", "ibc", "--chainA", "archway", "--chainB", "archway", "--chainAServiceName", "node-service-constantine-3", "--chainBServiceName", "node-service-archway-node-1")
			err := cmd.Run()
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
		})

		ginkgo.It("should start bridge between 2 custom archway chains", func() {
			dive.RunCustomArchwayNode0()
			dive.RunCustomArchwayNode1()
			cmd.Args = append(cmd.Args, "bridge", "ibc", "--chainA", "archway", "--chainB", "archway", "--chainAServiceName", "node-service-archway-node-0", "--chainBServiceName", "node-service-archway-node-1")
			err := cmd.Run()
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
		})

		ginkgo.It("should start bridge between 1 custom chain and running bridge command", func() {
			dive.RunCustomArchwayNode1()
			cmd.Args = append(cmd.Args, "bridge", "ibc", "--chainA", "archway", "--chainB", "archway", "--chainBServiceName", "node-service-archway-node-1")
			err := cmd.Run()
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
		})

		ginkgo.It("should start bridge between neutron and neutron by running one custom neutron chain", func() {
			dive.RunNeutronNode()
			dive.RunCustomNeutronNode1()
			cmd.Args = append(cmd.Args, "bridge", "ibc", "--chainA", "neutron", "--chainB", "neutron", "--chainAServiceName", "neutron-node-test-chain1", "--chainBServiceName", "neutron-node-test-chain2")
			err := cmd.Run()
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
		})

		ginkgo.It("should start bridge between 2 custom neutron chains", func() {
			dive.RunCustomNeutronNode0()
			dive.RunCustomNeutronNode1()
			cmd.Args = append(cmd.Args, "bridge", "ibc", "--chainA", "neutron", "--chainB", "neutron", "--chainAServiceName", "neutron-node-test-chain2", "--chainBServiceName", "neutron-node-test-chain2")
			err := cmd.Run()
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
		})

		ginkgo.It("should start bridge between 1 custom chain and running bridge command", func() {
			dive.RunCustomNeutronNode1()
			cmd.Args = append(cmd.Args, "bridge", "ibc", "--chainA", "neutron", "--chainB", "neutron", "--chainBServiceName", "neutron-node-test-chain2")
			err := cmd.Run()
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
		})

		ginkgo.It("should start bridge between archway and neutron chains", func() {

			cmd.Args = append(cmd.Args, "bridge", "ibc", "--chainA", "archway", "--chainB", "neutron")
			err := cmd.Run()
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
		})

		ginkgo.It("should start bridge between already running archway and neutron chains", func() {
			dive.RunArchwayNode()
			dive.RunNeutronNode()
			cmd.Args = append(cmd.Args, "bridge", "ibc", "--chainA", "archway", "--chainB", "neutron", "--chainAServiceName", "node-service-constantine-3", "--chainBServiceName", "neutron-node-test-chain1")
			err := cmd.Run()
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
		})

		ginkgo.It("should start bridge between already running archway and neutron chains with custom configuration", func() {
			dive.RunCustomNeutronNode0()
			dive.RunCustomArchwayNode0()
			cmd.Args = append(cmd.Args, "bridge", "ibc", "--chainA", "archway", "--chainB", "neutron", "--chainAServiceName", "node-service-archway-node-0", "--chainBServiceName", "neutron-node-test-chain2")
			err := cmd.Run()
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
		})

		ginkgo.It("should start IBC relay between icon and achway", func() {
			cmd.Args = append(cmd.Args, "bridge", "ibc", "--chainA", "icon", "--chainB", "archway")
			err := cmd.Run()
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
		})

		ginkgo.It("should start IBC relay between icon and neutron", func() {
			cmd.Args = append(cmd.Args, "bridge", "ibc", "--chainA", "icon", "--chainB", "neutron")
			err := cmd.Run()
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
		})

		ginkgo.It("should start bridge between icon and hardhat by running icon node first and running bridge command directly", func() {
			dive.RunDecentralizedCustomIconNode0()
			cmd.Args = append(cmd.Args, "bridge", "btp", "--chainA", "icon", "--chainB", "hardhat", "--chainAServiceName", "icon-node-0xacbc4e")
			err := cmd.Run()
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
		})

		ginkgo.It("should start bridge between icon and hardhat by running hardhat node first and running bridge command directly", func() {
			dive.RunHardhatNode()
			cmd.Args = append(cmd.Args, "bridge", "btp", "--chainA", "hardhat", "--chainB", "icon", "--chainAServiceName", "hardhat-node")
			err := cmd.Run()
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
		})

		ginkgo.It("should start bridge between icon and eth by running icon node first and running bridge command directly", func() {
			dive.RunDecentralizedCustomIconNode0()
			cmd.Args = append(cmd.Args, "bridge", "btp", "--chainA", "icon", "--chainB", "eth", "--chainAServiceName", "icon-node-0xacbc4e")
			err := cmd.Run()
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
		})

		ginkgo.It("should start bridge between icon and eth by running eth node first and running bridge command directly", func() {
			dive.RunEthNode()
			cmd.Args = append(cmd.Args, "bridge", "btp", "--chainA", "eth", "--chainB", "icon", "--chainAServiceName", "el-1-geth-lighthouse")
			err := cmd.Run()
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
		})

		ginkgo.It("should start bridge between icon and icon by running icon node first and running bridge command directly", func() {
			dive.RunDecentralizedCustomIconNode0()
			cmd.Args = append(cmd.Args, "bridge", "btp", "--chainA", "icon", "--chainB", "icon", "--chainAServiceName", "icon-node-0xacbc4e")
			err := cmd.Run()
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
		})
	})

	ginkgo.Describe("Other commands", func() {
		ginkgo.It("should handle error when trying to clean if no enclaves are running", func() {
			dive.Clean()
			dive.Clean()
		})

		ginkgo.It("should handle error when trying to clean if kurtosis engine is not running", func() {
			cmd1 := exec.Command("kurtosis", "engine", "stop")
			cmd1.Run()
			bin := dive.GetBinPath()
			cmd2 := exec.Command(bin, "clean")
			err := cmd2.Run()
			gomega.Expect(err).To(gomega.HaveOccurred())
			cmd3 := exec.Command("kurtosis", "engine", "start")
			cmd3.Run()
		})

		ginkgo.It("should handle invalid input for chain command", func() {
			cmd.Args = append(cmd.Args, "chain", "invalid_input")
			err := cmd.Run()
			gomega.Expect(err).To(gomega.HaveOccurred())
		})
	})

	ginkgo.Describe("Icon chain commands", func() {
		ginkgo.It("should run single icon node", func() {
			cmd.Args = append(cmd.Args, "chain", "icon")
			err := cmd.Run()
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
		})

		ginkgo.It("should run single icon node with verbose flag enabled", func() {
			cmd.Args = append(cmd.Args, "chain", "icon", "--verbose")
			err := cmd.Run()
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
		})

		ginkgo.It("should run single icon node along with decentralisation", func() {
			cmd.Args = append(cmd.Args, "chain", "icon", "-d")
			err := cmd.Run()
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
		})

		ginkgo.It("should run single icon node along with decentralisation with verbose flag enabled", func() {
			cmd.Args = append(cmd.Args, "chain", "icon", "-d", "--verbose")
			err := cmd.Run()
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
		})

		ginkgo.It("should run custom Icon node", func() {
			cmd.Args = append(cmd.Args, "chain", "icon", "-c", "../../cli/sample-jsons/config0.json", "-g", "../../services/jvm/icon/static-files/config/genesis-icon-0.zip")
			err := cmd.Run()
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
		})

		ginkgo.It("should run custom Icon node  with verbose flag enabled", func() {
			cmd.Args = append(cmd.Args, "chain", "icon", "-c", "../../cli/sample-jsons/config0.json", "-g", "../../services/jvm/icon/static-files/config/genesis-icon-0.zip", "--verbose")
			err := cmd.Run()
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
		})

		ginkgo.It("should run custom Icon node-1", func() {
			cmd.Args = append(cmd.Args, "chain", "icon", "-c", "../../cli/sample-jsons/config1.json", "-g", "../../services/jvm/icon/static-files/config/genesis-icon-1.zip")
			err := cmd.Run()
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
		})

		ginkgo.It("should run custom Icon node-1  with verbose flag enabled", func() {
			cmd.Args = append(cmd.Args, "chain", "icon", "-c", "../../cli/sample-jsons/config1.json", "-g", "../../services/jvm/icon/static-files/config/genesis-icon-1.zip", "--verbose")
			err := cmd.Run()
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
		})

		ginkgo.It("should run icon node first and then decentralise it", func() {
			dive.RunIconNode()
			cmd.Args = append(cmd.Args, "chain", "icon", "decentralize", "-p", "gochain", "-k", "keystores/keystore.json", "-n", "0x3", "-e", "http://172.16.0.3:9080/api/v3/icon_dex", "-s", "icon-node-0xacbc4e")
			err := cmd.Run()
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
		})

		ginkgo.It("should run icon node first and then decentralise it with verbose flag enabled", func() {
			dive.RunIconNode()
			cmd.Args = append(cmd.Args, "chain", "icon", "decentralize", "-p", "gochain", "-k", "keystores/keystore.json", "-n", "0x3", "-e", "http://172.16.0.3:9080/api/v3/icon_dex", "-s", "icon-node-0xacbc4e", "--verbose")
			err := cmd.Run()
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
		})

		ginkgo.It("should handle invalid input for chain command", func() {
			cmd.Args = append(cmd.Args, "chain", "icon", "-c", "invalid.json", "-g", "../../services/jvm/icon/static-files/config/genesis-icon-0.zip")
			err := cmd.Run()
			gomega.Expect(err).To(gomega.HaveOccurred())
		})

		ginkgo.It("should handle invalid input for chain command", func() {
			cmd.Args = append(cmd.Args, "chain", "icon", "-c", "../../cli/sample-jsons/config0.json", "-g", "../../services/jvm/icon/static-files/config/invalid-icon-3.zip")
			err := cmd.Run()
			gomega.Expect(err).To(gomega.HaveOccurred())
		})

		ginkgo.It("should handle invalid input for chain command", func() {
			cmd.Args = append(cmd.Args, "chain", "icon", "-c", "../../cli/sample-jsons/invalid_config.json", "-g", "../../services/jvm/icon/static-files/config/invalid-icon-3.zip")
			err := cmd.Run()
			gomega.Expect(err).To(gomega.HaveOccurred())
		})

		ginkgo.It("should handle invalid input for chain command", func() {
			dive.RunIconNode()
			cmd.Args = append(cmd.Args, "chain", "icon", "decentralize", "-p", "invalidPassword", "-k", "keystores/keystore.json", "-n", "0x3", "-e", "http://172.16.0.3:9080/api/v3/icon_dex", "-s", "icon-node-0xacbc4e")
			err := cmd.Run()
			gomega.Expect(err).To(gomega.HaveOccurred())
		})

		ginkgo.It("should handle invalid input for chain command", func() {
			dive.RunIconNode()
			cmd.Args = append(cmd.Args, "chain", "icon", "decentralize", "-p", "gochain", "-k", "keystores/invalid.json", "-n", "0x3", "-e", "http://172.16.0.3:9080/api/v3/icon_dex", "-s", "icon-node-0xacbc4e")
			err := cmd.Run()
			gomega.Expect(err).To(gomega.HaveOccurred())
		})

		ginkgo.It("should handle invalid input for chain command", func() {
			dive.RunIconNode()
			cmd.Args = append(cmd.Args, "chain", "icon", "decentralize", "-p", "gochain", "-k", "keystores/keystore.json", "-n", "0x9", "-e", "http://172.16.0.3:9080/api/v3/icon_dex", "-s", "icon-node-0xacbc4e")
			err := cmd.Run()
			gomega.Expect(err).To(gomega.HaveOccurred())
		})

		ginkgo.It("should handle invalid input for chain command", func() {
			dive.RunIconNode()
			cmd.Args = append(cmd.Args, "chain", "icon", "decentralize", "-p", "gochain", "-k", "keystores/keystore.json", "-n", "0x3", "-e", "http://172.16.0.3:9081/api/v3/icon_dex", "-s", "icon-node-0xacbc4e")
			err := cmd.Run()
			gomega.Expect(err).To(gomega.HaveOccurred())
		})

		ginkgo.It("should handle invalid input for chain command", func() {
			dive.RunIconNode()
			cmd.Args = append(cmd.Args, "chain", "icon", "decentralize", "-p", "gochain", "-k", "keystores/keystore.json", "-n", "0x3", "-e", "http://172.16.0.3:9080/api/v3/icon_dex", "-s", "icon-node")
			err := cmd.Run()
			gomega.Expect(err).To(gomega.HaveOccurred())
		})

		ginkgo.It("should output user that chain is already running when trying to run icon chain that is already running", func() {
			dive.RunIconNode()
			cmd.Args = append(cmd.Args, "chain", "icon")
			err := cmd.Run()
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
		})

	})

	ginkgo.Describe("Eth chain commands", func() {
		ginkgo.It("should run single eth node", func() {
			cmd.Args = append(cmd.Args, "chain", "eth")
			err := cmd.Run()
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
		})

		ginkgo.It("should run single eth node with verbose flag enabled", func() {
			cmd.Args = append(cmd.Args, "chain", "eth", "--verbose")
			err := cmd.Run()
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
		})

		ginkgo.It("should output user that chain is already running when trying to run eth chain that is already running", func() {
			dive.RunEthNode()
			cmd.Args = append(cmd.Args, "chain", "eth")
			err := cmd.Run()
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
		})

	})

	ginkgo.Describe("Hardhat chain commands", func() {
		ginkgo.It("should run single hardhat node", func() {
			cmd.Args = append(cmd.Args, "chain", "hardhat")
			err := cmd.Run()
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
		})

		ginkgo.It("should run single hardhat node with verbose flag enabled", func() {
			cmd.Args = append(cmd.Args, "chain", "hardhat", "--verbose")
			err := cmd.Run()
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
		})

		ginkgo.It("should output user that chain is already running when trying to run hardhat chain that is already running", func() {
			dive.RunHardhatNode()
			cmd.Args = append(cmd.Args, "chain", "hardhat")
			err := cmd.Run()
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
		})
	})

	ginkgo.Describe("Archway chain commands", func() {
		ginkgo.It("should run single archway node", func() {
			dive.RunArchwayNode()
		})

		ginkgo.It("should run single archway node with verbose flag enabled", func() {
			cmd.Args = append(cmd.Args, "chain", "archway", "--verbose")
			err := cmd.Run()
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
		})

		ginkgo.It("should run single custom archway node", func() {
			dive.RunCustomArchwayNode0()
		})

		ginkgo.It("should run single custom archway node with verbose flag enabled", func() {
			cmd.Args = append(cmd.Args, "chain", "archway", "-c", "../../cli/sample-jsons/archway.json", "--verbose")
			err := cmd.Run()
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
		})

		ginkgo.It("should run single custom archway node with invalid json path", func() {
			cmd.Args = append(cmd.Args, "chain", "archway", "-c", "../../cli/sample-jsons/archway4.json")
			err := cmd.Run()
			gomega.Expect(err).To(gomega.HaveOccurred())
		})

		ginkgo.It("should output user that chain is already running when trying to run archway chain that is already running", func() {
			dive.RunArchwayNode()
			cmd.Args = append(cmd.Args, "chain", "archway")
			err := cmd.Run()
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
		})
	})

	ginkgo.Describe("Neutron chain commands", func() {
		ginkgo.It("should run single neutron node", func() {
			dive.RunNeutronNode()
		})

		ginkgo.It("should run single nurtron node with verbose flag enabled", func() {
			cmd.Args = append(cmd.Args, "chain", "neutron", "--verbose")
			err := cmd.Run()
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
		})

		ginkgo.It("should run single custom neutron node", func() {
			dive.RunCustomNeutronNode0()
		})

		ginkgo.It("should run single custom neutron node with verbose flag enabled", func() {
			cmd.Args = append(cmd.Args, "chain", "neutron", "-c", "../../cli/sample-jsons/neutron.json", "--verbose")
			err := cmd.Run()
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
		})

		ginkgo.It("should run single custom neutron node with invalid json path", func() {
			cmd.Args = append(cmd.Args, "chain", "neutron", "-c", "../../cli/sample-jsons/neutron5.json")
			err := cmd.Run()
			gomega.Expect(err).To(gomega.HaveOccurred())
		})

		ginkgo.It("should output user that chain is already running when trying to run neutron chain that is already running", func() {
			dive.RunNeutronNode()
			cmd.Args = append(cmd.Args, "chain", "neutron")
			err := cmd.Run()
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
		})
	})
})
