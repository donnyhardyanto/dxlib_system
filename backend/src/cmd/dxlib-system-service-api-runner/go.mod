module github.com/donnyhardyanto/dxlib-system/service-api-runner

go 1.24.5

require (
	github.com/donnyhardyanto/dxlib v1.71.0
	github.com/donnyhardyanto/dxlib_module v1.49.0
	github.com/pkg/errors v0.9.1
	github.com/donnyhardyanto/dxlib-system/common v0.0.0-00010101000000-000000000000
	golang.org/x/image v0.27.0
)

replace dxlib_module => ../../library/dxlib_module

replace dxlib => ../../library/dxlib

replace github.com/donnyhardyanto/dxlib-system/common => ../../library/dxlib-system-common
