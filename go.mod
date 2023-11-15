module github.com/inaneverb/ekaweb/v2

// =============================================================================
//  These directories are also a part of "ekaweb", not any internal separated
//  module: /ekaweb/private, /ekaweb/websocket, /ekaweb/middleware
//  DO NOT FORGET TO UPDATE IMPORTS INSIDE INNER DIRECTORIES, THAT ARE PART
//  OF EKAWEB PACKAGE, WHEN YOU GOING TO RELEASE NEXT MAJOR VERSION.
// =============================================================================

go 1.21

require (
	github.com/cespare/xxhash/v2 v2.2.0
	github.com/inaneverb/ekacore/ekaarr/v4 v4.0.0
	github.com/inaneverb/ekacore/ekaunsafe/v4 v4.0.0
)

require github.com/inaneverb/ekacore/ekaext/v4 v4.0.0 // indirect

retract (
	v2.0.0 // not all import paths updated, mix of v1, v2; unbuildable
	v2.0.2 // non-optimal way to initialize UKVS inside packed middleware
	[v2.0.0, v2.0.4] // segfault bug in UKVS getter, also placeholder panic
)
