package sqlgraph

import (
	"database/sql"
	"encoding/json"
	"fmt"

	_ "github.com/jackc/pgx"
	"github.com/ruraomsk/ag-server/pudge"
)

//_ "github.com/lib/pq"

// DROP TABLE IF EXISTS public.ways;
// DROP TABLE IF EXISTS public.vertex;

var err error
var create = `
CREATE TABLE IF NOT EXISTS public.ways
(
	region	integer     NOT NULL,
    source 	bigint	 	NOT NULL,
    target 	bigint	 	NOT NULL,
	info    jsonb 		NOT NULL DEFAULT '{}'::jsonb
)
WITH (
    OIDS = FALSE,
    autovacuum_enabled = TRUE
)
TABLESPACE pg_default;

ALTER TABLE IF EXISTS public.ways
    OWNER to postgres;

COMMENT ON TABLE public.ways
    IS 'Таблица ребер';

CREATE TABLE IF NOT EXISTS public.vertex
	(
		region	integer     NOT NULL,
		uids 	bigint	 	NOT NULL Primary key,
		state   jsonb 		NOT NULL DEFAULT '{}'::jsonb
	)
	WITH (
		OIDS = FALSE,
		autovacuum_enabled = TRUE
	)
	TABLESPACE pg_default;
	
	ALTER TABLE IF EXISTS public.vertex
		OWNER to postgres;
	
	COMMENT ON TABLE public.vertex
		IS 'Таблица вершин';
	
	
`

//CreateGraph запускается один раз для переноса всех перекрестков как вершины графа
func CreateGraph(region int, dbb *sql.DB) error {
	if err = dbb.Ping(); err != nil {
		return fmt.Errorf("ping %s", err.Error())
	}
	_, err = dbb.Exec(create)
	if err != nil {
		fmt.Println(err.Error())
		return fmt.Errorf("запрос на создание графа %s", err.Error())
	}

	_, err = dbb.Exec("delete from public.vertex where region=$1", region)
	if err != nil {
		return fmt.Errorf("delete records from vertex %s", err.Error())
	}
	_, err = dbb.Exec("delete from public.ways where region=$1", region)
	if err != nil {
		return fmt.Errorf("delete records from ways %s", err.Error())
	}
	crs, err := dbb.Query("select state from public.\"cross\" where region=$1;", region)
	if err != nil {
		return fmt.Errorf("read cross %s", err.Error())
	}
	var state []byte
	var cross pudge.Cross
	for crs.Next() {
		crs.Scan(&state)
		err = json.Unmarshal(state, &cross)
		if err != nil {
			return fmt.Errorf("unmarhal cross %s", err.Error())
		}
		v := Vertex{Region: cross.Region, Area: cross.Area, ID: cross.ID, Dgis: cross.Dgis, Name: cross.Name, Scale: cross.Scale}
		state, _ := json.Marshal(v)
		uid := v.getUID()
		_, err = dbb.Exec("INSERT INTO public.vertex (region,uids,state) VALUES ($1, $2,$3);", cross.Region, uid, state)
		if err != nil {

			return fmt.Errorf("запрос на вставку %s", err.Error())
		}
	}
	crs.Close()
	return nil
}
