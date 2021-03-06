// Copyright ©2016-2020 by Richard A. Wilkes. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, version 2.0. If a copy of the MPL was not distributed with
// this file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// This Source Code Form is "Incompatible With Secondary Licenses", as
// defined by the Mozilla Public License, version 2.0.

package poly

type scanBeamTree struct {
	root    *scanBeamNode
	entries int
}

type scanBeamNode struct {
	y    float64
	less *scanBeamNode
	more *scanBeamNode
}

func (sbt *scanBeamTree) add(y float64) {
	sbt.addToScanBeamTreeAt(&sbt.root, y)
}

func (sbt *scanBeamTree) addToScanBeamTreeAt(node **scanBeamNode, y float64) {
	switch {
	case *node == nil:
		*node = &scanBeamNode{y: y}
		sbt.entries++
	case (*node).y > y:
		sbt.addToScanBeamTreeAt(&(*node).less, y)
	case (*node).y < y:
		sbt.addToScanBeamTreeAt(&(*node).more, y)
	default:
	}
}

func (sbt *scanBeamTree) buildScanBeamTable() []float64 {
	table := make([]float64, sbt.entries)
	if sbt.root != nil {
		sbt.root.buildScanBeamTableEntries(0, table)
	}
	return table
}

func (sbn *scanBeamNode) buildScanBeamTableEntries(index int, table []float64) int {
	if sbn.less != nil {
		index = sbn.less.buildScanBeamTableEntries(index, table)
	}
	table[index] = sbn.y
	index++
	if sbn.more != nil {
		index = sbn.more.buildScanBeamTableEntries(index, table)
	}
	return index
}
