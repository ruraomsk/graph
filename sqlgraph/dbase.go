package sqlgraph

import (
	"encoding/json"
	"fmt"
)

func (o *oneRegion) loadGraph() error {
	vs, err := dbv.Query("select state from public.vertex where region=$1;", o.region)
	o.load = false
	if err != nil {
		return err
	}
	defer vs.Close()
	for vs.Next() {
		var state []byte
		var v Vertex
		vs.Scan(&state)
		err = json.Unmarshal(state, &v)
		if err != nil {
			return err
		}
		e := new(extVertex)
		e.vertex = v
		e.ways = make(map[uint64]*Way)
		o.vertexs[v.getUID()] = e
		if v.Area == 0 {
			//Если область == 0 то это точка входа/выхода
			if o.freePin < v.ID {
				o.freePin = v.ID
			}
		}
		o.load = true
	}
	vs.Close()
	if !o.load {
		return nil
	}
	ws, err := dbv.Query("select info from public.ways where region=$1", o.region)
	if err != nil {
		return err
	}
	defer ws.Close()
	for ws.Next() {
		var info []byte
		w := new(Way)
		ws.Scan(&info)
		err = json.Unmarshal(info, &w)
		if err != nil {
			return err
		}
		e, is := o.vertexs[w.Source]
		if !is {
			return fmt.Errorf("нет такого перекрестка %s", w.getCrossSource())
		}
		e.ways[w.Target] = w
	}
	ws.Close()
	return nil
}
func (o *oneRegion) saveGraph() error {
	_, err := dbv.Exec("delete from public.vertex where region=$1;", o.region)
	if err != nil {
		return err
	}
	_, err = dbv.Exec("delete from public.ways where region=$1;", o.region)
	if err != nil {
		return err
	}
	// Выгружаем вершины
	for _, v := range o.vertexs {
		state, err := json.Marshal(v.vertex)
		if err != nil {
			return err
		}
		_, err = dbv.Exec("insert into public.vertex (region,uids,state) values ($1,$2,$3);", o.region, v.vertex.getUID(), state)
		if err != nil {
			return err
		}
		// Выгружаем ребра выходящие
		for _, w := range v.ways {
			info, err := json.Marshal(w)
			if err != nil {
				return err
			}
			_, err = dbv.Exec("insert into public.ways (region,source,target,info) values ($1,$2,$3,$4);", o.region, w.Source, w.Target, info)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
