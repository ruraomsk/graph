package tester

import (
	"database/sql"
	"fmt"
	"math/rand"
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
	rnd      = rand.New(rand.NewSource(time.Now().Unix()))
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
	err = sqlgraph.CreateGraph(region, dbb)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	err = sqlgraph.Open(dbb, dbb)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println("Data base loaded")
	vl, err := sqlgraph.ReadAllVertex(region)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println("Start make ways")
	for pos, v := range vl {
		oneCross, err = getCross(v.Region, v.Area, v.ID)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		for i := 0; i < 4; i++ {
			j := getPos(len(vl), pos)
			twoCross, _ = getCross(vl[j].Region, vl[j].Area, vl[j].ID)
			lens := 0
			for lens <= 0 {
				lens = rnd.Intn(3000)
			}
			err = sqlgraph.AddWay(oneCross, twoCross, 1, 2, lens, 100)
			if err != nil {
				fmt.Println(err.Error())
				return
			}
			err = sqlgraph.AddWay(twoCross, oneCross, 3, 4, lens, 100)
			if err != nil {
				fmt.Println(err.Error())
				return
			}
		}
	}
	fmt.Println("\nCreates ways end")
	sqlgraph.MakeCorols(region, vl[0])
	fmt.Println("Makers color end")
	fmt.Println(sqlgraph.ListUnlinked(region))
	fmt.Println(vl[0])
	fmt.Println(vl[len(vl)-1])
	fmt.Println(sqlgraph.MakePath(region, vl[0], vl[len(vl)-1]))
	fmt.Println("Make Path end")
	err = sqlgraph.Save(region)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
}
func getPos(max, pos int) int {
	r := pos
	for r == pos {
		r = rnd.Intn(max - 1)
	}
	return r
}
