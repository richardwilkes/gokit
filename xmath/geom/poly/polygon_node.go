package poly

import "github.com/richardwilkes/toolbox/xmath/geom"

type vertexNode struct {
	pt   geom.Point
	next *vertexNode
}

type polygonNode struct {
	left   *vertexNode
	right  *vertexNode
	next   *polygonNode
	proxy  *polygonNode
	active bool
}

func (p *polygonNode) addLeft(pt geom.Point) {
	p.proxy.left = &vertexNode{
		pt:   pt,
		next: p.proxy.left,
	}
}

func (p *polygonNode) addRight(pt geom.Point) {
	v := &vertexNode{pt: pt}
	p.proxy.right.next = v
	p.proxy.right = v
}

func (p *polygonNode) mergeLeft(other, list *polygonNode) {
	if p.proxy != other.proxy {
		p.proxy.right.next = other.proxy.left
		other.proxy.left = p.proxy.left
		for target := p.proxy; list != nil; list = list.next {
			if list.proxy == target {
				list.active = false
				list.proxy = other.proxy
			}
		}
	}
}

func (p *polygonNode) mergeRight(other, list *polygonNode) {
	if p.proxy != other.proxy {
		other.proxy.right.next = p.proxy.left
		other.proxy.right = p.proxy.right
		for target := p.proxy; list != nil; list = list.next {
			if list.proxy == target {
				list.active = false
				list.proxy = other.proxy
			}
		}
	}
}

func (p *polygonNode) generate() Polygon {
	contourCount := 0
	ptCounts := make([]int, 0, 32)

	// Count the points of each contour and disable any that don't have
	// enough points.
	for poly := p; poly != nil; poly = poly.next {
		if poly.active {
			var prev *vertexNode
			ptCount := 0
			for v := poly.proxy.left; v != nil; v = v.next {
				if prev == nil || prev.pt != v.pt {
					ptCount++
				}
				prev = v
			}
			if ptCount > 2 {
				ptCounts = append(ptCounts, ptCount)
				contourCount++
			} else {
				poly.active = false
			}
		}
	}
	if contourCount == 0 {
		return Polygon{}
	}

	// Create the polygon
	result := make([]Contour, contourCount)
	ci := 0
	for poly := p; poly != nil; poly = poly.next {
		if poly.active { //nolint:gocritic
			var prev *vertexNode
			result[ci] = make([]geom.Point, ptCounts[ci])
			v := len(result[ci]) - 1
			for vtx := poly.proxy.left; vtx != nil; vtx = vtx.next {
				if prev == nil || prev.pt != vtx.pt {
					result[ci][v] = vtx.pt
					v--
				}
				prev = vtx
			}
			ci++
		}
	}
	return result
}