package generate

import (
	"strings"

	"github.com/pkg/errors"
)

// EnvVariable "PODMAN_SYSTEMD_UNIT" is set in all generated systemd units and
// is set to the unit's (unique) name.
const EnvVariable = "PODMAN_SYSTEMD_UNIT"

// restartPolicies includes all valid restart policies to be used in a unit
// file.
var restartPolicies = []string{"no", "on-success", "on-failure", "on-abnormal", "on-watchdog", "on-abort", "always"}

// validateRestartPolicy checks that the user-provided policy is valid.
func validateRestartPolicy(restart string) error {
	for _, i := range restartPolicies {
		if i == restart {
			return nil
		}
	}
	return errors.Errorf("%s is not a valid restart policy", restart)
}

const headerTemplate = `# {{.ServiceName}}.service
# autogenerated by Podman {{.PodmanVersion}}
{{- if .TimeStamp}}
# {{.TimeStamp}}
{{- end}}

[Unit]
Description=Podman {{.ServiceName}}.service
Documentation=man:podman-generate-systemd(1)
Wants=network.target
After=network-online.target
`

// filterPodFlags removes --pod and --pod-id-file from the specified command.
func filterPodFlags(command []string) []string {
	processed := []string{}
	for i := 0; i < len(command); i++ {
		s := command[i]
		if s == "--pod" || s == "--pod-id-file" {
			i++
			continue
		}
		if strings.HasPrefix(s, "--pod=") || strings.HasPrefix(s, "--pod-id-file=") {
			continue
		}
		processed = append(processed, s)
	}
	return processed
}
