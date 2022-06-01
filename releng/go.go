package releng

import (
	"os"
	"os/exec"
)

// GoSpec configures a Go binary build.
type GoSpec struct {
	OS, Arch, Arm, LDFlags string
}

// BuildGo builds a Go binary.
func BuildGo(pkg, dest string, spec GoSpec) {
	log.Info("Building Go binary: ", dest)
	log.Debug("Package: ", pkg)
	log.Debug("Destination: ", dest)
	log.Debug("OS: ", spec.OS)
	log.Debug("Arch: ", spec.Arch)
	log.Debug("Arm: ", spec.Arm)
	log.Debug("LDFlags: ", spec.LDFlags)

	cmd := exec.Command("go", "build", "-o", dest, "-trimpath", "-ldflags="+spec.LDFlags, pkg)
	env := os.Environ()
	env = append(env, "GOOS="+spec.OS, "GOARCH="+spec.Arch, "GOARM="+spec.Arm)
	cmd.Env = env
	out, err := cmd.CombinedOutput()
	if len(out) > 0 {
		if err != nil {
			log.Warning("Output: ", string(out))
		} else {
			log.Debug("Output: ", string(out))
		}
	}
	Must(err)
}
