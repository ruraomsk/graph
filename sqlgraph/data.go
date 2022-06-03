package sqlgraph

import (
	"database/sql"
	"fmt"
)

//Vertex хранение вершины графа
type Vertex struct {
	Region int     `json:"region"` //Регион
	Area   int     `json:"area"`   //Район
	ID     int     `json:"id"`     //Номер перекрестка
	Dgis   string  `json:"dgis"`   //Координаты перекрестка
	Scale  float64 `json:"scale"`  //Масштаб
	Name   string  `json:"name"`
}

//Way описание ребра графа
type Way struct {
	Region int    `json:"region"` //Регион
	Source uint64 `json:"source"` //Источник код перекрестка
	Target uint64 `json:"target"` //Получатель уод перекрестка
	Start  point  `json:"starts"`
	Stop   point  `json:"stops"`
	Lenght int    `json:"lenght"`
	Time   int    `json:"time"`
}

//WayToWeb описание ребра для отражения в браузере
type WayToWeb struct {
	Region     int    `json:"region"`     //Регион
	SourceArea int    `json:"sourceArea"` //Источник код перекрестка
	SourceID   int    `json:"sourceID"`   //Источник код перекрестка
	TargetArea int    `json:"targetArea"` //Получатель уод перекрестка
	TargetID   int    `json:"targetID"`   //Получатель уод перекрестка
	Start      string `json:"starts"`
	Stop       string `json:"stops"`
	Lenght     int    `json:"lenght"`
	Time       int    `json:"time"`
}
type extVertex struct {
	vertex Vertex
	ways   map[uint64]*Way
}
type point struct {
	value string
}

type oneRegion struct {
	region  int
	load    bool
	modify  bool
	freePin int
	vertexs map[uint64]*extVertex
}

var regions map[int]*oneRegion
var db *sql.DB
var dbv *sql.DB

func (v *Vertex) getUID() uint64 {
	return uint64(v.Region<<32 + v.Area<<16 + v.ID)
}

// func setUID(region, area, id int) uint64 {
// 	return uint64(region<<32 + area<<16 + id)
// }
func (v *Vertex) getCross() string {
	return fmt.Sprintf("регион %d область %d ДК %d", v.Region, v.Area, v.ID)
}
func (w *Way) getCrossSource() string {
	return fmt.Sprintf("регион %d область %d ДК %d", w.Source>>32&0xffff, w.Source>>16&0xffff, w.Source&0xffff)
}

// func (w *Way) getCrossTarget() string {
// 	return fmt.Sprintf("регион %d область %d ДК %d", w.Target>>32&0xffff, w.Target>>16&0xffff, w.Target&0xffff)
// }
func getArea(uid uint64) int {
	return int(uid >> 16 & 0xffff)
}
func getID(uid uint64) int {
	return int(uid & 0xffff)
}
