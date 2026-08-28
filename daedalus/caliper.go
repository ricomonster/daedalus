package daedalus

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/ricomonster/daedalus/git"
)

const (
	VersionTypeVersionFile VersionType = "version-file"
	VersionTypeNodePackage VersionType = "node-package"
	VersionTypeGit         VersionType = "git"
)

var semverPattern = regexp.MustCompile(`^(v?)([0-9]+)\.([0-9]+)\.([0-9]+)$`)

type (
	VersionType string

	VersionSource struct {
		Type VersionType
		Path string
	}
)

func NewCaliperService() *CaliperService {
	return &CaliperService{}
}

type CaliperService struct{}

func (c *CaliperService) DetectVersionSource(path string) (*VersionSource, error) {
	root := path
	if root == "" {
		root = "."
	}

	info, err := os.Stat(root)
	if err != nil {
		return nil, err
	}

	if !info.IsDir() {
		switch filepath.Base(root) {
		case "VERSION":
			return &VersionSource{VersionTypeVersionFile, root}, nil
		case "package.json":
			return &VersionSource{VersionTypeNodePackage, root}, nil
		default:
			return nil, fmt.Errorf("path must be VERSION, package.json, or a directory")
		}
	}

	var matches []VersionSource
	for _, name := range []string{"VERSION", "package.json"} {
		file := filepath.Join(root, name)
		if _, err := os.Stat(file); err == nil {
			switch name {
			case "VERSION":
				matches = append(matches, VersionSource{
					Type: VersionTypeVersionFile,
					Path: file,
				})

			case "package.json":
				matches = append(matches, VersionSource{
					Type: VersionTypeNodePackage,
					Path: file,
				})
			}
		}
	}

	if len(matches) > 1 {
		return nil, fmt.Errorf("multiple version files found; specify --path explicitly")
	}

	if len(matches) == 1 {
		return &matches[0], nil
	}

	// Default git
	return &VersionSource{VersionTypeGit, root}, nil
}

func (vs VersionSource) Read() (string, error) {
	switch vs.Type {
	case VersionTypeVersionFile:
		data, err := os.ReadFile(vs.Path)
		return strings.TrimSpace(string(data)), err

	case VersionTypeNodePackage:
		data, err := os.ReadFile(vs.Path)
		if err != nil {
			return "", err
		}

		var packageJSON map[string]json.RawMessage
		if err := json.Unmarshal(data, &packageJSON); err != nil {
			return "", err
		}

		var version string
		if err := json.Unmarshal(packageJSON["version"], &version); err != nil {
			return "", fmt.Errorf("package.json has no valid version")
		}

		return version, nil

	case VersionTypeGit:
		// Run git
		tags, err := git.New(vs.Path).Tags()
		if err != nil {
			return "", err
		}

		for _, tag := range tags {
			matches := semverPattern.FindStringSubmatch(tag)
			if matches != nil {
				return matches[2] + "." + matches[3] + "." + matches[4], nil
			}
		}

		return "", fmt.Errorf("no semantic version git tag found")
	}

	return "", fmt.Errorf("unknown version source: %s", vs.Type)
}

func (vs VersionSource) Write(version string) error {
	switch vs.Type {
	case VersionTypeVersionFile:
		return os.WriteFile(vs.Path, []byte(version+"\n"), 0o644)

	case VersionTypeNodePackage:
		data, err := os.ReadFile(vs.Path)
		if err != nil {
			return err
		}

		var packageJSON map[string]json.RawMessage
		if err := json.Unmarshal(data, &packageJSON); err != nil {
			return err
		}

		value, _ := json.Marshal(version)
		packageJSON["version"] = value

		output, err := json.MarshalIndent(packageJSON, "", " ")
		if err != nil {
			return err
		}

		output = append(output, '\n')

		return os.WriteFile(vs.Path, output, 0o644)

	case VersionTypeGit:
		if version == "" {
			return fmt.Errorf("cannot write empty version")
		}

		tagName := "v" + strings.TrimPrefix(version, "v")

		// Create and push the tag
		return git.New(vs.Path).PushTag(tagName)
	}

	return fmt.Errorf("source %s is not writable", vs.Type)
}
