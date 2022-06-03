package sqlgraph

import (
	"fmt"

	"github.com/ruraomsk/ag-server/pudge"
)

//AddWay добавляет связь между двумя перекрестками
func AddWay(source pudge.Cross, target pudge.Cross, lenght int) error {
	if source.Region != target.Region {
		return fmt.Errorf("разные регионы источника и цели")
	}
	oneR, sext, err := verifyCross(source)
	if err != nil {
		return err
	}
	_, tart, err := verifyCross(target)
	if err != nil {
		return err
	}
	way := Way{Region: source.Region, Source: sext.vertex.getUID(), Target: tart.vertex.getUID(),
		Start: point{value: sext.vertex.Dgis}, Stop: point{value: tart.vertex.Dgis}, Lenght: lenght}
	oneR.vertexs[way.Source].ways[way.Target] = &way
	return nil
}

//AddWayToPoint добавляет связь от перекрестка к точке
func AddWayToPoint(source pudge.Cross, number int, lenght int) error {
	oneR, sext, err := verifyCross(source)
	if err != nil {
		return err
	}
	_, tart, err := verifyPoint(source.Region, number)
	if err != nil {
		return err
	}
	way := Way{Region: source.Region, Source: sext.vertex.getUID(), Target: tart.vertex.getUID(),
		Start: point{value: sext.vertex.Dgis}, Stop: point{value: tart.vertex.Dgis}, Lenght: lenght}
	oneR.vertexs[way.Source].ways[way.Target] = &way
	oneR.modify = true
	return nil
}

//AddWayFromPoint добавляет связь от точки к перекрестку
func AddWayFromPoint(number int, target pudge.Cross, lenght int) error {
	oneR, sext, err := verifyPoint(target.Region, number)
	if err != nil {
		return err
	}
	_, tart, err := verifyCross(target)
	if err != nil {
		return err
	}
	way := Way{Region: target.Region, Source: sext.vertex.getUID(), Target: tart.vertex.getUID(),
		Start: point{value: sext.vertex.Dgis}, Stop: point{value: tart.vertex.Dgis}, Lenght: lenght}
	oneR.vertexs[way.Source].ways[way.Target] = &way
	oneR.modify = true
	return nil
}

//DeleteWay удаляет связь между двумя перекрестками
func DeleteWay(source pudge.Cross, target pudge.Cross) error {
	if source.Region != target.Region {
		return fmt.Errorf("разные регионы источника и цели")
	}
	oneR, sext, err := verifyCross(source)
	if err != nil {
		return err
	}
	_, tart, err := verifyCross(target)
	if err != nil {
		return err
	}
	delete(oneR.vertexs[sext.vertex.getUID()].ways, tart.vertex.getUID())
	oneR.modify = true
	return nil
}

//DeleteWayToPoint удаляет путь от перекрестка к точке
func DeleteWayToPoint(source pudge.Cross, number int) error {
	oneR, sext, err := verifyCross(source)
	if err != nil {
		return err
	}
	_, tart, err := verifyPoint(source.Region, number)
	if err != nil {
		return err
	}
	delete(oneR.vertexs[sext.vertex.getUID()].ways, tart.vertex.getUID())
	oneR.modify = true
	return nil
}

//DeleteWayFromPoint удаляет путь от точки к перекрестку
func DeleteWayFromPoint(number int, target pudge.Cross) error {
	oneR, sext, err := verifyPoint(target.Region, number)
	if err != nil {
		return err
	}
	_, tart, err := verifyCross(target)
	if err != nil {
		return err
	}
	delete(oneR.vertexs[sext.vertex.getUID()].ways, tart.vertex.getUID())
	oneR.modify = true
	return nil
}
func verifyCross(cross pudge.Cross) (*oneRegion, *extVertex, error) {
	seek := Vertex{Region: cross.Region, Area: cross.Area, ID: cross.ID}
	oneR, is := regions[cross.Region]
	if !is {
		return nil, nil, fmt.Errorf("нет такого региона %d", cross.Region)
	}
	ext, is := oneR.vertexs[seek.getUID()]
	if !is {
		return nil, nil, fmt.Errorf("нет перекрестка %s", seek.getCross())
	}
	return oneR, ext, nil
}
func verifyPoint(region, number int) (*oneRegion, *extVertex, error) {
	np := Vertex{Region: region, Area: 0, ID: number}
	oneR, is := regions[np.Region]
	if !is {
		return nil, nil, fmt.Errorf("нет такого региона %d", np.Region)
	}
	ext, is := oneR.vertexs[np.getUID()]
	if !is {
		return nil, nil, fmt.Errorf("нет такой точки %s", np.getCross())
	}
	return oneR, ext, nil
}
