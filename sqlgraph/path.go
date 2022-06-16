package sqlgraph

import (
	"math"
)

func MakePath(region int, vStart, vEnd Vertex) []Vertex {
	var queue []uint64
	r := regions[region]
	path := make(map[uint64]uint64)
	for _, v := range r.vertexs {
		v.value = math.MaxInt
		path[v.vertex.getUID()] = 0
	}
	r.vertexs[vStart.getUID()].value = 0
	queue = append(queue, vStart.getUID())
	for len(queue) != 0 {
		next := queue[0]
		queue = queue[1:]
		source := r.vertexs[next]
		for _, w := range source.ways {
			vtarget := r.vertexs[w.Target]
			if vtarget.value > source.value {
				path[vtarget.vertex.getUID()] = source.vertex.getUID()
				queue = append(queue, w.Target)
				vtarget.value = source.value
			}
		}
	}
	result := make([]Vertex, 0)
	if r.vertexs[vEnd.getUID()].value == math.MaxInt {
		return result
	}
	t := vEnd.getUID()
	d := make([]uint64, 0)
	for t != 0 {
		d = append(d, t)
		t = path[t]
	}
	for i := len(d) - 1; i >= 0; i-- {
		result = append(result, r.vertexs[d[i]].vertex)
	}
	return result

}
