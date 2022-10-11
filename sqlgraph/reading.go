package sqlgraph

import (
	"fmt"

	"github.com/ruraomsk/ag-server/pudge"
)

// ReadAllVertex получить список всех перекрестков в графе
func ReadAllVertex(region int) ([]Vertex, error) {
	result := make([]Vertex, 0)
	oneR, is := regions[region]
	if !is {
		return result, fmt.Errorf("нет такого региона %d", region)
	}
	for _, v := range oneR.vertexs {
		if v.vertex.Area != 0 {
			result = append(result, v.vertex)
		}
	}
	return result, nil
}

// ReadAllPoints получить список всех точек в графе
func ReadAllPoints(region int) ([]Vertex, error) {
	result := make([]Vertex, 0)
	oneR, is := regions[region]
	if !is {
		return result, fmt.Errorf("нет такого региона %d", region)
	}
	for _, v := range oneR.vertexs {
		if v.vertex.Area == 0 {
			result = append(result, v.vertex)
		}
	}
	return result, nil
}

// GetWaysFromCross получить все пути из перекрестка
func GetWaysFromCross(cross pudge.Cross) ([]WayToWeb, error) {
	result := make([]WayToWeb, 0)
	_, sext, err := verifyCross(cross)
	if err != nil {
		return result, err
	}
	for _, v := range sext.ways {
		ww := WayToWeb{Region: v.Region, SourceArea: getArea(v.Source), SourceID: getID(v.Source), TargetArea: getArea(v.Target), TargetID: getID(v.Target),
			LineSource: v.LineSource, LineTarget: v.LineTarget,
			Start: v.Start, Stop: v.Stop, Lenght: v.Lenght, Time: v.Time}
		result = append(result, ww)
	}
	return result, nil
}

// GetWaysFromPoint получить все пути из точки
func GetWaysFromPoint(region int, number int) ([]WayToWeb, error) {
	result := make([]WayToWeb, 0)
	_, sext, err := verifyPoint(region, number)
	if err != nil {
		return result, err
	}
	for _, v := range sext.ways {
		ww := WayToWeb{Region: v.Region, SourceArea: getArea(v.Source), SourceID: getID(v.Source), TargetArea: getArea(v.Target), TargetID: getID(v.Target),
			Start: v.Start, Stop: v.Stop, Lenght: v.Lenght, Time: v.Time}
		result = append(result, ww)
	}
	return result, nil
}
