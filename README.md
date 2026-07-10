# Autofix CI

[![Step changelog](https://shields.io/github/v/release/bitrise-steplib/bitrise-step-autofix-ci?include_prereleases&label=changelog&color=blueviolet)](https://github.com/bitrise-steplib/bitrise-step-autofix-ci/releases)

Commits and pushes any local file changes made by previous steps, then fails the build to trigger a fresh run.

<details>
<summary>Description</summary>

Autofix CI detects uncommitted file changes left by previous steps (e.g. code formatters, linters with auto-fix, code generators) and automatically commits and pushes them back to the source branch. Inspired by [autofix.ci](https://autofix.ci/) for GitHub Actions.

After pushing, the step **intentionally fails the build** so that the unfixed commit doesn't pass any quality gates downstream. The newly pushed commit triggers a fresh CI run, which succeeds because the fixes are now in place.

#### Typical workflow placement

Place this step **after** any steps that may modify files (formatters, linters, generators) and **before** any steps that enforce quality gates (tests, lint checks). When the step finds changes to push, the build fails at this point and the quality gates run on the corrected commit in the next build.

#### How it works

1. Detects changed files via `git status` (including untracked files by default)
2. Aborts if any changed file is a Bitrise CI config (`bitrise.yml`, `bitrise.yaml`, `.bitrise/**`) to prevent privilege escalation
3. Commits all changes using a bot identity (`Bitrise Autofix`)
4. Pushes to the source branch (see **Authentication** below)
5. Exits with failure so CI gates don't pass on the unfixed commit

#### What triggers a skip

The step exits successfully without doing anything in these cases:

- **Not a PR build**: the `BITRISE_PULL_REQUEST` environment variable is not set. The step is designed for PR workflows; push builds are left untouched.
- **Fork PR**: the PR source repository differs from the target repository. The step cannot push to a forked repository using the provided credentials.
- **No changes detected**: there are no uncommitted modifications to commit.

#### Authentication

The step supports both HTTPS and SSH remotes — whichever the preceding Git Clone step configured.

- **HTTPS**: The step uses the `git_token` (and optionally `git_username`) inputs to authenticate the push via git's credential helper protocol. Credentials are passed to the `git push` subprocess through environment variables and never written to disk or embedded in the remote URL.
- **SSH**: If the repository was cloned over SSH and an SSH key is already loaded in the agent, the push uses SSH automatically. The `git_token` and `git_username` inputs are ignored for SSH remotes. Note that SSH deploy keys are typically read-only — if the push fails with a permission error, switch to HTTPS authentication using the `git_remote_url` and `git_token` inputs.

#### Security

If any changed file is `bitrise.yml`, `bitrise.yaml`, or anything under `.bitrise/`, the step aborts with an error instead of committing. This prevents a malicious PR from using the autofix mechanism to sneak CI configuration changes through an auto-commit.

#### Outputs

- `AUTOFIX_NEEDED`: `true` if uncommitted changes were detected
- `AUTOFIX_PUSHED`: `true` if the autofix commit was pushed successfully
- `AUTOFIX_FILE_COUNT`: number of files included in the autofix commit

</details>

## 🧩 Get started

Add this step directly to your workflow in the [Bitrise Workflow Editor](https://docs.bitrise.io/en/bitrise-ci/workflows-and-pipelines/steps/adding-steps-to-a-workflow.html).

You can also run this step directly with [Bitrise CLI](https://github.com/bitrise-io/bitrise).

## ⚙️ Configuration

<details>
<summary>Inputs</summary>

| Key | Description | Flags | Default |
| --- | --- | --- | --- |
| `commit_subject` | Subject line of the autofix commit message. Additional context (changed files, step URL) is automatically appended as the commit body. | required | `Bitrise CI Autofix` |
| `include_untracked` | When enabled (default), new files that are not yet tracked by git are included in the autofix commit. This is the right behavior for most cases, such as code generators and formatters that create new files.  When disabled, only modifications to already-tracked files are committed. Useful when the previous steps may create temporary or build output files that should not be committed.  | required | `true` |
| `git_username` | Username for HTTPS git push authentication. Defaults to the Bitrise-provided HTTP username. Not used for SSH remotes. | sensitive | `$GIT_HTTP_USERNAME` |
| `git_token` | Token or password for HTTPS git push authentication. Defaults to the Bitrise-provided HTTP password. Not used for SSH remotes. | sensitive | `$GIT_HTTP_PASSWORD` |
| `git_remote_url` | When set, the step updates the `origin` remote to this URL before fetching and pushing. This allows overriding an SSH remote (set by the Git Clone step) with an HTTPS URL so that token-based authentication works.  Defaults to `$BITRISEIO_BASE_REPOSITORY_URL`, the HTTPS URL of the base repository provided by Bitrise. Leave empty to use the remote URL already configured in the local git repository.  **SSH vs HTTPS:** SSH deploy keys are typically read-only, so a push over SSH will fail with a permission error. Switching to HTTPS with a `git_token` is usually the right fix: keep this input at its default HTTPS value and set `git_token` to a token with write access.  |  | `$BITRISEIO_BASE_REPOSITORY_URL` |
| `dry_run` | When enabled, the step detects changes, runs security checks, and creates a local commit — but skips the `git push`.  The step exits with success (instead of the usual intentional failure) so the build is not affected.  Note: `AUTOFIX_PUSHED` will always be `false` in dry run mode.  | required | `false` |
| `verbose` | Enable logging additional information for troubleshooting. | required | `false` |
</details>

<details>
<summary>Outputs</summary>

| Environment Variable | Description |
| --- | --- |
| `AUTOFIX_NEEDED` | Whether uncommitted changes were detected. `true` or `false`. |
| `AUTOFIX_PUSHED` | Whether the autofix commit was successfully pushed. `true` or `false`. |
| `AUTOFIX_FILE_COUNT` | Number of files included in the autofix commit. |
</details>

## 🙋 Contributing

We welcome [pull requests](https://github.com/bitrise-steplib/bitrise-step-autofix-ci/pulls) and [issues](https://github.com/bitrise-steplib/bitrise-step-autofix-ci/issues) against this repository.

For pull requests, work on your changes in a forked repository and use the Bitrise CLI to [run step tests locally](https://docs.bitrise.io/en/bitrise-ci/bitrise-cli/running-your-first-local-build-with-the-cli.html).

Learn more about developing steps:

- [Create your own step](https://docs.bitrise.io/en/bitrise-ci/workflows-and-pipelines/developing-your-own-bitrise-step/developing-a-new-step.html)
