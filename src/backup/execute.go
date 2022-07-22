package backup

import (
	"os"

	"github.com/yuk7/wsldl/lib/utils"
	"github.com/yuk7/wsllib-go"
)

type Options struct {
	Tar    bool
	Tgz    bool
	Vhdx   bool
	Vhdxgz bool
	Reg    bool
}

//Execute is default run entrypoint.
func Execute(name string, opts *Options) {
	opttar := false
	opttgz := false
	optvhdx := false
	optvhdxgz := false
	optreg := false

	ok := false

	if opts != nil {
		opttar = opts.Tar
		opttgz = opts.Tgz
		optvhdx = opts.Vhdx
		optvhdxgz = opts.Vhdxgz
		optreg = opts.Reg

		ok = opttar || opttgz || optvhdx || optvhdxgz || optreg
	}

	if !ok {
		_, _, flags, _ := wsllib.WslGetDistributionConfiguration(name)
		if flags&wsllib.FlagEnableWsl2 == wsllib.FlagEnableWsl2 {
			optvhdxgz = true
			optreg = true
		} else {
			opttgz = true
			optreg = true
		}

		ok = true
	}

	if !ok {
		utils.ErrorExit(os.ErrInvalid, true, true, false)
	}

	if optreg {
		err := backupReg(name, "backup.reg")
		if err != nil {
			utils.ErrorExit(err, true, true, false)
		}
	}

	if opttar {
		err := backupTar(name, "backup.tar")
		if err != nil {
			utils.ErrorExit(err, true, true, false)
		}
	}

	if opttgz {
		err := backupTar(name, "backup.tar.gz")
		if err != nil {
			utils.ErrorExit(err, true, true, false)
		}
	}

	if optvhdx {
		err := backupExt4Vhdx(name, "backup.ext4.vhdx")
		if err != nil {
			utils.ErrorExit(err, true, true, false)
		}
	}

	if optvhdxgz {
		err := backupExt4Vhdx(name, "backup.ext4.vhdx.gz")
		if err != nil {
			utils.ErrorExit(err, true, true, false)
		}
	}
}
