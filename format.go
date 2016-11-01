// Copyright 2014 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package trace

type ViewerData struct {
	Events   []*ViewerEvent         `json:"traceEvents"`
}

type ViewerEvent struct {
	Name       string      `json:"name,omitempty"`
	Phase      string      `json:"ph"`
	Categories string      `json:"cat,omitempty"`
	Scope      string      `json:"s,omitempty"`
	Time       float64     `json:"ts"`
	Dur        float64     `json:"dur,omitempty"`
	Pid        uint64      `json:"pid"`
	Tid        uint64      `json:"tid"`
	ID         uint64      `json:"id,omitempty"`
	Stack      int         `json:"sf,omitempty"`
	EndStack   int         `json:"esf,omitempty"`
	Arg        interface{} `json:"args,omitempty"`
}
