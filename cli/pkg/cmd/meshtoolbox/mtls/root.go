package mtls

import (
	"fmt"

	"github.com/solo-io/solo-kit/pkg/api/v1/clients"
	// "github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
	"github.com/solo-io/supergloo/cli/pkg/cmd/options"
	"github.com/solo-io/supergloo/cli/pkg/common"
	"github.com/solo-io/supergloo/cli/pkg/nsutil"
	superglooV1 "github.com/solo-io/supergloo/pkg/api/v1"
	"github.com/spf13/cobra"
)

var validRootArgs = []string{"enable", "disable", "toggle"} // for bash completion

func Root(opts *options.Options) *cobra.Command {
	cmd := &cobra.Command{
		Use:       "mtls",
		Short:     `set mTLS status`,
		Long:      `set mTLS status`,
		ValidArgs: validRootArgs,
		Args:      rootArgValidation,
		RunE: func(c *cobra.Command, args []string) error {
			// this function does nothing but it triggers validation
			return nil
		},
	}
	cmd.AddCommand(
		Enable(opts),
		Disable(opts),
		Toggle(opts),
	)
	return cmd
}

func rootArgValidation(c *cobra.Command, args []string) error {
	expectedArgCount := 1
	if len(args) != expectedArgCount {
		return fmt.Errorf("Too many args (%v given, %v expected)", len(args), expectedArgCount)
	}
	subCommandName := args[0]
	if !common.Contains(validRootArgs, subCommandName) {
		return fmt.Errorf("%v is not a valid argument", subCommandName)
	}
	return nil
}

func Enable(opts *options.Options) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "enable",
		Short: `enable mTLS`,
		Long:  `enable mTLS`,
		RunE: func(c *cobra.Command, args []string) error {
			if err := enableMtls(opts); err != nil {
				return err
			}
			return nil
		},
	}
	return cmd
}

func Disable(opts *options.Options) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "disable",
		Short: `disable mTLS`,
		Long:  `disable mTLS`,
		RunE: func(c *cobra.Command, args []string) error {
			if err := disableMtls(opts); err != nil {
				return err
			}
			return nil
		},
	}
	return cmd
}

func Toggle(opts *options.Options) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "toggle",
		Short: `toggle mTLS`,
		Long:  `toggle mTLS`,
		RunE: func(c *cobra.Command, args []string) error {
			if err := toggleMtls(opts); err != nil {
				return err
			}
			return nil
		},
	}
	return cmd
}

func enableMtls(opts *options.Options) error {

	if _, err := updateMtls("enable", opts); err != nil {
		return err
	}
	fmt.Printf("Enabled mTLS on mesh %v", opts.MeshTool.Mesh.Name)

	return nil
}

func disableMtls(opts *options.Options) error {
	if _, err := updateMtls("disable", opts); err != nil {
		return err
	}
	fmt.Printf("Disabled mTLS on mesh %v", opts.MeshTool.Mesh.Name)
	return nil
}

func toggleMtls(opts *options.Options) error {
	mesh, err := updateMtls("toggle", opts)
	if err != nil {
		return err
	}
	status := "disabled"
	if mesh.Encryption.TlsEnabled {
		status = "enabled"
	}
	fmt.Printf("Toggled (%v) mTLS on mesh %v", status, opts.MeshTool.Mesh.Name)
	return nil
}

// Ensure that all the needed user-specified values have been provided
func ensureFlags(operation string, opts *options.Options) error {

	// all operations require a target mesh spec
	meshRef := &(opts.MeshTool).Mesh
	if err := nsutil.EnsureMesh(meshRef, opts); err != nil {
		return err
	}

	return nil
}

func updateMtls(operation string, opts *options.Options) (*superglooV1.Mesh, error) {
	// 1. validate/aquire arguments
	if err := ensureFlags(operation, opts); err != nil {
		return nil, err
	}

	// 2. read the existing mesh
	meshClient, err := common.GetMeshClient()
	if err != nil {
		return nil, err
	}
	meshRef := &(opts.MeshTool).Mesh
	mesh, err := (*meshClient).Read(meshRef.Namespace, meshRef.Name, clients.ReadOpts{})
	if err != nil {
		return nil, err
	}

	// 3. mutate the mesh structure
	switch operation {
	case "enable":
		if mesh.Encryption == nil {
			mesh.Encryption = &superglooV1.Encryption{
				TlsEnabled: true,
			}
		} else {
			mesh.Encryption.TlsEnabled = true

		}
	case "disable":
		if mesh.Encryption == nil {
			mesh.Encryption = &superglooV1.Encryption{
				TlsEnabled: false,
			}
		} else {
			mesh.Encryption.TlsEnabled = false

		}
	case "toggle":
		// if encryption has not been specified, "toggle" will enable it
		if mesh.Encryption == nil {
			mesh.Encryption = &superglooV1.Encryption{
				TlsEnabled: true,
			}
		} else {
			mesh.Encryption.TlsEnabled = !mesh.Encryption.TlsEnabled

		}
	default:
		panic(fmt.Errorf("Operation %v not recognized", operation))
	}

	// 4. write the changes
	writtenMesh, err := (*meshClient).Write(mesh, clients.WriteOpts{OverwriteExisting: true})
	if err != nil {
		return nil, err
	}
	return writtenMesh, nil
}
