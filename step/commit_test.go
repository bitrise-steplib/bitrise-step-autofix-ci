package step

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_buildCommitMessage(t *testing.T) {
	msg := buildCommitMessage("Bitrise CI Autofix", []string{"main.go", "step/step.go"})

	assert.True(t, strings.HasPrefix(msg, "Bitrise CI Autofix\n"), "message should start with the subject line")
	assert.Contains(t, msg, "Created by the Autofix CI step")
	assert.Contains(t, msg, stepRepoURL)
	assert.Contains(t, msg, "- main.go\n")
	assert.Contains(t, msg, "- step/step.go\n")

	// Files come first, the explanation and URL follow as a footer.
	headerPos := strings.Index(msg, "Modified files:")
	filesPos := strings.Index(msg, "- main.go")
	urlPos := strings.Index(msg, stepRepoURL)
	assert.Less(t, headerPos, filesPos)
	assert.Less(t, filesPos, urlPos, "file list should appear before the step URL")
}
