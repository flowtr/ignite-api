package pkg

import (
	"github.com/weaveworks/ignite/cmd/ignite/cmd/cmdutil"
	"github.com/weaveworks/ignite/pkg/apis/ignite"
	"github.com/weaveworks/ignite/pkg/apis/ignite/validation"
	meta "github.com/weaveworks/ignite/pkg/apis/meta/v1alpha1"
	"github.com/weaveworks/ignite/pkg/config"
	"github.com/weaveworks/ignite/pkg/dmlegacy"
	"github.com/weaveworks/ignite/pkg/metadata"
	"github.com/weaveworks/ignite/pkg/operations"
	"github.com/weaveworks/ignite/pkg/providers"
)

type VMData struct {
	Image    string    `json:"image"`
	Sandbox  string    `json:"sandbox"`
	Kernel   string    `json:"kernel"`
	CPUs     uint64    `json:"cpus"`
	Memory   meta.Size `json:"memory"`
	DiskSize meta.Size `json:"diskSize"`
	// TODO: Implement working omitempty without pointers for the following entries
	// Currently both will show in the JSON output as empty arrays. Making them
	// pointers requires plenty of nil checks (as their contents are accessed directly)
	// and is very risky for stability. APIMachinery potentially has a solution.
	Network ignite.VMNetworkSpec `json:"network,omitempty"`
	Storage ignite.VMStorageSpec `json:"storage,omitempty"`
	// This will be done at either "ignite start" or "ignite create" time
	// TODO: We might revisit this later
	CopyFiles []ignite.FileMapping `json:"copyFiles,omitempty"`
	// SSH specifies how the SSH setup should be done
	// nil here means "don't do anything special"
	// If SSH.Generate is set, Ignite will generate a new SSH key and copy it in to authorized_keys in the VM
	// Specifying a path in SSH.Generate means "use this public key"
	// If SSH.PublicKey is set, this struct will marshal as a string using that path
	// If SSH.Generate is set, this struct will marshal as a bool => true
	SSH *ignite.SSH `json:"ssh,omitempty"`
}

func CreateVM(data VMData) error {
	vm := providers.Client.VMs().New()

	// Resolve registry configuration used for pulling image if required.
	cmdutil.ResolveRegistryConfigDir()

	// Initialize the VM's Prefixer
	vm.Status.IDPrefix = providers.IDPrefix
	// Set the runtime and network-plugin on the VM, then override the global config.
	vm.Status.Runtime.Name = providers.RuntimeName
	vm.Status.Network.Plugin = providers.NetworkPluginName
	// Populate the runtime and network-plugin providers.
	if err := config.SetAndPopulateProviders(providers.RuntimeName, providers.NetworkPluginName); err != nil {
		return err
	}

	// Generate a random UID and Name
	if err := metadata.SetNameAndUID(vm, providers.Client); err != nil {
		return err
	}
	// Set VM labels.
	if err := metadata.SetLabels(vm, []string{}); err != nil {
		return err
	}

	ociRef, err := meta.NewOCIImageRef(data.Image)
	if err != nil {
		return err
	}
	vm.Spec.Image.OCI = ociRef

	ociRef, err = meta.NewOCIImageRef(data.Kernel)
	if err != nil {
		return err
	}
	vm.Spec.Kernel.OCI = ociRef

	ociRef, err = meta.NewOCIImageRef(data.Sandbox)
	if err != nil {
		return err
	}
	vm.Spec.Sandbox.OCI = ociRef

	// Get the image, or import it if it doesn't exist.
	image, err := operations.FindOrImportImage(providers.Client, vm.Spec.Image.OCI)
	if err != nil {
		return err
	}

	// Populate relevant data from the Image on the VM object.
	vm.SetImage(image)

	// Get the kernel, or import it if it doesn't exist.
	kernel, err := operations.FindOrImportKernel(providers.Client, vm.Spec.Kernel.OCI)
	if err != nil {
		return err
	}

	// Populate relevant data from the Kernel on the VM object.
	vm.SetKernel(kernel)

	if err := validation.ValidateVM(vm).ToAggregate(); err != nil {
		return err
	}

	// Create the vm
	if err := providers.Client.VMs().Set(vm); err != nil {
		return err
	}

	// Allocate and populate the overlay file
	if err := dmlegacy.AllocateAndPopulateOverlay(vm); err != nil {
		return err
	}

	if err := metadata.Success(vm); err != nil {
		return err
	}

	return nil
}
