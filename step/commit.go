package step

import "strings"

const stepRepoURL = "https://github.com/bitrise-steplib/bitrise-step-autofix-ci"

func buildCommitMessage(subject string, changedFiles []string) string {
	var sb strings.Builder
	sb.WriteString(subject)
	sb.WriteString("\n\nModified files:\n")
	for _, f := range changedFiles {
		sb.WriteString("- ")
		sb.WriteString(f)
		sb.WriteString("\n")
	}
	sb.WriteString("\nCreated by the Autofix CI step to capture changes from earlier\n")
	sb.WriteString("workflow steps (e.g. a formatter, linter or code generator).\n")
	sb.WriteString(stepRepoURL)
	sb.WriteString("\n")
	return sb.String()
}
