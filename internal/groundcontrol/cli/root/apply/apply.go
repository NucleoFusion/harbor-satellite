package apply

import (
	"fmt"

	"github.com/container-registry/harbor-satellite/internal/groundcontrol/cli/common"
	"github.com/spf13/cobra"
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

			for _, v := range files {
				fmt.Println(v)
			}
			return nil
		},
	}

	applyCmd.Flags().StringArrayVarP(&argFiles, "filename", "f", nil, "configuration filepaths")
	applyCmd.MarkFlagRequired("filename")

	applyCmd.Flags().BoolVar(&dryRun, "dry-run", false, "dry run configuration files input")
	applyCmd.Flags().BoolVar(&recursive, "recursive", false, "recursively execute all files in directories passed as args")
	return applyCmd
}

func getAllFiles(argFiles []string, recursive bool) ([]string, error) {
	files := argFiles

	for _, v := range files {
		isDir, err := isDirectory(v)
		if err != nil {
			return nil, err
		}

		if isDir {
			dirFiles, err := WalkDirectory(v, recursive)
			if err != nil {
				return nil, err
			}

			files = append(files, dirFiles...)
		}
	}

	return files, nil
}
