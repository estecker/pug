package workspace

import (
	"fmt"
	"log/slog"
	"net/url"
	"os"
	"path/filepath"

	"github.com/leg100/pug/internal"
	"github.com/leg100/pug/internal/module"
	"github.com/leg100/pug/internal/resource"
)

type Workspace struct {
	ID         resource.MonotonicID
	Name       string
	ModuleID   resource.MonotonicID
	ModulePath string
	Cost       *float64
}

func New(mod *module.Module, name string) (*Workspace, error) {
	if name != url.PathEscape(name) {
		return nil, fmt.Errorf("invalid workspace name: %s", name)
	}
	return &Workspace{
		ID:         resource.NewMonotonicID(resource.Workspace),
		Name:       name,
		ModuleID:   mod.ID,
		ModulePath: mod.Path,
	}, nil
}

func (ws *Workspace) String() string {
	return ws.Name
}

func (ws *Workspace) TerraformEnv() string {
	return TerraformEnv(ws.Name)
}

func (ws *Workspace) LogValue() slog.Value {
	return slog.GroupValue(
		slog.String("name", ws.Name),
	)
}

// VarFiles returns the path to the workspace's terraform variables file
// and whether it exists or not. If tfvars is provided, it looks in that
// directory instead of the module directory and returns the absolute path.
// Otherwise it returns just the filename (relative to the module directory).
func (ws *Workspace) VarFiles(workdir internal.Workdir, tfvars string) (string, bool) {
	fname := fmt.Sprintf("%s.tfvars", ws.Name)
	var path string
	if tfvars != "" {
		path = filepath.Join(tfvars, fname)
	} else {
		path = filepath.Join(workdir.String(), ws.ModulePath, fname)
	}
	_, err := os.Stat(path)
	// When using a custom tfvars directory, return the full path
	// Otherwise return just the filename (terraform will look in the module dir)
	if tfvars != "" {
		return path, err == nil
	}
	return fname, err == nil
}

func TerraformEnv(workspaceName string) string {
	return fmt.Sprintf("TF_WORKSPACE=%s", workspaceName)
}
