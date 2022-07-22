package clean

import (
	"fmt"
	"os"

	"github.com/yuk7/wsldl/lib/utils"
)

//Execute is default run entrypoint.
func Execute(name string, confirm bool) {
	showProgress := true
	if !confirm {
		var in string
		fmt.Printf("This will remove this distro (%s) from the filesystem.\n", name)
		fmt.Printf("Are you sure you would like to proceed? (This cannot be undone)\n")
		fmt.Printf("Type \"y\" to continue:")
		fmt.Scan(&in)

		if in != "y" {
			fmt.Fprintf(os.Stderr, "Accepting is required to proceed.")
			utils.ErrorExit(os.ErrInvalid, false, true, false)
		}
	} else {
		showProgress = false
	}

	Clean(name, showProgress)
}
