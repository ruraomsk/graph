package tester

import (
	"encoding/json"

	"github.com/ruraomsk/ag-server/pudge"
)

func getCross(region, area, id int) (pudge.Cross, error) {
	crs, err := dbb.Query("select state from public.\"cross\" where region=$1 and area=$2 and id=$3;", region, area, id)
	if err != nil {
		return pudge.Cross{}, err
	}
	var cross pudge.Cross
	for crs.Next() {
		var state []byte
		crs.Scan(&state)
		err = json.Unmarshal(state, &cross)
		if err != nil {
			return pudge.Cross{}, err
		}
	}
	return cross, nil
}
