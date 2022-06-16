package sqlgraph

func MakeCorols(region int, vStart Vertex) {
	r := regions[region]
	var queue []uint64
	for _, v := range r.vertexs {
		v.value = 0
	}
	queue = append(queue, vStart.getUID())
	for len(queue) != 0 {
		next := queue[0]
		queue = queue[1:]
		source := r.vertexs[next]
		for _, w := range source.ways {
			vtarget := r.vertexs[w.Target]
			if vtarget.value == 0 {
				queue = append(queue, w.Target)
				vtarget.value = 1
			}
		}
	}
}
func GetFull(region int) bool {
	r := regions[region]
	count := 0
	for _, v := range r.vertexs {
		if v.value == 0 {
			count++
		}
	}
	return count == 0
}
func ListUnlinked(region int) []Vertex {
	result := make([]Vertex, 0)
	r := regions[region]
	for _, v := range r.vertexs {
		if v.value == 0 {
			result = append(result, v.vertex)
		}
	}
	return result

}
