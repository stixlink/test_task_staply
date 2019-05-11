package utility

import (
	"os"
	"strings"
)

func RemoveErrorFiles(paths ...string) []error {
	var (
		errs []error
		err  error
	)
	for _, f := range paths {
		err = os.Remove(f)
		if err != nil {
			errs = append(errs, err)
		}
	}
	return errs
}

func JoinError(errors []error, separate string) string {
	var messages []string
	for _, e := range errors {
		messages = append(messages, e.Error())
	}
	return strings.Join(messages, separate)
}
