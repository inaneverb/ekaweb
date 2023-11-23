module github.com/inaneverb/ekaweb/middleware/otel/v2

go 1.21.0

require (
	github.com/inaneverb/ekaweb/v2 v2.0.5
	go.opentelemetry.io/otel v1.17.0
	go.opentelemetry.io/otel/trace v1.17.0
)

require (
	github.com/cespare/xxhash/v2 v2.2.0 // indirect
	github.com/inaneverb/ekacore/ekaarr/v4 v4.0.0 // indirect
	github.com/inaneverb/ekacore/ekaext/v4 v4.0.0 // indirect
	github.com/inaneverb/ekacore/ekaunsafe/v4 v4.0.0 // indirect
)

retract (
	v2.0.0 // Bug: Request, Response are not included to span if required
	v2.0.3 // Bug: Incorrect re-check changed HTTP method, path; fixed in 2.0.4
)
