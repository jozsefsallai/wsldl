package main

import (
	"os"
	"path/filepath"

	"github.com/urfave/cli/v2"
	"github.com/yuk7/wsldl/backup"
	"github.com/yuk7/wsldl/clean"
	"github.com/yuk7/wsldl/config"
	"github.com/yuk7/wsldl/get"
	"github.com/yuk7/wsldl/install"
	"github.com/yuk7/wsldl/isregd"
	"github.com/yuk7/wsldl/lib/utils"
	"github.com/yuk7/wsldl/run"
	"github.com/yuk7/wsldl/version"
	"github.com/yuk7/wsllib-go"
)

func main() {
	efPath, _ := os.Executable()
	name := filepath.Base(efPath[:len(efPath)-len(filepath.Ext(efPath))])

	app := &cli.App{
		Name:  "wsldl",
		Usage: "Advanced WSL Distribution Launcher / Installer",
		Action: func(*cli.Context) error {
			if !wsllib.WslIsDistributionRegistered(name) {
				install.Execute(name, nil)
			} else {
				run.ExecuteNoArgs(name)
			}

			return nil
		},
		Commands: []*cli.Command{
			{
				Name:    "version",
				Aliases: []string{"-v", "--version"},
				Usage:   "Print version",
				Action: func(*cli.Context) error {
					version.Execute()
					return nil
				},
			},

			{
				Name:  "isregd",
				Usage: "Check if default distribution is registered via exit code",
				Action: func(*cli.Context) error {
					isregd.Execute(name)
					return nil
				},
			},

			{
				Name:  "install",
				Usage: "Install WSL distribution",
				Flags: []cli.Flag{
					&cli.BoolFlag{
						Name:  "root",
						Usage: "Detect rootfs files",
					},
				},
				Action: func(ctx *cli.Context) error {
					detectRootfsFiles := ctx.Bool("root")
					backupLocation := ctx.Args().First()

					install.Execute(name, &install.Options{
						DetectRootfsFiles: detectRootfsFiles,
						BackupLocation:    backupLocation,
					})

					return nil
				},
			},

			{
				Name:    "run",
				Aliases: []string{"-c", "/c"},
				Usage:   "Run the given command line in that instance. Inherit current directory.",
				Action: func(ctx *cli.Context) error {
					args := ctx.Args().Slice()
					run.Execute(name, args)
					return nil
				},
			},

			{
				Name:    "runp",
				Aliases: []string{"-p", "/p"},
				Usage:   "Run the given command line in that instance after converting its path.",
				Action: func(ctx *cli.Context) error {
					args := ctx.Args().Slice()
					run.ExecuteP(name, args)
					return nil
				},
			},

			{
				Name:    "config",
				Aliases: []string{"set"},
				Usage:   "Set config value",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:  string(config.DefaultUser),
						Usage: "Set the default user of this instance to `USER`",
					},
					&cli.Uint64Flag{
						Name:  string(config.DefaultUID),
						Usage: "Set the default user uid of this instance to `UID`",
					},
					&cli.StringFlag{
						Name:  string(config.AppendPath),
						Usage: "Switch of Append Windows PATH to $PATH",
					},
					&cli.StringFlag{
						Name:  string(config.MountDrive),
						Usage: "Switch of Mount drives",
					},
					&cli.IntFlag{
						Name:  string(config.WSLVersion),
						Usage: "Set the WSL version of this instance to `1 or 2`",
					},
					&cli.StringFlag{
						Name:  string(config.DefaultTerm),
						Usage: "Set default type of terminal window.",
					},
					&cli.IntFlag{
						Name:  string(config.FlagsVal),
						Usage: "Set the flag value directly.",
					},
				},
				Action: func(ctx *cli.Context) error {
					flags := ctx.LocalFlagNames()
					if len(flags) == 0 {
						config.Execute(name, nil, nil)
					}

					flag := flags[0]
					setting := config.Setting(flag)
					value := ctx.Generic(flag)

					config.Execute(name, &setting, value)
					return nil
				},
			},

			{
				Name:  "get",
				Usage: "Get config value",
				Flags: []cli.Flag{
					&cli.BoolFlag{
						Name: string(get.DefaultUID),
					},
					&cli.BoolFlag{
						Name: string(get.AppendPath),
					},
					&cli.BoolFlag{
						Name: string(get.MountDrive),
					},
					&cli.BoolFlag{
						Name: string(get.WSLVersion),
					},
					&cli.BoolFlag{
						Name:    string(get.LXGUID),
						Aliases: []string{string(get.LXUID)},
					},
					&cli.BoolFlag{
						Name:    string(get.DefaultTerm),
						Aliases: []string{string(get.DefaultTerminal)},
					},
					&cli.BoolFlag{
						Name:    string(get.WTProfileName),
						Aliases: []string{string(get.WTProfileName2), string(get.WTProfileName3)},
					},
					&cli.BoolFlag{
						Name: string(get.FlagsVal),
					},
					&cli.BoolFlag{
						Name: string(get.FlagsBits),
					},
				},
				Action: func(ctx *cli.Context) error {
					flags := ctx.LocalFlagNames()
					if len(flags) == 0 {
						get.Execute(name, nil)
					}

					flag := flags[0]
					setting := get.Setting(flag)

					get.Execute(name, &setting)
					return nil
				},
			},

			{
				Name:  "backup",
				Usage: "Backup WSL instance",
				Flags: []cli.Flag{
					&cli.BoolFlag{
						Name:  "tar",
						Usage: "Backup as tar file",
					},
					&cli.BoolFlag{
						Name:  "tgz",
						Usage: "Backup as tar.gz file",
					},
					&cli.BoolFlag{
						Name:  "vhdx",
						Usage: "Backup as vhdx file (WSL2 only)",
					},
					&cli.BoolFlag{
						Name:  "vhdxgz",
						Usage: "Backup as vhdx.gz file (WSL2 only)",
					},
					&cli.BoolFlag{
						Name:  "reg",
						Usage: "Backup settings registry file",
					},
				},
				Action: func(ctx *cli.Context) error {
					backup.Execute(name, &backup.Options{
						Tar:    ctx.Bool("tar"),
						Tgz:    ctx.Bool("tgz"),
						Vhdx:   ctx.Bool("vhdx"),
						Vhdxgz: ctx.Bool("vhdxgz"),
						Reg:    ctx.Bool("reg"),
					})

					return nil
				},
			},

			{
				Name:  "clean",
				Usage: "Uninstall WSL instance",
				Flags: []cli.Flag{
					&cli.BoolFlag{
						Name:  "y",
						Usage: "Confirm uninstallation of instance",
					},
				},
				Action: func(ctx *cli.Context) error {
					clean.Execute(name, ctx.Bool("y"))
					return nil
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		utils.ErrorExit(err, true, false, false)
	}
}
