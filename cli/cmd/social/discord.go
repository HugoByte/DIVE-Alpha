package social

import (
	"github.com/hugobyte/dive-core/cli/common"
	"github.com/spf13/cobra"
)

const diveURL = "https://discord.gg/GyRQSBN3Cu"

var DiscordCmd = common.NewDiveCommandBuilder().
	SetUse("discord").
	SetShort("Opens DIVE discord channel").
	SetLong(`The command opens the Discord channel for DIVE, providing a direct link or launching the Discord application to access the dedicated DIVE community. It allows users to engage in discussions, seek support, share insights, and collaborate with other members of the DIVE community within the Discord platform.`).
	SetRunE(discord).Build()

func discord(cmd *cobra.Command, args []string) error {

	cli := common.GetCli()

	err := common.ValidateArgs(args)

	if err != nil {
		cli.Logger().Fatal(common.CodeOf(err), err.Error())
		return err
	}

	cli.Logger().Info("Redirecting to DIVE discord channel...")

	if err := common.OpenFile(diveURL); err != nil {
		cli.Logger().Fatalf(common.CodeOf(err), "Failed to open Dive discord channel with error %v", err)
		return err
	}
	return nil
}
