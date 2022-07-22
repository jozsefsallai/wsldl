package install

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strconv"

	"github.com/yuk7/wsldl/lib/preset"
	"github.com/yuk7/wsldl/lib/utils"
	"github.com/yuk7/wsllib-go"
)

type Options struct {
	BackupLocation    string
	DetectRootfsFiles bool
}

//Execute is default install entrypoint
func Execute(name string, options *Options) {
	if !wsllib.WslIsDistributionRegistered(name) {
		var rootPath string
		var showProgress bool
		jsonPreset, _ := preset.ReadParsePreset()

		if options != nil {
			if len(options.BackupLocation) != 0 {
				rootPath = options.BackupLocation
			} else {
				rootPath = detectRootfsFiles()
				if jsonPreset.InstallFile != "" {
					rootPath = jsonPreset.InstallFile
				}

				showProgress = options.DetectRootfsFiles
			}
		}

		if options == nil {
			if isInstalledFilesExist() {
				var in string
				fmt.Printf("An old installation file was found.\n")
				fmt.Printf("Do you want to rewrite and repair the installation infomation?\n")
				fmt.Printf("Type y/n:")
				fmt.Scan(&in)

				if in == "y" {
					err := repairRegistry(name)
					if err != nil {
						utils.ErrorExit(err, showProgress, true, showProgress)
					}
					utils.StdoutGreenPrintln("done.")
					return
				}
			}
		}

		err := Install(name, rootPath, showProgress)
		if err != nil {
			utils.ErrorExit(err, showProgress, true, options == nil)
		}

		if jsonPreset.WslVersion == 1 || jsonPreset.WslVersion == 2 {
			wslexe := utils.GetWindowsDirectory() + "\\System32\\wsl.exe"
			_, err = exec.Command(wslexe, "--set-version", name, strconv.Itoa(jsonPreset.WslVersion)).Output()
		}

		if err == nil {
			if showProgress {
				utils.StdoutGreenPrintln("Installation complete")
			}
		} else {
			utils.ErrorExit(err, showProgress, true, options == nil)
		}

		if options == nil {
			fmt.Fprintf(os.Stdout, "Press enter to continue...")
			bufio.NewReader(os.Stdin).ReadString('\n')
		}

	} else {
		utils.ErrorRedPrintln("ERR: [" + name + "] is already installed.")
		utils.ErrorExit(os.ErrInvalid, false, true, false)
	}
}
