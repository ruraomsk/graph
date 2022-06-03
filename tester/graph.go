package tester

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/ruraomsk/ag-server/pudge"
	"github.com/ruraomsk/graph/setup"
	"github.com/ruraomsk/graph/sqlgraph"
)

var (
	dbb      *sql.DB
	err      error
	region   = 1
	oneCross pudge.Cross
	twoCross pudge.Cross
)

func getDB() {
	dbinfo := fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable",
		setup.Set.DataBase.Host, setup.Set.DataBase.User,
		setup.Set.DataBase.Password, setup.Set.DataBase.DBname)
	for {
		dbb, err = sql.Open("postgres", dbinfo)

		if err != nil {
			fmt.Printf("Запрос на открытие %s %s\n", dbinfo, err.Error())
			time.Sleep(time.Second * 10)
			continue
		}
		if err = dbb.Ping(); err != nil {
			fmt.Printf("Ping %s\n", err.Error())
			time.Sleep(time.Second * 10)
			continue
		}
		break
	}
}
func FullTest() {
	getDB()
	// err = sqlgraph.CreateGraph(region, dbb)
	// if err != nil {
	// 	fmt.Println(err.Error())
	// 	return
	// }
	err = sqlgraph.Open(dbb, dbb)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	// fmt.Println(sqlgraph.GetListRegions())
	number, err := sqlgraph.GetNumberPoint(region)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println(number)
	err = sqlgraph.AddPoint(region, 1, "0.11111,0.333333", "Point1")
	if err != nil {
		fmt.Println(err.Error())
	}
	err = sqlgraph.AddPoint(region, 2, "0.999999,0.5555555", "Point2")
	if err != nil {
		fmt.Println(err.Error())
	}
	oneCross, err = getCross(region, 1, 1)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	twoCross, err = getCross(region, 1, 7)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	err = sqlgraph.AddWay(oneCross, twoCross, 1100)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	err = sqlgraph.AddWay(oneCross, twoCross, 22222)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	err = sqlgraph.AddWayToPoint(oneCross, 1, 1555)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	err = sqlgraph.AddWayToPoint(twoCross, 2, 19555)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	err = sqlgraph.AddWayFromPoint(2, twoCross, 1999)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	err = sqlgraph.AddWayFromPoint(1, oneCross, 2999)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	err = sqlgraph.AddWay(twoCross, oneCross, 11000)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	// err = sqlgraph.DeleteVertex(twoCross)
	// if err != nil {
	// 	fmt.Println(err.Error())
	// 	return
	// }
	// err = sqlgraph.DeletePoint(region, 2)
	// if err != nil {
	// 	fmt.Println(err.Error())
	// 	return
	// }
	list, err := sqlgraph.ReadAllPoints(region)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println(list)

	// err = sqlgraph.DeleteWay(twoCross, oneCross)
	// if err != nil {
	// 	fmt.Println(err.Error())
	// 	return
	// }
	err = sqlgraph.Save(region)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
}
