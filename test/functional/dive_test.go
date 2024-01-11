package dive_test

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"testing"

	dive "github.com/HugoByte/DIVE/test/functional"
	"github.com/hugobyte/dive-core/cli/cmd/utility"
	"github.com/hugobyte/dive-core/cli/common"
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

	ginkgo.BeforeEach(func() {
		cmd = dive.GetBinaryCommand()
		cmd.Stdout = &testWriter{}
		cmd.Stderr = &testWriter{}
	})


	ginkgo.Describe("Smoke Tests", func() {
		ginkgo.It("should display the correct version", func() {
			cmd.Args = append(cmd.Args, "version")
			cmd.Stdout = &stdout
			err := cmd.Run()
			fmt.Println(stdout.String())
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
			enclaveName := dive.GenerateRandomName()
			cli := common.GetCli(enclaveName)
			latestVersion := utility.GetLatestVersion(cli)
			gomega.Expect(stdout.String()).To(gomega.ContainSubstring(latestVersion))
		})

		ginkgo.It("should start bridge between icon and eth correctly", func() {
			enclaveName := dive.GenerateRandomName()
			cmd.Args = append(cmd.Args, "bridge", "btp", "--chainA", "icon", "--chainB", "eth", "--enclaveName", enclaveName)
			err := cmd.Run()
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
			dive.Clean(enclaveName)
		})

		ginkgo.It("should start bridge between icon and hardhat but with icon bridge set to true", func() {
			enclaveName := dive.GenerateRandomName()
			cmd.Args = append(cmd.Args, "bridge", "btp", "--chainA", "icon", "--chainB", "hardhat", "--bmvbridge", "--enclaveName", enclaveName)
			err := cmd.Run()
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
			dive.Clean(enclaveName)
		})

		ginkgo.It("should start bridge between icon and icon", func() {
			enclaveName := dive.GenerateRandomName()
			cmd.Args = append(cmd.Args, "bridge", "btp", "--chainA", "icon", "--chainB", "icon", "--enclaveName", enclaveName)
			err := cmd.Run()
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
			dive.Clean(enclaveName)
		})

		ginkgo.It("should start bridge between archway and archway using ibc", func() {
			enclaveName := dive.GenerateRandomName()
			cmd.Args = append(cmd.Args, "bridge", "ibc", "--chainA", "archway", "--chainB", "archway", "--enclaveName", enclaveName)
			err := cmd.Run()
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
			dive.Clean(enclaveName)
		})

	})

	ginkgo.Describe("Bridge command Test", func() {
		ginkgo.It("should start bridge between icon and eth but with icon bridge set to true", func() {
			enclaveName := dive.GenerateRandomName()
			cmd.Args = append(cmd.Args, "bridge", "btp", "--chainA", "icon", "--chainB", "eth", "--bmvbridge", "--enclaveName", enclaveName)
			err := cmd.Run()
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
			dive.Clean(enclaveName)
		})

		ginkgo.It("should start bridge between icon and hardhat", func() {
			enclaveName := dive.GenerateRandomName()
			cmd.Args = append(cmd.Args, "bridge", "btp", "--chainA", "icon", "--chainB", "hardhat", "--enclaveName", enclaveName)
			err := cmd.Run()
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
			dive.Clean(enclaveName)
		})

		ginkgo.It("should start bridge between icon and eth by running each chain individually", func() {
			enclaveName := dive.GenerateRandomName()
			dive.RunDecentralizedIconNode(enclaveName)
			dive.RunEthNode(enclaveName)
			cmd.Args = append(cmd.Args, "bridge", "btp", "--chainA", "icon", "--chainB", "eth", "--chainAServiceName", dive.ICON_CONFIG0_SERVICENAME, "--chainBServiceName", dive.ETH_SERVICENAME, "--enclaveName", enclaveName)
			err := cmd.Run()
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
			dive.Clean(enclaveName)
		})

		ginkgo.It("should start bridge between icon and hardhat by running each chain individually", func() {
			enclaveName := dive.GenerateRandomName()
			dive.RunDecentralizedIconNode(enclaveName)
			dive.RunHardhatNode(enclaveName)
			cmd.Args = append(cmd.Args, "bridge", "btp", "--chainA", "icon", "--chainB", "hardhat", "--chainAServiceName", dive.ICON_CONFIG0_SERVICENAME, "--chainBServiceName", dive.HARDHAT_SERVICENAME, "--enclaveName", enclaveName)
			err := cmd.Run()
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
			dive.Clean(enclaveName)
		})

		ginkgo.It("should start bridge between icon and eth by running icon node first and then decentralising it", func() {
			enclaveName := dive.GenerateRandomName()
			dive.RunIconNode(enclaveName)

			serviceName0, endpoint0, nid0 := dive.GetServiceDetails(fmt.Sprintf("services_%s.json", enclaveName), dive.ICON_CONFIG0_SERVICENAME)
			dive.DecentralizeCustomIconNode(nid0, endpoint0, serviceName0, enclaveName)

			dive.RunEthNode(enclaveName)

			cmd.Args = append(cmd.Args, "bridge", "btp", "--chainA", "icon", "--chainB", "eth", "--chainAServiceName", dive.ICON_CONFIG0_SERVICENAME, "--chainBServiceName", dive.ETH_SERVICENAME, "--enclaveName", enclaveName)
			err := cmd.Run()
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
			defer os.Remove(fmt.Sprintf("updated-config-%s.json", enclaveName))
			dive.Clean(enclaveName)
		})

		ginkgo.It("should start bridge between icon and hardhat by running icon node first and then decentralising it", func() {
			enclaveName := dive.GenerateRandomName()

			dive.RunIconNode(enclaveName)
			serviceName0, endpoint0, nid0 := dive.GetServiceDetails(fmt.Sprintf("services_%s.json", enclaveName), dive.ICON_CONFIG0_SERVICENAME)
			dive.DecentralizeCustomIconNode(nid0, endpoint0, serviceName0, enclaveName)

			dive.RunHardhatNode(enclaveName)
			cmd.Args = append(cmd.Args, "bridge", "btp", "--chainA", "icon", "--chainB", "hardhat", "--chainAServiceName", dive.ICON_CONFIG0_SERVICENAME, "--chainBServiceName", dive.HARDHAT_SERVICENAME, "--enclaveName", enclaveName)
			err := cmd.Run()
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
			defer os.Remove(fmt.Sprintf("updated-config-%s.json", enclaveName))
			dive.Clean(enclaveName)
		})

		ginkgo.It("should start bridge between icon and icon by running one custom icon chain", func() {
			enclaveName := dive.GenerateRandomName()
			dive.RunDecentralizedIconNode(enclaveName)
			dive.RunDecentralizedCustomIconNode1(enclaveName)
			cmd.Args = append(cmd.Args, "bridge", "btp", "--chainA", "icon", "--chainB", "icon", "--chainAServiceName", dive.ICON_CONFIG0_SERVICENAME, "--chainBServiceName", dive.ICON_CONFIG1_SERVICENAME, "--enclaveName", enclaveName)
			err := cmd.Run()
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
			dive.Clean(enclaveName)
		})

		ginkgo.It("should start bridge between icon and running custom icon later decentralising it", func() {
			enclaveName := dive.GenerateRandomName()
			dive.RunDecentralizedIconNode(enclaveName)
			dive.RunCustomIconNode1(enclaveName)

			serviceName, endpoint, nid := dive.GetServiceDetails(fmt.Sprintf("services_%s.json", enclaveName), dive.ICON_CONFIG1_SERVICENAME)
			dive.DecentralizeCustomIconNode(nid, endpoint, serviceName, enclaveName)

			cmd.Args = append(cmd.Args, "bridge", "btp", "--chainA", "icon", "--chainB", "icon", "--chainAServiceName", dive.ICON_CONFIG0_SERVICENAME, "--chainBServiceName", dive.ICON_CONFIG1_SERVICENAME, "--enclaveName", enclaveName)
			err := cmd.Run()
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
			defer os.Remove(fmt.Sprintf("updated-config-%s.json", enclaveName))
			dive.Clean(enclaveName)
		})

		ginkgo.It("should start bridge between icon and icon by running one icon chain and later decentralsing it. Running another custom icon chain and then decentralising it", func() {
			enclaveName := dive.GenerateRandomName()
			dive.RunIconNode(enclaveName)
			dive.RunCustomIconNode1(enclaveName)

			serviceName0, endpoint0, nid0 := dive.GetServiceDetails(fmt.Sprintf("services_%s.json", enclaveName), dive.ICON_CONFIG0_SERVICENAME)
			dive.DecentralizeCustomIconNode(nid0, endpoint0, serviceName0, enclaveName)

			serviceName1, endpoint1, nid1 := dive.GetServiceDetails(fmt.Sprintf("services_%s.json", enclaveName), dive.ICON_CONFIG1_SERVICENAME)
			dive.DecentralizeCustomIconNode(nid1, endpoint1, serviceName1, enclaveName)

			cmd.Args = append(cmd.Args, "bridge", "btp", "--chainA", "icon", "--chainB", "icon", "--chainAServiceName", dive.ICON_CONFIG0_SERVICENAME, "--chainBServiceName", dive.ICON_CONFIG1_SERVICENAME, "--enclaveName", enclaveName)
			err := cmd.Run()
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
			defer os.Remove(fmt.Sprintf("updated-config-%s.json", enclaveName))
			dive.Clean(enclaveName)
		})

		ginkgo.It("should start bridge between 2 custom icon chains", func() {
			enclaveName := dive.GenerateRandomName()
			dive.RunDecentralizedCustomIconNode0(enclaveName)
			dive.RunDecentralizedCustomIconNode1(enclaveName)
			cmd.Args = append(cmd.Args, "bridge", "btp", "--chainA", "icon", "--chainB", "icon", "--chainAServiceName", dive.ICON_CONFIG0_SERVICENAME, "--chainBServiceName", dive.ICON_CONFIG1_SERVICENAME, "--enclaveName", enclaveName)
			err := cmd.Run()
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
			dive.Clean(enclaveName)
		})

		ginkgo.It("should start bridge between 2 custom icon chains by running them first and then decentralising it later", func() {
			enclaveName := dive.GenerateRandomName()
			dive.RunCustomIconNode0(enclaveName)
			dive.RunCustomIconNode1(enclaveName)

			serviceName0, endpoint0, nid0 := dive.GetServiceDetails(fmt.Sprintf("services_%s.json", enclaveName), dive.ICON_CONFIG0_SERVICENAME)
			dive.DecentralizeCustomIconNode(nid0, endpoint0, serviceName0, enclaveName)

			serviceName1, endpoint1, nid1 := dive.GetServiceDetails(fmt.Sprintf("services_%s.json", enclaveName), dive.ICON_CONFIG1_SERVICENAME)
			dive.DecentralizeCustomIconNode(nid1, endpoint1, serviceName1, enclaveName)

			cmd.Args = append(cmd.Args, "bridge", "btp", "--chainA", "icon", "--chainB", "icon", "--chainAServiceName", dive.ICON_CONFIG0_SERVICENAME, "--chainBServiceName", dive.ICON_CONFIG1_SERVICENAME, "--enclaveName", enclaveName)
			err := cmd.Run()
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
			defer os.Remove(fmt.Sprintf("updated-config-%s.json", enclaveName))
			dive.Clean(enclaveName)
		})

		ginkgo.It("should start bridge between 2 chains when all nodes are running", func() {
			enclaveName := dive.GenerateRandomName()
			dive.RunDecentralizedIconNode(enclaveName)
			dive.RunEthNode(enclaveName)
			dive.RunHardhatNode(enclaveName)
			cmd.Args = append(cmd.Args, "bridge", "btp", "--chainA", "icon", "--chainB", "eth", "--chainAServiceName", dive.ICON_CONFIG0_SERVICENAME, "--chainBServiceName", dive.ETH_SERVICENAME, "--enclaveName", enclaveName)
			err := cmd.Run()
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
			dive.Clean(enclaveName)
		})

		ginkgo.It("should handle invalid input for bridge command", func() {
			enclaveName := dive.GenerateRandomName()
			cmd.Args = append(cmd.Args, "bridge", "btp", "--chainA", "icon", "--chainB", "invalid_input", "--enclaveName", enclaveName)
			err := cmd.Run()
			gomega.Expect(err).To(gomega.HaveOccurred())
		})

		ginkgo.It("should handle invalid input for bridge command", func() {
			enclaveName := dive.GenerateRandomName()
			cmd.Args = append(cmd.Args, "bridge", "btp", "--chainA", "invalid_input", "--chainB", "eth", "--enclaveName", enclaveName)
			err := cmd.Run()
			gomega.Expect(err).To(gomega.HaveOccurred())
		})

		ginkgo.It("should handle invalid input ibc bridge command", func() {
			enclaveName := dive.GenerateRandomName()
			cmd.Args = append(cmd.Args, "bridge", "ibc", "--chainA", "archway", "--chainB", "invalid", "--enclaveName", enclaveName)
			err := cmd.Run()
			gomega.Expect(err).To(gomega.HaveOccurred())
		})

		ginkgo.It("should start bridge between archway and archway by running one custom archway chain", func() {
			enclaveName := dive.GenerateRandomName()
			dive.RunArchwayNode(enclaveName)
			dive.RunCustomArchwayNode1(enclaveName)
			cmd.Args = append(cmd.Args, "bridge", "ibc", "--chainA", "archway", "--chainB", "archway", "--chainAServiceName", dive.DEFAULT_ARCHWAY_SERVICENAME, "--chainBServiceName", dive.ARCHWAY_CONFIG1_SERVICENAME, "--enclaveName", enclaveName)
			err := cmd.Run()
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
			defer os.Remove(fmt.Sprintf("updated-config-%s.json", enclaveName))
			dive.Clean(enclaveName)
		})

		ginkgo.It("should start bridge between 2 custom archway chains", func() {
			enclaveName := dive.GenerateRandomName()
			dive.RunCustomArchwayNode0(enclaveName)
			dive.RunCustomArchwayNode1(enclaveName)
			cmd.Args = append(cmd.Args, "bridge", "ibc", "--chainA", "archway", "--chainB", "archway", "--chainAServiceName", dive.ARCHWAY_CONFIG0_SERVICENAME, "--chainBServiceName", dive.ARCHWAY_CONFIG1_SERVICENAME, "--enclaveName", enclaveName)
			err := cmd.Run()
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
			defer os.Remove(fmt.Sprintf("updated-config-%s.json", enclaveName))
			dive.Clean(enclaveName)
		})

		ginkgo.It("should start bridge between archway to archway with 1 custom chain parameter", func() {
			enclaveName := dive.GenerateRandomName()
			dive.RunCustomArchwayNode1(enclaveName)
			cmd.Args = append(cmd.Args, "bridge", "ibc", "--chainA", "archway", "--chainB", "archway", "--chainBServiceName", dive.ARCHWAY_CONFIG1_SERVICENAME, "--enclaveName", enclaveName)
			err := cmd.Run()
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
			defer os.Remove(fmt.Sprintf("updated-config-%s.json", enclaveName))
			dive.Clean(enclaveName)
		})

		ginkgo.It("should start bridge between neutron and neutron by running one custom neutron chain", func() {
			enclaveName := dive.GenerateRandomName()
			dive.RunNeutronNode(enclaveName)
			dive.RunCustomNeutronNode1(enclaveName)
			cmd.Args = append(cmd.Args, "bridge", "ibc", "--chainA", "neutron", "--chainB", "neutron", "--chainAServiceName", dive.DEFAULT_NEUTRON_SERVICENAME, "--chainBServiceName", dive.NEUTRON_CONFIG1_SERVICENAME, "--enclaveName", enclaveName)
			err := cmd.Run()
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
			defer os.Remove(fmt.Sprintf("updated-config-%s.json", enclaveName))
			dive.Clean(enclaveName)
		})

		ginkgo.It("should start bridge between 2 custom neutron chains", func() {
			enclaveName := dive.GenerateRandomName()
			dive.RunCustomNeutronNode0(enclaveName)
			dive.RunCustomNeutronNode1(enclaveName)
			cmd.Args = append(cmd.Args, "bridge", "ibc", "--chainA", "neutron", "--chainB", "neutron", "--chainAServiceName", dive.NEUTRON_CONFIG0_SERVICENAME, "--chainBServiceName", dive.NEUTRON_CONFIG1_SERVICENAME, "--enclaveName", enclaveName)
			err := cmd.Run()
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
			defer os.Remove(fmt.Sprintf("updated-config-%s.json", enclaveName))
			dive.Clean(enclaveName)
		})

		ginkgo.It("should start IBC relay between nuetron to neutron with one 1 custom chain.", func() {
			enclaveName := dive.GenerateRandomName()
			dive.RunCustomNeutronNode1(enclaveName)
			cmd.Args = append(cmd.Args, "bridge", "ibc", "--chainA", "neutron", "--chainB", "neutron", "--chainBServiceName", dive.NEUTRON_CONFIG1_SERVICENAME, "--enclaveName", enclaveName)
			err := cmd.Run()
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
			defer os.Remove(fmt.Sprintf("updated-config-%s.json", enclaveName))
			dive.Clean(enclaveName)
		})

		ginkgo.It("should start bridge between archway and neutron chains", func() {
			enclaveName := dive.GenerateRandomName()
			cmd.Args = append(cmd.Args, "bridge", "ibc", "--chainA", "archway", "--chainB", "neutron", "--enclaveName", enclaveName)
			err := cmd.Run()
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
			dive.Clean(enclaveName)
		})

		ginkgo.It("should start bridge between already running archway and neutron chains", func() {
			enclaveName := dive.GenerateRandomName()
			dive.RunArchwayNode(enclaveName)
			dive.RunNeutronNode(enclaveName)
			cmd.Args = append(cmd.Args, "bridge", "ibc", "--chainA", "archway", "--chainB", "neutron", "--chainAServiceName", dive.DEFAULT_ARCHWAY_SERVICENAME, "--chainBServiceName", dive.DEFAULT_NEUTRON_SERVICENAME, "--enclaveName", enclaveName)
			err := cmd.Run()
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
			dive.Clean(enclaveName)
		})

		ginkgo.It("should start bridge between already running archway and neutron chains with custom configuration", func() {
			enclaveName := dive.GenerateRandomName()
			dive.RunCustomNeutronNode0(enclaveName)
			dive.RunCustomArchwayNode0(enclaveName)
			cmd.Args = append(cmd.Args, "bridge", "ibc", "--chainA", "archway", "--chainB", "neutron", "--chainAServiceName", dive.ARCHWAY_CONFIG0_SERVICENAME, "--chainBServiceName", dive.NEUTRON_CONFIG0_SERVICENAME, "--enclaveName", enclaveName)
			err := cmd.Run()
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
			defer os.Remove(fmt.Sprintf("updated-config-%s.json", enclaveName))
			dive.Clean(enclaveName)
		})

		ginkgo.It("should start IBC relay between icon and archway", func() {
			enclaveName := dive.GenerateRandomName()
			cmd.Args = append(cmd.Args, "bridge", "ibc", "--chainA", "icon", "--chainB", "archway", "--enclaveName", enclaveName)
			err := cmd.Run()
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
			dive.Clean(enclaveName)
		})

		ginkgo.It("should start IBC relay between icon and neutron", func() {
			enclaveName := dive.GenerateRandomName()
			cmd.Args = append(cmd.Args, "bridge", "ibc", "--chainA", "icon", "--chainB", "neutron", "--enclaveName", enclaveName)
			err := cmd.Run()
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
			dive.Clean(enclaveName)
		})

		ginkgo.It("should start IBC relay between already running icon and archway chain", func() {
			enclaveName := dive.GenerateRandomName()
			dive.RunIconNode(enclaveName)
			dive.RunArchwayNode(enclaveName)
			cmd.Args = append(cmd.Args, "bridge", "ibc", "--chainA", "icon", "--chainB", "archway", "--chainAServiceName", dive.ICON_CONFIG0_SERVICENAME, "--chainBServiceName", dive.DEFAULT_ARCHWAY_SERVICENAME, "--enclaveName", enclaveName)
			err := cmd.Run()
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
			dive.Clean(enclaveName)
		})

		ginkgo.It("should start IBC relay between already running icon and neutron chain", func() {
			enclaveName := dive.GenerateRandomName()
			dive.RunIconNode(enclaveName)
			dive.RunNeutronNode(enclaveName)
			cmd.Args = append(cmd.Args, "bridge", "ibc", "--chainA", "icon", "--chainB", "neutron", "--chainAServiceName", dive.ICON_CONFIG0_SERVICENAME, "--chainBServiceName", dive.DEFAULT_NEUTRON_SERVICENAME, "--enclaveName", enclaveName)
			err := cmd.Run()
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
			dive.Clean(enclaveName)
		})

		ginkgo.It("should start IBC relay between already running icon and custom archway chain", func() {
			enclaveName := dive.GenerateRandomName()
			dive.RunIconNode(enclaveName)
			dive.RunCustomArchwayNode0(enclaveName)
			cmd.Args = append(cmd.Args, "bridge", "ibc", "--chainA", "icon", "--chainB", "archway", "--chainAServiceName", dive.ICON_CONFIG0_SERVICENAME, "--chainBServiceName", dive.ARCHWAY_CONFIG0_SERVICENAME, "--enclaveName", enclaveName)
			err := cmd.Run()
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
			defer os.Remove(fmt.Sprintf("updated-config-%s.json", enclaveName))
			dive.Clean(enclaveName)
		})

		ginkgo.It("should start IBC relay between already running icon and custom neutron chain", func() {
			enclaveName := dive.GenerateRandomName()
			dive.RunIconNode(enclaveName)
			dive.RunCustomNeutronNode0(enclaveName)
			cmd.Args = append(cmd.Args, "bridge", "ibc", "--chainA", "icon", "--chainB", "neutron", "--chainAServiceName", dive.ICON_CONFIG0_SERVICENAME, "--chainBServiceName", dive.NEUTRON_CONFIG0_SERVICENAME, "--enclaveName", enclaveName)
			err := cmd.Run()
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
			defer os.Remove(fmt.Sprintf("updated-config-%s.json", enclaveName))
			dive.Clean(enclaveName)
		})

		ginkgo.It("should start IBC relay between already running custom icon and archway chain", func() {
			enclaveName := dive.GenerateRandomName()
			dive.RunCustomIconNode0(enclaveName)
			dive.RunArchwayNode(enclaveName)
			cmd.Args = append(cmd.Args, "bridge", "ibc", "--chainA", "icon", "--chainB", "archway", "--chainAServiceName", dive.ICON_CONFIG0_SERVICENAME, "--chainBServiceName", dive.DEFAULT_ARCHWAY_SERVICENAME, "--enclaveName", enclaveName)
			err := cmd.Run()
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
			defer os.Remove(fmt.Sprintf("updated-config-%s.json", enclaveName))
			dive.Clean(enclaveName)
		})

		ginkgo.It("should start IBC relay between already running custom icon and neutron chain", func() {
			enclaveName := dive.GenerateRandomName()
			dive.RunCustomIconNode1(enclaveName)
			dive.RunNeutronNode(enclaveName)
			cmd.Args = append(cmd.Args, "bridge", "ibc", "--chainA", "icon", "--chainB", "neutron", "--chainAServiceName", dive.ICON_CONFIG1_SERVICENAME, "--chainBServiceName", dive.DEFAULT_NEUTRON_SERVICENAME, "--enclaveName", enclaveName)
			err := cmd.Run()
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
			defer os.Remove(fmt.Sprintf("updated-config-%s.json", enclaveName))
			dive.Clean(enclaveName)
		})

		ginkgo.It("should start IBC relay between already running custom icon and custom archway chain", func() {
			enclaveName := dive.GenerateRandomName()
			dive.RunCustomIconNode1(enclaveName)
			dive.RunCustomArchwayNode0(enclaveName)
			cmd.Args = append(cmd.Args, "bridge", "ibc", "--chainA", "icon", "--chainB", "archway", "--chainAServiceName", dive.ICON_CONFIG1_SERVICENAME, "--chainBServiceName", dive.ARCHWAY_CONFIG0_SERVICENAME, "--enclaveName", enclaveName)
			err := cmd.Run()
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
			defer os.Remove(fmt.Sprintf("updated-config-%s.json", enclaveName))
			dive.Clean(enclaveName)
		})

		ginkgo.It("should start IBC relay between already running custom icon and custom neutron chain", func() {
			enclaveName := dive.GenerateRandomName()
			dive.RunCustomIconNode1(enclaveName)
			dive.RunCustomNeutronNode0(enclaveName)
			cmd.Args = append(cmd.Args, "bridge", "ibc", "--chainA", "icon", "--chainB", "neutron", "--chainAServiceName", dive.ICON_CONFIG1_SERVICENAME, "--chainBServiceName", dive.NEUTRON_CONFIG0_SERVICENAME, "--enclaveName", enclaveName)
			err := cmd.Run()
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
			defer os.Remove(fmt.Sprintf("updated-config-%s.json", enclaveName))
			dive.Clean(enclaveName)
		})

		ginkgo.It("should start bridge between icon and hardhat by running icon node first and running bridge command directly", func() {
			enclaveName := dive.GenerateRandomName()
			dive.RunDecentralizedCustomIconNode0(enclaveName)
			cmd.Args = append(cmd.Args, "bridge", "btp", "--chainA", "icon", "--chainB", "hardhat", "--chainAServiceName", dive.ICON_CONFIG0_SERVICENAME, "--enclaveName", enclaveName)
			err := cmd.Run()
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
			defer os.Remove(fmt.Sprintf("updated-config-%s.json", enclaveName))
			dive.Clean(enclaveName)
		})

		ginkgo.It("should start bridge between icon and hardhat by running hardhat node first and running bridge command directly", func() {
			enclaveName := dive.GenerateRandomName()
			dive.RunHardhatNode(enclaveName)
			cmd.Args = append(cmd.Args, "bridge", "btp", "--chainA", "hardhat", "--chainB", "icon", "--chainAServiceName", dive.HARDHAT_SERVICENAME, "--enclaveName", enclaveName)
			err := cmd.Run()
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
			dive.Clean(enclaveName)
		})

		ginkgo.It("should start bridge between icon and eth by running icon node first and running bridge command directly", func() {
			enclaveName := dive.GenerateRandomName()
			dive.RunDecentralizedCustomIconNode0(enclaveName)
			cmd.Args = append(cmd.Args, "bridge", "btp", "--chainA", "icon", "--chainB", "eth", "--chainAServiceName", dive.ICON_CONFIG0_SERVICENAME, "--enclaveName", enclaveName)
			err := cmd.Run()
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
			defer os.Remove(fmt.Sprintf("updated-config-%s.json", enclaveName))
			dive.Clean(enclaveName)
		})

		ginkgo.It("should start bridge between icon and eth by running eth node first and running bridge command directly", func() {
			enclaveName := dive.GenerateRandomName()
			dive.RunEthNode(enclaveName)
			cmd.Args = append(cmd.Args, "bridge", "btp", "--chainA", "eth", "--chainB", "icon", "--chainAServiceName", dive.ETH_SERVICENAME, "--enclaveName", enclaveName)
			err := cmd.Run()
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
			dive.Clean(enclaveName)
		})

		ginkgo.It("should start bridge between icon and icon by running icon node first and running bridge command directly", func() {
			enclaveName := dive.GenerateRandomName()
			dive.RunDecentralizedCustomIconNode1(enclaveName)
			cmd.Args = append(cmd.Args, "bridge", "btp", "--chainA", "icon", "--chainB", "icon", "--chainAServiceName", dive.ICON_CONFIG1_SERVICENAME, "--enclaveName", enclaveName)
			err := cmd.Run()
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
			defer os.Remove(fmt.Sprintf("updated-config-%s.json", enclaveName))
			dive.Clean(enclaveName)
		})
	})

	ginkgo.Describe("Other commands", func() {

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
		ginkgo.It("should run single icon node testing", func() {
			enclaveName := dive.GenerateRandomName()
			cmd.Args = append(cmd.Args, "chain", "icon", "--enclaveName", enclaveName)
			err := cmd.Run()
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
			dive.Clean(enclaveName)
		})

		ginkgo.It("should run single icon node along with decentralisation", func() {
			enclaveName := dive.GenerateRandomName()
			cmd.Args = append(cmd.Args, "chain", "icon", "-d", "--enclaveName", enclaveName)
			err := cmd.Run()
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
			dive.Clean(enclaveName)
		})

		ginkgo.It("should run custom Icon node-0", func() {
			enclaveName := dive.GenerateRandomName()
			updated_path := dive.UpdatePublicPort(enclaveName, dive.ICON_CONFIG0)
			cmd.Args = append(cmd.Args, "chain", "icon", "-c", updated_path, "-g", dive.ICON_GENESIS0, "--enclaveName", enclaveName)
			defer os.Remove(updated_path)
			err := cmd.Run()
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
			dive.Clean(enclaveName)
		})

		ginkgo.It("should run custom Icon node-1", func() {
			enclaveName := dive.GenerateRandomName()
			filepath := dive.ICON_CONFIG1
			updated_path := dive.UpdatePublicPort(enclaveName, filepath)
			cmd.Args = append(cmd.Args, "chain", "icon", "-c", updated_path, "-g", dive.ICON_GENESIS1, "--enclaveName", enclaveName)
			defer os.Remove(updated_path)
			err := cmd.Run()
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
			dive.Clean(enclaveName)
		})

		ginkgo.It("should run icon node first and then decentralise it", func() {
			enclaveName := dive.GenerateRandomName()
			dive.RunIconNode(enclaveName)
			serviceName0, endpoint0, nid0 := dive.GetServiceDetails(fmt.Sprintf("services_%s.json", enclaveName), dive.ICON_CONFIG0_SERVICENAME)
			cmd.Args = append(cmd.Args, "chain", "icon", "decentralize", "-p", "gochain", "-k", "keystores/keystore.json", "-n", nid0, "-e", endpoint0, "-s", serviceName0, "--enclaveName", enclaveName)
			err := cmd.Run()
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
			dive.Clean(enclaveName)
		})

		ginkgo.It("should handle invalid input for chain command", func() {
			enclaveName := dive.GenerateRandomName()
			cmd.Args = append(cmd.Args, "chain", "icon", "-c", "invalid.json", "-g", dive.ICON_GENESIS0, "--enclaveName", enclaveName)
			err := cmd.Run()
			gomega.Expect(err).To(gomega.HaveOccurred())
			dive.Clean(enclaveName)
		})

		ginkgo.It("should handle invalid input for chain command", func() {
			enclaveName := dive.GenerateRandomName()
			cmd.Args = append(cmd.Args, "chain", "icon", "-c", dive.ICON_CONFIG0, "-g", "./config/invalid-icon-3.zip", "--enclaveName", enclaveName)
			err := cmd.Run()
			gomega.Expect(err).To(gomega.HaveOccurred())
			dive.Clean(enclaveName)
		})

		ginkgo.It("should handle invalid input for chain command", func() {
			enclaveName := dive.GenerateRandomName()
			cmd.Args = append(cmd.Args, "chain", "icon", "-c", "../../cli/sample-jsons/invalid_config.json", "-g", "./config/invalid-icon-3.zip", "--enclaveName", enclaveName)
			err := cmd.Run()
			gomega.Expect(err).To(gomega.HaveOccurred())
			dive.Clean(enclaveName)
		})

		ginkgo.It("should handle invalid input for chain command", func() {
			enclaveName := dive.GenerateRandomName()
			dive.RunIconNode(enclaveName)
			serviceName0, endpoint0, nid0 := dive.GetServiceDetails(fmt.Sprintf("services_%s.json", enclaveName), dive.ICON_CONFIG0_SERVICENAME)
			cmd.Args = append(cmd.Args, "chain", "icon", "decentralize", "-p", "invalidPassword", "-k", "keystores/keystore.json", "-n", nid0, "-e", endpoint0, "-s", serviceName0, "--enclaveName", enclaveName)
			err := cmd.Run()
			gomega.Expect(err).To(gomega.HaveOccurred())
			dive.Clean(enclaveName)
		})

		ginkgo.It("should handle invalid input for chain command", func() {
			enclaveName := dive.GenerateRandomName()
			dive.RunIconNode(enclaveName)
			serviceName0, endpoint0, nid0 := dive.GetServiceDetails(fmt.Sprintf("services_%s.json", enclaveName), dive.ICON_CONFIG0_SERVICENAME)
			cmd.Args = append(cmd.Args, "chain", "icon", "decentralize", "-p", "gochain", "-k", "keystores/invalid.json", "-n", nid0, "-e", endpoint0, "-s", serviceName0, "--enclaveName", enclaveName)
			err := cmd.Run()
			gomega.Expect(err).To(gomega.HaveOccurred())
			dive.Clean(enclaveName)
		})

		ginkgo.It("should handle invalid input for chain command", func() {
			enclaveName := dive.GenerateRandomName()
			dive.RunIconNode(enclaveName)
			serviceName0, endpoint0, _ := dive.GetServiceDetails(fmt.Sprintf("services_%s.json", enclaveName), dive.ICON_CONFIG0_SERVICENAME)
			cmd.Args = append(cmd.Args, "chain", "icon", "decentralize", "-p", "gochain", "-k", "keystores/keystore.json", "-n", "0x9", "-e", endpoint0, "-s", serviceName0, "--enclaveName", enclaveName)
			err := cmd.Run()
			gomega.Expect(err).To(gomega.HaveOccurred())
			dive.Clean(enclaveName)
		})

		ginkgo.It("should handle invalid input for chain command", func() {
			enclaveName := dive.GenerateRandomName()
			dive.RunIconNode(enclaveName)
			serviceName0, _, nid0 := dive.GetServiceDetails(fmt.Sprintf("services_%s.json", enclaveName), dive.ICON_CONFIG0_SERVICENAME)
			cmd.Args = append(cmd.Args, "chain", "icon", "decentralize", "-p", "gochain", "-k", "keystores/keystore.json", "-n", nid0, "-e", "http://172.16.0.3:9081/api/v3/icon_dex", "-s", serviceName0, "--enclaveName", enclaveName)
			err := cmd.Run()
			gomega.Expect(err).To(gomega.HaveOccurred())
			dive.Clean(enclaveName)
		})

		ginkgo.It("should handle invalid input for chain command", func() {
			enclaveName := dive.GenerateRandomName()
			dive.RunIconNode(enclaveName)
			_, endpoint0, nid0 := dive.GetServiceDetails(fmt.Sprintf("services_%s.json", enclaveName), dive.ICON_CONFIG0_SERVICENAME)
			cmd.Args = append(cmd.Args, "chain", "icon", "decentralize", "-p", "gochain", "-k", "keystores/keystore.json", "-n", nid0, "-e", endpoint0, "-s", "icon-node", "--enclaveName", enclaveName)
			err := cmd.Run()
			gomega.Expect(err).To(gomega.HaveOccurred())
			dive.Clean(enclaveName)
		})

		ginkgo.It("should output user that chain is already running when trying to run icon chain that is already running", func() {
			enclaveName := dive.GenerateRandomName()
			dive.RunIconNode(enclaveName)
			cmd.Args = append(cmd.Args, "chain", "icon", "--enclaveName", enclaveName)
			err := cmd.Run()
			gomega.Expect(err).To(gomega.HaveOccurred())
			dive.Clean(enclaveName)
		})

	})

	ginkgo.Describe("Eth chain commands", func() {
		ginkgo.It("should run single eth node", func() {
			enclaveName := dive.GenerateRandomName()
			cmd.Args = append(cmd.Args, "chain", "eth", "--enclaveName", enclaveName)
			err := cmd.Run()
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
			dive.Clean(enclaveName)
		})

		ginkgo.It("should output user that chain is already running when trying to run eth chain that is already running", func() {
			enclaveName := dive.GenerateRandomName()
			dive.RunEthNode(enclaveName)
			cmd.Args = append(cmd.Args, "chain", "eth", "--enclaveName", enclaveName)
			err := cmd.Run()
			gomega.Expect(err).To(gomega.HaveOccurred())
			dive.Clean(enclaveName)
		})

	})

	ginkgo.Describe("Hardhat chain commands", func() {
		ginkgo.It("should run single hardhat node-1", func() {
			enclaveName := dive.GenerateRandomName()
			cmd.Args = append(cmd.Args, "chain", "hardhat", "--enclaveName", enclaveName)
			err := cmd.Run()
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
			dive.Clean(enclaveName)
		})

		ginkgo.It("should output user that chain is already running when trying to run hardhat chain that is already running", func() {
			enclaveName := dive.GenerateRandomName()
			dive.RunHardhatNode(enclaveName)
			cmd.Args = append(cmd.Args, "chain", "hardhat", "--enclaveName", enclaveName)
			err := cmd.Run()
			gomega.Expect(err).To(gomega.HaveOccurred())
			dive.Clean(enclaveName)
		})
	})

	ginkgo.Describe("Archway chain commands", func() {
		ginkgo.It("should run single archway node", func() {
			enclaveName := dive.GenerateRandomName()
			cmd.Args = append(cmd.Args, "chain", "archway", "--enclaveName", enclaveName)
			err := cmd.Run()
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
			dive.Clean(enclaveName)

		})

		ginkgo.It("should run single custom archway node-1", func() {
			filepath1 := "../../cli/sample-jsons/archway.json"
			updated_path1 := dive.UpdatePublicPorts(filepath1)
			enclaveName := dive.GenerateRandomName()
			cmd.Args = append(cmd.Args, "chain", "archway", "-c", updated_path1, "--enclaveName", enclaveName)
			err := cmd.Run()
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
			defer os.Remove(updated_path1)
			dive.Clean(enclaveName)

		})

		ginkgo.It("should run single custom archway node with invalid json path", func() {
			enclaveName := dive.GenerateRandomName()
			cmd.Args = append(cmd.Args, "chain", "archway", "-c", "../../cli/sample-jsons/invalid_archway.json", "--enclaveName", enclaveName)
			err := cmd.Run()
			gomega.Expect(err).To(gomega.HaveOccurred())
			dive.Clean(enclaveName)

		})

		ginkgo.It("should output user that chain is already running when trying to run archway chain that is already running", func() {
			enclaveName := dive.GenerateRandomName()
			dive.RunArchwayNode(enclaveName)
			cmd.Args = append(cmd.Args, "chain", "archway", "--enclaveName", enclaveName)
			err := cmd.Run()
			gomega.Expect(err).To(gomega.HaveOccurred())
			dive.Clean(enclaveName)
		})
	})

	ginkgo.Describe("Neutron chain commands", func() {
		ginkgo.It("should run single neutron node", func() {
			enclaveName := dive.GenerateRandomName()
			cmd.Args = append(cmd.Args, "chain", "neutron", "--enclaveName", enclaveName)
			err := cmd.Run()
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
			dive.Clean(enclaveName)

		})

		ginkgo.It("should run single custom neutron node-1", func() {
			updated_path2 := dive.UpdateNeutronPublicPorts(dive.NEUTRON_CONFIG0)
			enclaveName := dive.GenerateRandomName()
			cmd.Args = append(cmd.Args, "chain", "neutron", "-c", updated_path2, "--enclaveName", enclaveName)
			defer os.Remove(updated_path2)
			err := cmd.Run()
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
			dive.Clean(enclaveName)
		})

		ginkgo.It("should run single custom neutron node with invalid json path", func() {
			enclaveName := dive.GenerateRandomName()
			cmd.Args = append(cmd.Args, "chain", "neutron", "-c", "../../cli/sample-jsons/neutron5.json", "--enclaveName", enclaveName)
			err := cmd.Run()
			gomega.Expect(err).To(gomega.HaveOccurred())
			dive.Clean(enclaveName)
		})

		ginkgo.It("should output user that chain is already running when trying to run neutron chain that is already running", func() {
			enclaveName := dive.GenerateRandomName()
			dive.RunNeutronNode(enclaveName)
			cmd.Args = append(cmd.Args, "chain", "neutron", "--enclaveName", enclaveName)
			err := cmd.Run()
			gomega.Expect(err).To(gomega.HaveOccurred())
			dive.Clean(enclaveName)
		})
	})

	ginkgo.Describe("Relaychain commands", func() {
		var selectedChain string
		if envChain := os.Getenv("relayChain"); envChain != "" {
			selectedChain = envChain
		} else {
			selectedChain = "default" // Provide a default value if not set
		}
		relayChainNames := []string{"kusama", "polkadot"}

		// Add a flag to check if the selected chain is valid
		validChainSelected := false

		for _, relayChainName := range relayChainNames {

			relayChainName := relayChainName // Capture the loop variable

			if selectedChain != "default" && selectedChain != relayChainName {
				// Skip tests for other chains
				continue
			}

			// Set the flag to indicate that a valid chain is selected
			validChainSelected = true

			ginkgo.It("should run single relaychain "+relayChainName, func() {
				enclaveName := dive.GenerateRandomName()
				cmd.Args = append(cmd.Args, "chain", relayChainName, "--enclaveName", enclaveName)
				err := cmd.Run()
				gomega.Expect(err).NotTo(gomega.HaveOccurred())
				dive.Clean(enclaveName)
			})

			ginkgo.It("should run single relaychain in mainnet for "+relayChainName, func() {
				enclaveName := dive.GenerateRandomName()
				cmd.Args = append(cmd.Args, "chain", relayChainName, "-n", "mainnet", "--enclaveName", enclaveName)
				err := cmd.Run()
				gomega.Expect(err).NotTo(gomega.HaveOccurred())
				dive.Clean(enclaveName)
			})

			ginkgo.It("should run single relaychain in testnet for "+relayChainName, func() {
				enclaveName := dive.GenerateRandomName()
				cmd.Args = append(cmd.Args, "chain", relayChainName, "-n", "testnet", "--enclaveName", enclaveName)
				err := cmd.Run()
				gomega.Expect(err).NotTo(gomega.HaveOccurred())
				dive.Clean(enclaveName)
			})

			ginkgo.It("should run custom relaychain in localnet for "+relayChainName, func() {
				enclaveName := dive.GenerateRandomName()
				config := dive.UpdateRelayChain(dive.LOCAL_CONFIG0, "local", "rococo-local", enclaveName, false, false)
				cmd.Args = append(cmd.Args, "chain", relayChainName, "-c", config, "--enclaveName", enclaveName)
				defer os.Remove(config)
				err := cmd.Run()
				gomega.Expect(err).NotTo(gomega.HaveOccurred())
				dive.Clean(enclaveName)
			})

			ginkgo.It("should run custom relaychain in testnet "+relayChainName, func() {
				enclaveName := dive.GenerateRandomName()
				config := dive.UpdateRelayChain(dive.LOCAL_CONFIG0, "testnet", "rococo", enclaveName, false, false)
				cmd.Args = append(cmd.Args, "chain", relayChainName, "-c", config, "--enclaveName", enclaveName)
				defer os.Remove(config)
				err := cmd.Run()
				gomega.Expect(err).NotTo(gomega.HaveOccurred())
				dive.Clean(enclaveName)
			})

			ginkgo.It("should run custom relaychain in mainnet"+relayChainName, func() {
				enclaveName := dive.GenerateRandomName()
				config := dive.UpdateRelayChain(dive.LOCAL_CONFIG0, "mainnet", "kusama", enclaveName, false, false)
				cmd.Args = append(cmd.Args, "chain", relayChainName, "-c", config, "--enclaveName", enclaveName)
				defer os.Remove(config)
				err := cmd.Run()
				gomega.Expect(err).NotTo(gomega.HaveOccurred())
				dive.Clean(enclaveName)
			})

			ginkgo.It("should run single relaychain with explorer and metrix service for "+relayChainName, func() {
				enclaveName := dive.GenerateRandomName()
				cmd.Args = append(cmd.Args, "chain", relayChainName, "--explorer", "--metrics", "--enclaveName", enclaveName)
				err := cmd.Run()
				gomega.Expect(err).NotTo(gomega.HaveOccurred())
				dive.Clean(enclaveName)
			})

			ginkgo.It("should run custom relaychain with explorer services and metrics service in testnet for "+relayChainName, func() {
				enclaveName := dive.GenerateRandomName()
				config := dive.UpdateRelayChain(dive.LOCAL_CONFIG0, "testnet", "rococo", enclaveName, true, true)
				cmd.Args = append(cmd.Args, "chain", relayChainName, "-c", config, "--enclaveName", enclaveName)
				defer os.Remove(config)
				err := cmd.Run()
				gomega.Expect(err).NotTo(gomega.HaveOccurred())
				dive.Clean(enclaveName)
			})

		}
		if !validChainSelected && selectedChain != "default" {
			// Print an error message if an invalid chain is selected
			fmt.Printf("Error: Invalid relayChain selected: %s. Expected 'kusama' or 'polkadot'\n", selectedChain)
			log.Fatal("Tests cannot be run because an invalid relayChain is selected.")
		}
	})

	ginkgo.Describe("Parachain commands", func() {
		var selectedRelayChain string
		var selectedParaChain string

		if envChain := os.Getenv("relayChain"); envChain != "" {
			selectedRelayChain = envChain
		} else {
			selectedRelayChain = "default" // Provide a default value if not set
		}

		if envSelectedParaChain := os.Getenv("paraChain"); envSelectedParaChain != "" {
			selectedParaChain = envSelectedParaChain
		} else {
			// If not provided, set a default value
			selectedParaChain = "default"
		}

		relayChainNames := []string{"kusama", "polkadot"}

		if selectedRelayChain != "default" {
			validRelayChain := false
			for _, relayChainName := range relayChainNames {
				if selectedRelayChain == relayChainName {
					validRelayChain = true
					break
				}
			}
			if !validRelayChain {
				fmt.Printf("Error: Invalid relayChain selected: %s. Expected one of %v\n", selectedRelayChain, relayChainNames)
				log.Fatal("Tests cannot be run because an invalid relayChain is selected.")
			}
		}

		for _, relayChainName := range relayChainNames {
			if selectedRelayChain != "default" && selectedRelayChain != relayChainName {
				continue
			}

			relayChainName := relayChainName // Capture the loop variable

			var paraChainNames []string

			if relayChainName == "kusama" {
				paraChainNames = []string{"karura", "kintsugi-btc", "altair", "bifrost", "mangata", "robonomics", "turing-network", "encointer-network", "bajun-networkc", "calamari", "Khala-network", "litmus", "moonriver", "subzero"}
			} else if relayChainName == "polkadot" {
				paraChainNames = []string{"polkadex", "zeitgeist", "acala", "bifrost", "clover", "integritee-shell", "integritee-shell", "litentry", "moonbeam", "nodle", "pendulum", "ajuna-network", "centrifuge", "frequency", "interlay", "kylin", "manta", "moonsama", "parallel", "phala-network", "subsocial"}
			} else {
				paraChainNames = []string{"polkadex", "zeitgeist", "karura", "kintsugi-btc", "altair", "bifrost", "mangata", "robonomics", "turing-network", "encointer-network", "bajun-networkc", "calamari", "khala-network", "litmus", "moonriver", "subzero", "acala", "bifrost", "clover", "integritee-shell", "integritee-shell", "litentry", "moonbeam", "nodle", "pendulum", "ajuna-network", "centrifuge", "frequency", "interlay", "kylin", "manta", "moonsama", "parallel", "phala-network", "subsocial"}
			}

			// Validate paraChain before running tests
			if selectedParaChain != "default" {
				validParaChain := false
				for _, paraChainName := range paraChainNames {
					if selectedParaChain == paraChainName {
						validParaChain = true
						break
					}
				}
				if !validParaChain {
					fmt.Printf("Error: Invalid paraChain selected: %s. Expected a valid parachain name\n", selectedParaChain)
					log.Fatal("Tests cannot be run because an invalid paraChain is selected.")
				}
			}

			for _, paraChainName := range paraChainNames {
				if selectedParaChain != "default" && selectedParaChain != paraChainName {
					continue
				}

				paraChainName := paraChainName

				ginkgo.It("should run single parachain  in testnet for "+relayChainName+" and "+paraChainName, func() {
					enclaveName := dive.GenerateRandomName()
					cmd.Args = append(cmd.Args, "chain", relayChainName, "-p", paraChainName, "--no-relay", "-n", "testnet", "--enclaveName", enclaveName)
					err := cmd.Run()
					gomega.Expect(err).NotTo(gomega.HaveOccurred())
					dive.Clean(enclaveName)
				})
				ginkgo.It("should run single parachain in mainnet for "+relayChainName+" and "+paraChainName, func() {
					enclaveName := dive.GenerateRandomName()
					cmd.Args = append(cmd.Args, "chain", relayChainName, "-p", paraChainName, "--no-relay", "-n", "mainnet", "--enclaveName", enclaveName)
					err := cmd.Run()
					gomega.Expect(err).NotTo(gomega.HaveOccurred())
					dive.Clean(enclaveName)
				})
				ginkgo.It("should run single parachain in mainnet with explorer services for "+relayChainName+" and "+paraChainName, func() {
					enclaveName := dive.GenerateRandomName()
					cmd.Args = append(cmd.Args, "chain", relayChainName, "-p", paraChainName, "--no-relay", "-n", "mainnet", "--explorer", "--enclaveName", enclaveName)
					err := cmd.Run()
					gomega.Expect(err).NotTo(gomega.HaveOccurred())
					dive.Clean(enclaveName)
				})
				ginkgo.It("should run single parachain in mainnet with metrics services for "+relayChainName+" and "+paraChainName, func() {
					enclaveName := dive.GenerateRandomName()
					cmd.Args = append(cmd.Args, "chain", relayChainName, "-p", paraChainName, "--no-relay", "-n", "mainnet", "--metrics", "--enclaveName", enclaveName)
					err := cmd.Run()
					gomega.Expect(err).NotTo(gomega.HaveOccurred())
					dive.Clean(enclaveName)
				})
				ginkgo.It("should run custom parachain in testnet with for "+relayChainName+" and "+paraChainName, func() {
					enclaveName := dive.GenerateRandomName()
					config := dive.UpdateParaChain(dive.LOCAL_CONFIG0, "karura", false, false)
					cmd.Args = append(cmd.Args, "chain", relayChainName, "--no-relay", "-c", config, "--enclaveName", enclaveName)
					defer os.Remove(config)
					err := cmd.Run()
					gomega.Expect(err).NotTo(gomega.HaveOccurred())
					dive.Clean(enclaveName)
				})
				ginkgo.It("should run custom parachain in mainnet for "+relayChainName+" and "+paraChainName, func() {
					enclaveName := dive.GenerateRandomName()
					config := dive.UpdateParaChain(dive.LOCAL_CONFIG0, "karura", false, false)
					cmd.Args = append(cmd.Args, "chain", relayChainName, "--no-relay", "-c", config, "--enclaveName", enclaveName)
					defer os.Remove(config)
					err := cmd.Run()
					gomega.Expect(err).NotTo(gomega.HaveOccurred())
					dive.Clean(enclaveName)
				})
				ginkgo.It("should run custom parachain in mainnet with explorer services for "+relayChainName+" and "+paraChainName, func() {
					enclaveName := dive.GenerateRandomName()
					config := dive.UpdateParaChain(dive.LOCAL_CONFIG0, "karura", true, false)
					cmd.Args = append(cmd.Args, "chain", relayChainName, "--no-relay", "-c", config, "--enclaveName", enclaveName)
					defer os.Remove(config)
					err := cmd.Run()
					gomega.Expect(err).NotTo(gomega.HaveOccurred())
					dive.Clean(enclaveName)
				})
				ginkgo.It("should run custom parachain in mainnet with metrics services for "+relayChainName+" and "+paraChainName, func() {
					enclaveName := dive.GenerateRandomName()
					config := dive.UpdateParaChain(dive.LOCAL_CONFIG0, "karura", false, true)
					cmd.Args = append(cmd.Args, "chain", relayChainName, "--no-relay", "-c", config, "--enclaveName", enclaveName)
					defer os.Remove(config)
					err := cmd.Run()
					gomega.Expect(err).NotTo(gomega.HaveOccurred())
					dive.Clean(enclaveName)
				})
				ginkgo.It("should run custom parachain in mainnet with explorer and metrics services for "+relayChainName+" and "+paraChainName, func() {
					enclaveName := dive.GenerateRandomName()
					config := dive.UpdateParaChain(dive.LOCAL_CONFIG0, "karura", true, true)
					cmd.Args = append(cmd.Args, "chain", relayChainName, "--no-relay", "-c", config, "--enclaveName", enclaveName)
					defer os.Remove(config)
					err := cmd.Run()
					gomega.Expect(err).NotTo(gomega.HaveOccurred())
					dive.Clean(enclaveName)
				})
			}
		}
	})

	ginkgo.Describe("Relaychain and parachain commands", func() {

		var selectedRelayChain string
		var selectedParaChain string

		if envChain := os.Getenv("relayChain"); envChain != "" {
			selectedRelayChain = envChain
		} else {
			selectedRelayChain = "default" // Provide a default value if not set
		}

		if envSelectedParaChain := os.Getenv("paraChain"); envSelectedParaChain != "" {
			selectedParaChain = envSelectedParaChain
		} else {
			// If not provided, set a default value
			selectedParaChain = "default"
		}

		relayChainNames := []string{"kusama", "polkadot"}

		// Validate relayChain before running tests
		if selectedRelayChain != "default" {
			validRelayChain := false
			for _, relayChainName := range relayChainNames {
				if selectedRelayChain == relayChainName {
					validRelayChain = true
					break
				}
			}
			if !validRelayChain {
				fmt.Printf("Error: Invalid relayChain selected: %s. Expected one of %v\n", selectedRelayChain, relayChainNames)
				log.Fatal("Tests cannot be run because an invalid relayChain is selected.")
			}
		}

		// if selectedParaChain == "default" && selectedRelayChain == "default" {
		// 	fmt.Println("Error: Atleast relay chain should be given. ")
		// 	log.Fatal("Tests cannot be run because relayChain is missing.")
		// 	return // Added return to stop further execution
		// }

		for _, relayChainName := range relayChainNames {
			if selectedRelayChain != "default" && selectedRelayChain != relayChainName {
				continue

			}

			relayChainName := relayChainName // Capture the loop variable

			var paraChainNames []string

			if relayChainName == "kusama" {
				paraChainNames = []string{"karura", "kintsugi-btc", "altair", "bifrost", "mangata", "robonomics", "turing-network", "encointer-Network", "bajun-networkc", "calamari", "khala-network", "litmus", "moonriver", "subzero"}
			} else if relayChainName == "polkadot" {
				paraChainNames = []string{"Polkadex", "zeitgeist", "acala", "bifrost", "clover", "integritee-shell", "integritee-shell", "litentry", "moonbeam", "nodle", "pendulum", "ajuna-network", "centrifuge", "frequency", "interlay", "kylin", "manta", "moonsama", "parallel", "phala-network", "subsocial"}
			} else {
				paraChainNames = []string{"polkadex", "zeitgeist", "karura", "kintsugi-btc", "altair", "bifrost", "mangata", "robonomics", "turing-network", "encointer-network", "bajun-networkc", "calamari", "khala-network", "litmus", "moonriver", "subzero", "acala", "bifrost", "clover", "integritee-shell", "integritee-shell", "litentry", "moonbeam", "nodle", "pendulum", "ajuna-network", "centrifuge", "frequency", "interlay", "kylin", "manta", "moonsama", "parallel", "phala-network", "subsocial"}
			}

			// Validate paraChain before running tests
			if selectedParaChain != "default" {
				validParaChain := false
				for _, paraChainName := range paraChainNames {
					if selectedParaChain == paraChainName {
						validParaChain = true
						break
					}
				}
				if !validParaChain {
					fmt.Printf("Error: Invalid paraChain selected: %s. Expected a valid parachain name\n", selectedParaChain)
					log.Fatal("Tests cannot be run because an invalid paraChain is selected.")
				}
			}

			for _, paraChainName := range paraChainNames {
				if selectedParaChain != "default" && selectedParaChain != paraChainName {
					continue
				}

				paraChainName := paraChainName
				ginkgo.It("should run single relaychain and parachain  in testnet for "+relayChainName+" and "+paraChainName, func() {
					enclaveName := dive.GenerateRandomName()
					cmd.Args = append(cmd.Args, "chain", relayChainName, "-p", paraChainName, "-n", "testnet", "--enclaveName", enclaveName)
					fmt.Println(cmd)
					err := cmd.Run()
					gomega.Expect(err).NotTo(gomega.HaveOccurred())
					dive.Clean(enclaveName)
				})
				ginkgo.It("should run single relaychain and parachain in mainnet for "+relayChainName+" and "+paraChainName, func() {
					enclaveName := dive.GenerateRandomName()
					cmd.Args = append(cmd.Args, "chain", relayChainName, "-p", paraChainName, "-n", "mainnet", "--enclaveName", enclaveName)
					err := cmd.Run()
					gomega.Expect(err).NotTo(gomega.HaveOccurred())
					dive.Clean(enclaveName)
				})
				ginkgo.It("should run single relaychain and parachain in local for "+relayChainName+" and "+paraChainName, func() {
					enclaveName := dive.GenerateRandomName()
					cmd.Args = append(cmd.Args, "chain", relayChainName, "-p", paraChainName, "-n", "local", "--enclaveName", enclaveName)
					err := cmd.Run()
					gomega.Expect(err).NotTo(gomega.HaveOccurred())
					dive.Clean(enclaveName)
				})
				ginkgo.It("should run single relaychain and parachain in mainnet with explorer services for "+relayChainName+" and "+paraChainName, func() {
					enclaveName := dive.GenerateRandomName()
					cmd.Args = append(cmd.Args, "chain", relayChainName, "-p", paraChainName, "-n", "mainnet", "--explorer", "--enclaveName", enclaveName)
					err := cmd.Run()
					gomega.Expect(err).NotTo(gomega.HaveOccurred())
					dive.Clean(enclaveName)
				})
				ginkgo.It("should run single relaychain and  parachain in mainnet with metrics services for "+relayChainName+" and "+paraChainName, func() {
					enclaveName := dive.GenerateRandomName()
					cmd.Args = append(cmd.Args, "chain", relayChainName, "-p", paraChainName, "-n", "mainnet", "--metrics", "--enclaveName", enclaveName)
					err := cmd.Run()
					gomega.Expect(err).NotTo(gomega.HaveOccurred())
					dive.Clean(enclaveName)
				})
				ginkgo.It("should run custom relaychain and parachain in testnet for "+relayChainName+" and "+paraChainName, func() {
					enclaveName := dive.GenerateRandomName()
					config := dive.UpdateChainInfo(dive.LOCAL_CONFIG0, "testnet", "rococo", "karura", false, false)
					cmd.Args = append(cmd.Args, "chain", relayChainName, "-c", config, "--enclaveName", enclaveName)
					defer os.Remove(config)
					err := cmd.Run()
					gomega.Expect(err).NotTo(gomega.HaveOccurred())
					dive.Clean(enclaveName)
				})
				ginkgo.It("should run custom relaychain and  parachain in mainnet for "+relayChainName+" and "+paraChainName, func() {
					enclaveName := dive.GenerateRandomName()
					config := dive.UpdateChainInfo(dive.LOCAL_CONFIG0, "mainnet", "polkadot", "karura", false, false)
					cmd.Args = append(cmd.Args, "chain", relayChainName, "-c", config, "--enclaveName", enclaveName)
					defer os.Remove(config)
					err := cmd.Run()
					gomega.Expect(err).NotTo(gomega.HaveOccurred())
					dive.Clean(enclaveName)
				})
				ginkgo.It("should run custom relaychain and  parachain in local for "+relayChainName+" and "+paraChainName, func() {
					enclaveName := dive.GenerateRandomName()
					config := dive.UpdateChainInfo(dive.LOCAL_CONFIG0, "local", "rococo-local", "karura", false, false)
					cmd.Args = append(cmd.Args, "chain", relayChainName, "-c", config, "--enclaveName", enclaveName)
					defer os.Remove(config)
					err := cmd.Run()
					gomega.Expect(err).NotTo(gomega.HaveOccurred())
					dive.Clean(enclaveName)
				})
				ginkgo.It("should run custom relaychain and  parachain in mainnet  with explorer services for "+relayChainName+" and "+paraChainName, func() {
					enclaveName := dive.GenerateRandomName()
					config := dive.UpdateChainInfo(dive.LOCAL_CONFIG0, "mainnet", "polkadot", "frequency", true, false)
					cmd.Args = append(cmd.Args, "chain", relayChainName, "-c", config, "--enclaveName", enclaveName)
					defer os.Remove(config)
					err := cmd.Run()
					gomega.Expect(err).NotTo(gomega.HaveOccurred())
					dive.Clean(enclaveName)
				})
				ginkgo.It("should run custom relaychain and parachain in mainnet  with metrics services for "+relayChainName+" and "+paraChainName, func() {
					enclaveName := dive.GenerateRandomName()
					config := dive.UpdateChainInfo(dive.LOCAL_CONFIG0, "mainnet", "polkadot", "frequency", true, false)
					cmd.Args = append(cmd.Args, "chain", relayChainName, "-c", config, "--enclaveName", enclaveName)
					defer os.Remove(config)
					err := cmd.Run()
					gomega.Expect(err).NotTo(gomega.HaveOccurred())
					dive.Clean(enclaveName)
				})
				ginkgo.It("should run custom relaychain and parachain in mainnet with explorer and metrics services for "+relayChainName+" and "+paraChainName, func() {
					enclaveName := dive.GenerateRandomName()
					config := dive.UpdateChainInfo(dive.LOCAL_CONFIG0, "mainnet", "polkadot", "frequency", true, true)
					cmd.Args = append(cmd.Args, "chain", relayChainName, "-c", config, "--enclaveName", enclaveName)
					defer os.Remove(config)
					err := cmd.Run()
					gomega.Expect(err).NotTo(gomega.HaveOccurred())
					dive.Clean(enclaveName)
				})
			}
		}
	})
})
