package excubitor

import (
	"errors"
	"fmt"
	"os/user"
	"runtime"
)

var WrongOSError error = errors.New("unsupported OS")
var SoftfailError error = errors.New("soft fail")

func check() error {
	if runtime.GOOS != "linux" {
		return fmt.Errorf("%w: %s", WrongOSError, runtime.GOOS)
	}

	currUser, err := user.Current()
	if err != nil {
		return err
	}

	if currUser.Uid != "0" {
		return fmt.Errorf("%w: %s", SoftfailError, fmt.Sprintf("the application is not running as root. You may experience inconviniences such as only being able to log in as your current account. To circumvent this, you may use an LDAP server for authentication through PAM."))
	}

	return nil
}
