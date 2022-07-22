package config

import (
	"errors"
	"fmt"
	"os"
	"strconv"

	"github.com/yuk7/wsldl/get"
	"github.com/yuk7/wsldl/lib/utils"
	"github.com/yuk7/wsldl/run"
	"github.com/yuk7/wsllib-go"
	wslreg "github.com/yuk7/wslreglib-go"
)

type Setting string

const (
	DefaultUID  Setting = "default-uid"
	DefaultUser Setting = "default-user"
	AppendPath  Setting = "append-path"
	MountDrive  Setting = "mount-drive"
	WSLVersion  Setting = "wsl-version"
	DefaultTerm Setting = "default-term"
	FlagsVal    Setting = "flags-val"
)

//Execute is default install entrypoint
func Execute(name string, setting *Setting, value interface{}) {
	var err error
	ok := true
	uid, flags := get.WslGetConfig(name)

	if setting != nil && value != nil {
		switch *setting {
		case DefaultUID:
			var intUID uint64
			intUID, ok = value.(uint64)

			if ok {
				uid = uint64(intUID)
			}

		case DefaultUser:
			user := fmt.Sprintf("%v", value)

			str, _, errtmp := run.ExecRead(name, "id -u "+utils.DQEscapeString(user))
			err = errtmp
			if err == nil {
				var intUID int
				intUID, err = strconv.Atoi(str)
				uid = uint64(intUID)
				if err != nil {
					err = errors.New(str)
				}
			}

		case AppendPath:
			var b bool
			arg := fmt.Sprintf("%v", value)

			b, err = strconv.ParseBool(arg)
			if b {
				flags |= wsllib.FlagAppendNTPath
			} else {
				flags ^= wsllib.FlagAppendNTPath
			}

		case MountDrive:
			var b bool
			arg := fmt.Sprintf("%v", value)

			b, err = strconv.ParseBool(arg)
			if b {
				flags |= wsllib.FlagEnableDriveMounting
			} else {
				flags ^= wsllib.FlagEnableDriveMounting
			}

		case WSLVersion:
			var intWslVer int
			intWslVer, ok = value.(int)

			if ok {
				if intWslVer == 1 || intWslVer == 2 {
					err = wslreg.SetWslVersion(name, intWslVer)
				} else {
					err = os.ErrInvalid
					break
				}
			}

		case DefaultTerm:
			termValue := 0
			term := fmt.Sprintf("%v", value)

			switch term {
			case "default", strconv.Itoa(wslreg.FlagWsldlTermDefault):
				termValue = wslreg.FlagWsldlTermDefault
			case "wt", strconv.Itoa(wslreg.FlagWsldlTermWT):
				termValue = wslreg.FlagWsldlTermWT
			case "flute", strconv.Itoa(wslreg.FlagWsldlTermFlute):
				termValue = wslreg.FlagWsldlTermFlute
			default:
				err = os.ErrInvalid
			}

			if err == nil {
				break
			}

			profile, tempErr := wslreg.GetProfileFromName(name)
			if tempErr != nil {
				err = tempErr
				break
			}
			profile.WsldlTerm = termValue
			err = wslreg.WriteProfile(profile)

		case FlagsVal:
			var intFlags int
			intFlags, ok = value.(int)

			if ok {
				flags = uint32(intFlags)
			}

		default:
			err = os.ErrInvalid
		}

		if !ok {
			err = os.ErrInvalid
		}

		if err != nil {
			utils.ErrorExit(err, true, true, false)
		}

		wsllib.WslConfigureDistribution(name, uid, flags)
	} else {
		utils.ErrorExit(os.ErrInvalid, true, true, false)
	}
}
