module github.com/inaneverb/ekaweb/v2

go 1.21

require (
	github.com/cespare/xxhash/v2 v2.2.0
	github.com/inaneverb/ekacore/ekaarr/v4 v4.0.0
	github.com/inaneverb/ekacore/ekaunsafe/v4 v4.0.0
)

require (
	github.com/inaneverb/ekacore/ekaext/v4 v4.0.0 // indirect
)

retract (
	[v1.1.0, v1.9.9] // breaking changes in major, use v2 instead
)