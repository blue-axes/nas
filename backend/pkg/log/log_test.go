package log

import "testing"

func TestLog(t *testing.T) {
	Info("info\n")
	Infof("info f\n")
	Infoln("info ln")
}
