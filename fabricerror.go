package fabric

import (
	"fmt"

	"github.com/pkg/errors"
)

func (c FabricErrorCode) Error() string {
	if c == 0 {
		return ""
	}

	return fmt.Sprintf("error [%v] [0x%x]", c.String(), uint64(c))
}

func errno(hr uintptr, syserr error) error {
	if hr == 0 {
		return nil
	}

	return errors.Wrap(FabricErrorCode(hr), fabricGetLastError())
}
