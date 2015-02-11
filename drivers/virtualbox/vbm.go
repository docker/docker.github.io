package virtualbox

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"

	log "github.com/Sirupsen/logrus"
)

var (
	reVMNameUUID      = regexp.MustCompile(`"(.+)" {([0-9a-f-]+)}`)
	reVMInfoLine      = regexp.MustCompile(`(?:"(.+)"|(.+))=(?:"(.*)"|(.*))`)
	reColonLine       = regexp.MustCompile(`(.+):\s+(.*)`)
	reMachineNotFound = regexp.MustCompile(`Could not find a registered machine named '(.+)'`)
)

var (
	ErrMachineExist    = errors.New("machine already exists")
	ErrMachineNotExist = errors.New("machine does not exist")
	ErrVBMNotFound     = errors.New("VBoxManage not found")
	vboxManageCmd      = setVBoxManageCmd()
)

// detect the VBoxManage cmd's path if needed
func setVBoxManageCmd() string {
	cmd := "VBoxManage"
	if path, err := exec.LookPath(cmd); err == nil {
		return path
	}
	if runtime.GOOS == "windows" {
		if p := os.Getenv("VBOX_INSTALL_PATH"); p != "" {
			if path, err := exec.LookPath(filepath.Join(p, cmd)); err == nil {
				return path
			}
		}
		if p := os.Getenv("VBOX_MSI_INSTALL_PATH"); p != "" {
			if path, err := exec.LookPath(filepath.Join(p, cmd)); err == nil {
				return path
			}
		}
		// look at HKEY_LOCAL_MACHINE\SOFTWARE\Oracle\VirtualBox\InstallDir
		p := "C:\\Program Files\\Oracle\\VirtualBox"
		if path, err := exec.LookPath(filepath.Join(p, cmd)); err == nil {
			return path
		}
	}
	return cmd
}

func vbm(args ...string) error {
	cmd := exec.Command(vboxManageCmd, args...)
	if os.Getenv("DEBUG") != "" {
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	}
	log.Debugf("executing: %v %v", vboxManageCmd, strings.Join(args, " "))
	if err := cmd.Run(); err != nil {
		if ee, ok := err.(*exec.Error); ok && ee == exec.ErrNotFound {
			return ErrVBMNotFound
		}
		return fmt.Errorf("%v %v failed: %v", vboxManageCmd, strings.Join(args, " "), err)
	}
	return nil
}

func vbmOut(args ...string) (string, error) {
	cmd := exec.Command(vboxManageCmd, args...)
	if os.Getenv("DEBUG") != "" {
		cmd.Stderr = os.Stderr
	}
	log.Debugf("executing: %v %v", vboxManageCmd, strings.Join(args, " "))

	b, err := cmd.Output()
	if err != nil {
		if ee, ok := err.(*exec.Error); ok && ee == exec.ErrNotFound {
			err = ErrVBMNotFound
		}
	}
	return string(b), err
}

func vbmOutErr(args ...string) (string, string, error) {
	cmd := exec.Command(vboxManageCmd, args...)
	log.Debugf("executing: %v %v", vboxManageCmd, strings.Join(args, " "))
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		if ee, ok := err.(*exec.Error); ok && ee == exec.ErrNotFound {
			err = ErrVBMNotFound
		}
	}
	return stdout.String(), stderr.String(), err
}
