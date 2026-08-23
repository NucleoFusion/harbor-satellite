package apply

import (
	"bytes"
	"errors"
	"fmt"

	"github.com/container-registry/harbor-satellite/internal/groundcontrol/cli/common"
	"github.com/container-registry/harbor-satellite/internal/groundcontrol/cli/models"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

func NewApplyCommand(runtime *common.Runtime) *cobra.Command {
	var (
		argFiles  []string
		dryRun    bool
		recursive bool
	)
	applyCmd := &cobra.Command{
		Use:   "apply",
		Short: "Apply configuration",
		Long:  "Apply a state configuration to manage harbor-satellite objects",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			files, err := getAllFiles(argFiles, recursive)
			if err != nil {
				return err
			}

			var errs []error
			for _, v := range files {
				content, err := readFile(v)
				if err != nil {
					errs = append(errs, err)
					continue
				}

				resources, err := ParseYAML(content)
				if err != nil {
					errs = append(errs, err)
					continue
				}

				for _, resource := range resources {
					ApplyResource(runtime, resource)
				}
			}

			return errors.Join(errs...)
		},
	}

	applyCmd.Flags().StringArrayVarP(&argFiles, "filename", "f", nil, "configuration filepaths")
	applyCmd.MarkFlagRequired("filename")

	applyCmd.Flags().BoolVar(&dryRun, "dry-run", false, "dry run configuration files input")
	applyCmd.Flags().BoolVar(&recursive, "recursive", false, "recursively execute all files in directories passed as args")
	return applyCmd
}

func ApplyResource(rt *common.Runtime, resource models.YAMLResource) models.ApplyResult {
	switch resource.Kind {
	// case "Satellite":
	case "Group":
		return applyGroupResource(rt, resource)
		// case "Config":
	}

	return errorToResult(resource.Kind, resource.Metadata.Name, fmt.Errorf("invalid resource"))
}

func ParseYAML(data []byte) ([]models.YAMLResource, error) {
	dec := yaml.NewDecoder(bytes.NewReader(data))

	var resources []models.YAMLResource
	for {
		var env models.Envelope
		err := dec.Decode(&env)
		if err != nil {
			if err.Error() == "EOF" {
				break
			}
			return nil, fmt.Errorf("decoding document: %w", err)
		}

		spec, err := decodeSpec(env.Kind, env.Spec)
		if err != nil {
			return nil, fmt.Errorf("resource %s/%s: %w", env.Kind, env.Metadata.Name, err)
		}

		resources = append(resources, models.YAMLResource{
			APIVersion: env.APIVersion,
			Kind:       env.Kind,
			Metadata:   env.Metadata,
			Spec:       spec,
		})
	}
	return resources, nil
}

func decodeSpec(kind string, raw yaml.Node) (any, error) {
	switch kind {
	case "Group":
		var s models.GroupSpec
		if err := raw.Decode(&s); err != nil {
			return nil, err
		}
		return &s, nil
	case "Config":
		var s models.ConfigSpec
		if err := raw.Decode(&s); err != nil {
			return nil, err
		}
		return &s, nil
	case "Satellite":
		var s models.SatelliteSpec
		if err := raw.Decode(&s); err != nil {
			return nil, err
		}
		return &s, nil
	default:
		return nil, fmt.Errorf("unknown kind %q", kind)
	}
}

func errorToResult(kind, name string, err error) models.ApplyResult {
	return models.ApplyResult{
		Kind:  kind,
		Name:  name,
		Error: err,
	}
}
