package utils

import (
	"bytes"
	"github.com/pkg/errors"
	"os/exec"
)

// todo:
// run everything in the test with kube injection on - add injection to all the Deploy commands
// or turn on auto inject for the ns
// then simply run cURL from the testrunner against reviews... try a curl w/ guaranteed response
//
func IstioInject(istioNamespace, input string) (string, error) {
	cmd := exec.Command("istioctl", "kube-inject", "-i", istioNamespace, "-f", "-")
	cmd.Stdin = bytes.NewBuffer([]byte(input))
	output := &bytes.Buffer{}
	cmd.Stdout = output
	cmd.Stderr = output
	err := cmd.Run()
	if err != nil {
		return "", errors.Wrapf(err, "kube inject failed: %v", output.String())
	}
	return output.String(), nil
}
