package sqlgraph

func MakeCorols(region int, v Vertex) {
	r := regions[region]
	var queue []uint64
	for _, v := range r.vertexs {
		v.color = 0
	}
	queue = append(queue, v.getUID())
	for len(queue) != 0 {
		next := queue[0]
		queue = queue[1:]
		source := r.vertexs[next]
		for _, w := range source.ways {
			vtarget := r.vertexs[w.Target]
			if vtarget.color == 0 {
				queue = append(queue, w.Target)
				vtarget.color = 1
			}
		}
	}
}
func GetFull(region int) bool {
	r := regions[region]
	count := 0
	for _, v := range r.vertexs {
		if v.color == 0 {
			count++
		}
	}
	return count == 0
}
func ListUnlinked(region int) []Vertex {
	result := make([]Vertex, 0)
	r := regions[region]
	for _, v := range r.vertexs {
		if v.color == 0 {
			result = append(result, v.vertex)
		}
	}
	return result

}
