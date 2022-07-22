package get

import (
	"errors"
	"fmt"
	"os"

	"github.com/yuk7/wsldl/lib/utils"
	"github.com/yuk7/wsldl/lib/wtutils"
	"github.com/yuk7/wsllib-go"
	wslreg "github.com/yuk7/wslreglib-go"
)

type Setting string

const (
	DefaultUID Setting = "default-uid"
	AppendPath Setting = "append-path"
	MountDrive Setting = "mount-drive"
	WSLVersion Setting = "wsl-version"
	LXGUID     Setting = "lxguid"
	LXUID      Setting = "lxuid"

	DefaultTerm     Setting = "default-term"
	DefaultTerminal Setting = "default-terminal"

	WTProfileName  Setting = "wt-profile-name"
	WTProfileName2 Setting = "wt-profilename"
	WTProfileName3 Setting = "wt-pn"

	FlagsVal  Setting = "flags-val"
	FlagsBits Setting = "flags-bits"
)

//Execute is default install entrypoint
func Execute(name string, setting *Setting) {
	uid, flags := WslGetConfig(name)
	profile, proferr := wslreg.GetProfileFromName(name)
	if setting != nil {
		switch *setting {
		case DefaultUID:
			print(uid)

		case AppendPath:
			print(flags&wsllib.FlagAppendNTPath == wsllib.FlagAppendNTPath)

		case MountDrive:
			print(flags&wsllib.FlagEnableDriveMounting == wsllib.FlagEnableDriveMounting)

		case WSLVersion:
			if flags&wsllib.FlagEnableWsl2 == wsllib.FlagEnableWsl2 {
				print("2")
			} else {
				print("1")
			}

		case LXGUID, LXUID:
			if profile.UUID == "" {
				if proferr != nil {
					utils.ErrorExit(proferr, true, true, false)
				}
				utils.ErrorExit(errors.New("lxguid get failed"), true, true, false)
			}
			print(profile.UUID)

		case DefaultTerm, DefaultTerminal:
			switch profile.WsldlTerm {
			case wslreg.FlagWsldlTermWT:
				print("wt")
			case wslreg.FlagWsldlTermFlute:
				print("flute")
			default:
				print("default")
			}

		case WTProfileName, WTProfileName2, WTProfileName3:
			if profile.DistributionName != "" {
				name = profile.DistributionName
			}

			conf, err := wtutils.ReadParseWTConfig()
			if err != nil {
				utils.ErrorExit(err, true, true, false)
			}
			guid := "{" + wtutils.CreateProfileGUID(name) + "}"
			profileName := ""
			for _, profile := range conf.Profiles.ProfileList {
				if profile.GUID == guid {
					profileName = profile.Name
					break
				}
			}
			if profileName != "" {
				print(profileName)
			} else {
				utils.ErrorExit(errors.New("profile not found"), true, true, false)
			}

		case FlagsVal:
			print(flags)

		case FlagsBits:
			fmt.Printf("%04b", flags)

		default:
			utils.ErrorExit(os.ErrInvalid, true, true, false)
		}
	} else {
		utils.ErrorExit(os.ErrInvalid, true, true, false)
	}
}
