package sqlgraph

import (
	"database/sql"
	"fmt"

	"github.com/ruraomsk/ag-server/pudge"
)

//Open - начало работы с графовой подстемой
func Open(dbb *sql.DB, dbbb *sql.DB) error {
	db = dbb
	dbv = dbbb
	regions = make(map[int]*oneRegion)
	regs, err := db.Query("select distinct region from public.region order by region;")
	if err != nil {
		return err
	}
	defer regs.Close()
	for regs.Next() {
		var region int
		regs.Scan(&region)
		or := oneRegion{region: region, vertexs: make(map[uint64]*extVertex), freePin: 0}
		err := or.loadGraph()
		if err != nil {
			return err
		}
		regions[region] = &or
	}
	return nil
}

//GetListRegions возвращает список регионов для которых есть граф
func GetListRegions() []int {
	result := make([]int, 0)
	for _, v := range regions {
		if v.load {
			result = append(result, v.region)
		}
	}
	return result
}

//Save - сохраняет граф одного региона
func Save(region int) error {
	or, is := regions[region]
	if !is {
		return fmt.Errorf("нет такого региона %d", region)
	}
	if or.modify {
		err := or.saveGraph()
		if err != nil {
			return err
		}
		or.modify = false
	}
	return nil
}

//SaveAll - сохраняет все изменения в БД
func SaveAll() error {
	for _, v := range regions {
		if v.modify {
			err := v.saveGraph()
			if err != nil {
				return err
			}
			v.modify = false
		}
	}
	return nil
}

//AddVertex - добавляет перекресток в граф
func AddVertex(cross pudge.Cross) error {
	seek := Vertex{Region: cross.Region, Area: cross.Area, ID: cross.ID}
	oneR, is := regions[cross.Region]
	if !is {
		return fmt.Errorf("нет такого региона %d", cross.Region)
	}
	_, is = oneR.vertexs[seek.getUID()]
	if is {
		return fmt.Errorf("уже есть такой перекресток %s", seek.getCross())
	}
	seek.Dgis = cross.Dgis
	seek.Name = cross.Name
	seek.Scale = cross.Scale
	ext := extVertex{vertex: seek, ways: make(map[uint64]*Way)}
	oneR.vertexs[seek.getUID()] = &ext
	oneR.modify = true
	return nil
}

//DeleteVertex - удаляет перекресток из графа
func DeleteVertex(cross pudge.Cross) error {
	seek := Vertex{Region: cross.Region, Area: cross.Area, ID: cross.ID}
	oneR, is := regions[cross.Region]
	if !is {
		return fmt.Errorf("нет такого региона %d", cross.Region)
	}
	_, is = oneR.vertexs[seek.getUID()]
	if !is {
		return fmt.Errorf("нет перекрестка %s", seek.getCross())
	}
	oneR.modify = true
	delete(oneR.vertexs, seek.getUID())
	//Теперь бежим по все вершинам и удаляем пути на удаленную вершину
	for _, v := range oneR.vertexs {
		delete(v.ways, seek.getUID())
	}
	return nil
}

//AddPoint - добавляет точку в регионе
func AddPoint(region int, number int, position string, name string) error {
	np := Vertex{Region: region, Area: 0, ID: number, Dgis: position, Name: name, Scale: 1}
	oneR, is := regions[np.Region]
	if !is {
		return fmt.Errorf("нет такого региона %d", np.Region)
	}
	_, is = oneR.vertexs[np.getUID()]
	if is {
		return fmt.Errorf("уже есть такая точка %s", np.getCross())
	}
	ext := extVertex{vertex: np, ways: make(map[uint64]*Way)}
	oneR.vertexs[np.getUID()] = &ext
	oneR.modify = true
	return nil
}

//DeletePoint удаляет точку в регионе
func DeletePoint(region int, number int) error {
	np := Vertex{Region: region, Area: 0, ID: number}
	oneR, is := regions[np.Region]
	if !is {
		return fmt.Errorf("нет такого региона %d", np.Region)
	}
	_, is = oneR.vertexs[np.getUID()]
	if !is {
		return fmt.Errorf("нет такой точки %s", np.getCross())
	}
	oneR.modify = true
	delete(oneR.vertexs, np.getUID())
	//Теперь бежим по все вершинам и удаляем пути на удаленную вершину
	for _, v := range oneR.vertexs {
		delete(v.ways, np.getUID())
	}
	return nil
}

//GetNumberPoint - возвращает номер для точки данного региона
func GetNumberPoint(region int) (int, error) {
	oneR, is := regions[region]
	if !is {
		return 0, fmt.Errorf("нет такого региона %d", region)
	}
	number := 0
	for _, v := range oneR.vertexs {
		if v.vertex.Area == 0 {
			if v.vertex.ID > number {
				number = v.vertex.ID
			}
		}
	}
	return number + 1, nil
}
