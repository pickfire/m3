// Copyright (c) 2018 Uber Technologies, Inc.
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

// This file was automatically generated by genny.
// Any changes will be lost if this file is regenerated.
// see https://github.com/mauricelam/genny

package storage

import (
	"github.com/m3db/m3/src/x/checked"
	"github.com/m3db/m3/src/x/ident"

	"github.com/cespare/xxhash/v2"
)

// Copyright (c) 2018 Uber Technologies, Inc.
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

// shardMapOptions provides options used when created the map.
type shardMapOptions struct {
	InitialSize int
	ByteOpts    checked.BytesOptions
}

// newShardMap returns a new byte keyed map.
func newShardMap(opts shardMapOptions) *shardMap {
	var (
		copyFn     shardMapCopyFn
		finalizeFn shardMapFinalizeFn
	)
	copyFn = func(k *ident.TagsOrID) *ident.TagsOrID {
		// Do not need a key copy pool since we
		// will be relying on the string table.
		// TODO: make clone not actually copy.
		return k
	}
	return _shardMapAlloc(_shardMapOptions{
		hash: func(id *ident.TagsOrID) shardMapHash {
			if id.Tags != nil {
				//b := id.Tags.ToID().Bytes()
				d := xxhash.New()
				for _, t := range id.Tags.Values() {
					d.Write(t.Name.Bytes())
					d.Write(t.Value.Bytes())
				}
				return shardMapHash(d.Sum64())
			}
			return shardMapHash(xxhash.Sum64(id.ID.Bytes()))
		},
		equals: func(x, y *ident.TagsOrID) bool {
			if x.Tags != nil && y.Tags != nil {
				return x.Tags.Equal(*y.Tags)
			}
			if x.ID != nil && y.ID != nil {
				return x.ID.Equal(y.ID)
			}
			if x.ID != nil {
				return ident.IDMatchesTags(x.ID, y.Tags)
			}
			return ident.IDMatchesTags(y.ID, x.Tags)
			// var xID []byte
			// if x.ID == nil {
			// 	xID = x.Tags.ToIDBytes()
			// } else {
			// 	xID = x.ID.Bytes()
			// }
			// var yID []byte
			// if y.ID == nil {
			// 	yID = y.Tags.ToIDBytes()
			// } else {
			// 	yID = y.ID.Bytes()
			// }
			// return bytes.Equal(xID, yID)
		},
		copy:        copyFn,
		finalize:    finalizeFn,
		initialSize: opts.InitialSize,
	})
}
