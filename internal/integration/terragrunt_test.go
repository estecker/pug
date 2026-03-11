package app

import (
	"os/exec"
	"strings"
	"testing"

	tea "charm.land/bubbletea/v2"
)

func TestTerragrunt_SingleInit(t *testing.T) {
	t.Parallel()
	skipIfTerragruntNotFound(t)

	_ = setupAndInitTerragruntModule(t)
}

func TestTerragrunt_SinglePlan(t *testing.T) {
	t.Parallel()
	skipIfTerragruntNotFound(t)

	tm := setupAndInitTerragruntModule(t)

	// Create plan on first module
	tm.Type("p")
	waitFor(t, tm, func(s string) bool {
		return strings.Contains(s, "plan 󰠱 modules/a  default") &&
			strings.Contains(s, "exited +10~0-0") &&
			strings.Contains(s, "terragrunt")
	})
}

func TestTerragrunt_SingleApply(t *testing.T) {
	t.Parallel()
	skipIfTerragruntNotFound(t)

	tm := setupAndInitTerragruntModule(t)

	// Create apply for module.
	tm.Type("a")

	// Give approval
	waitFor(t, tm, func(s string) bool {
		return strings.Contains(s, "Auto-apply 1 workspaces? (y/N):")
	})
	tm.Type("y")

	// Send to apply task page
	waitFor(t, tm, func(s string) bool {
		return strings.Contains(s, "apply 󰠱 modules/a  default") &&
			strings.Contains(s, "exited +10~0-0") &&
			strings.Contains(s, "terragrunt")
	})
}

// TestTerragrunt_Dependencies tests that terragrunt dependencies are
// respected.
func TestTerragrunt_Dependencies(t *testing.T) {
	t.Parallel()
	skipIfTerragruntNotFound(t)

	tm := setupAndInitTerragruntModulesWithDependencies(t)

	// Select all modules and create apply on each
	tm.Send(tea.KeyPressMsg{Code: 'a', Mod: tea.ModCtrl})
	tm.Type("a")

	// Give approval
	waitFor(t, tm, func(s string) bool {
		return strings.Contains(s, "Auto-apply 6 workspaces? (y/N):")
	})
	tm.Type("y")

	// Expect 6 applies. The "." module fails because it doesn't have any config
	// files.
	waitFor(t, tm, func(s string) bool {
		return strings.Contains(s, "apply 1/5/6") &&
			matchPattern(t, `modules/vpc.*default.*\+0~0-0`, s) &&
			matchPattern(t, `modules/redis.*default.*\+0~0-0`, s) &&
			matchPattern(t, `modules/mysql.*default.*\+0~0-0`, s) &&
			matchPattern(t, `modules/backend-app.*default.*\+0~0-0`, s) &&
			matchPattern(t, `modules/frontend-app.*default.*\+0~0-0`, s) &&
			matchPattern(t, `\..*default.*errored`, s) &&
			// Expect several modules to now have some resources
			strings.Contains(s, "└  default ✓ 3") &&
			strings.Contains(s, "└  default ✓ 2") &&
			strings.Contains(s, "└  default ✓ 1") &&
			strings.Contains(s, "└  default ✓ 1") &&
			strings.Contains(s, "└  default ✓ 0") &&
			strings.Contains(s, "└  default ✓ 0")
	})

	// Go back to explorer.
	tm.Type("0")

	// Destroy resources in all modules (they should still all be selected).
	tm.Type("D")

	// Give approval
	waitFor(t, tm, func(s string) bool {
		return strings.Contains(s, "Destroy resources of 6 workspaces? (y/N):")
	})
	tm.Type("y")

	// Expect 6 apply tasks.
	waitFor(t, tm, func(s string) bool {
		return strings.Contains(s, "apply (destroy) 6/6") &&
			matchPattern(t, `modules/vpc.*default.*apply \(destroy\).*exited.*\+0~0-0`, s) &&
			matchPattern(t, `modules/redis.*default.*apply \(destroy\).*exited.*\+0~0-0`, s) &&
			matchPattern(t, `modules/mysql.*default.*apply \(destroy\).*exited.*\+0~0-0`, s) &&
			matchPattern(t, `modules/backend-app.*default.*apply \(destroy\).*exited.*\+0~0-0`, s) &&
			matchPattern(t, `modules/frontend-app.*default.*apply \(destroy\).*exited.*\+0~0-0`, s) &&
			matchPattern(t, `\..*default.*apply \(destroy\).*exited.*\+0~0-0`, s) &&
			// Expect modules to now have some 0 resources
			strings.Count(s, "└  default ✓ 0") >= 6
	})
}

func skipIfTerragruntNotFound(t *testing.T) {
	if _, err := exec.LookPath("terragrunt"); err != nil {
		t.Skip("skipping test: terragrunt not found")
	}
}

func setupAndInitTerragruntModule(t *testing.T) *testModel {
	tm := setup(t, "./testdata/single_terragrunt_module", withTerragrunt())

	// Expect single module to be listed
	waitFor(t, tm, func(s string) bool {
		return strings.Contains(s, "└ 󰠱 a")
	})

	// Initialize module
	tm.Type("i")
	// Expect init to succeed, and to populate pug with one workspace with 0
	// resources
	waitFor(t, tm, func(s string) bool {
		return strings.Contains(s, "Terraform has been successfully initialized!") &&
			strings.Contains(s, "init 󰠱 modules/a") &&
			strings.Contains(s, "exited") &&
			strings.Contains(s, "└  default ✓ 0")
	})

	// Show task info sidebar so tests can check that terragrunt is indeed being
	// executed.
	tm.Type("I")
	waitFor(t, tm, func(s string) bool {
		return strings.Contains(s, "terragrunt")
	})

	// Go back to explorer
	tm.Type("0")

	return tm
}

func setupAndInitTerragruntModulesWithDependencies(t *testing.T) *testModel {
	tm := setup(t, "./testdata/terragrunt_modules_with_dependencies", withTerragrunt())

	waitFor(t, tm, func(s string) bool {
		return strings.Contains(s, "├ 󰠱 backend-app") &&
			strings.Contains(s, "├ 󰠱 frontend-app") &&
			strings.Contains(s, "├ 󰠱 mysql") &&
			strings.Contains(s, "├ 󰠱 redis") &&
			strings.Contains(s, "└ 󰠱 vpc") &&
			strings.Contains(s, "└ 󰠱 .")
	})

	// Select all modules and init
	tm.Send(tea.KeyPressMsg{Code: 'a', Mod: tea.ModCtrl})
	tm.Type("i")
	waitFor(t, tm, func(s string) bool {
		return strings.Contains(s, "init 6/6") &&
			matchPattern(t, `modules/vpc.*init.*exited`, s) &&
			matchPattern(t, `modules/redis.*init.*exited`, s) &&
			matchPattern(t, `modules/mysql.*init.*exited`, s) &&
			matchPattern(t, `modules/frontend-app.*init.*exited`, s) &&
			matchPattern(t, `modules/backend-app.*init.*exited`, s) &&
			matchPattern(t, `\..*init.*exited`, s) &&
			// Expect modules to be listed along with their default workspace.
			strings.Count(s, "└  default ✓ 0") >= 6
	})

	// Go back to explorer and clear selection.
	tm.Type("0")
	tm.Send(tea.KeyPressMsg{Code: '\\', Mod: tea.ModCtrl})

	return tm
}
