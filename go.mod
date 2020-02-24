module github.com/jsirianni/zfsmon

go 1.13

require (
	// v10 release
	github.com/asaskevich/govalidator v0.0.0-20200108200545-475eaeb16496
	github.com/bicomsystems/go-libzfs v0.3.3 // indirect

	github.com/hashicorp/go-multierror v1.0.0
	// zfs 7.x bindings
	github.com/jsirianni/go-libzfs v1.0.0-zfs-7-zfsmon
	github.com/pkg/errors v0.9.1
	github.com/spf13/cobra v0.0.6
)
