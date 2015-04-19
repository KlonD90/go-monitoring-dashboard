package main

import (
	"fmt"
	"net"
	"os"
	"github.com/shirou/gopsutil/mem"
	"time"
	"encoding/json"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

func CheckError(err error) {
	if err  != nil {
        fmt.Println("Error: " , err)
		os.Exit(0)
	}
}


type CurrentMemory struct {
	Total uint64
	Free uint64
	Used uint64
	UsedPercent float64
	SwapUsage uint64
}

type ActiveQuery struct {
	Time int
	Query string
}

func main() {
	v, _ := mem.VirtualMemory()
	s, _ := mem.SwapMemory()
	ServerAddr,err := net.ResolveUDPAddr("udp","127.0.0.1:10001")
	CheckError(err)
	Conn, err := net.DialUDP("udp", nil, ServerAddr)
	CheckError(err)
	defer Conn.Close()
	go func() {
		for {
			msg := CurrentMemory{v.Total, v.Free, v.Used,
				v.UsedPercent, s.Used}
			buf, err := json.Marshal(msg)
			CheckError(err)
			_, err = Conn.Write(buf)
			if err != nil {
				fmt.Println(msg, err)
			}
			time.Sleep(time.Second * 1)
		}
	}()
	go func() {
		db, err := sql.Open("mysql", "root:mendyaev@/mysql")
		CheckError(err)
		defer db.Close()
		stmt, err := db.Prepare("SELECT TIME, INFO "+
			"FROM information_schema.processlist "+
			"WHERE COMMAND = \"Query\" AND TIME > 0")
		var (
			timeQuery int
			query string
		)
		for {
			rows, err := stmt.Query()
			CheckError(err)
			for rows.Next() {
				err := rows.Scan(&timeQuery, &query)
				CheckError(err)
				fmt.Println("Time: ", timeQuery, ", query: ", query)
				msg := ActiveQuery{timeQuery, query}
				buf, err := json.Marshal(msg)
				CheckError(err)
				_, err = Conn.Write(buf)
				CheckError(err)
			}
			rows.Close()
			time.Sleep(time.Second * 1)
		}
	}()
	ServerAddrListen,err := net.ResolveUDPAddr("udp", ":10002")
	if err != nil {
		fmt.Println("Error: ", err)
		os.Exit(0)
	}
	ServerConn, err := net.ListenUDP("udp", ServerAddrListen)
	if err != nil {
		fmt.Println("Error: ", err)
		os.Exit(0)
	}
	defer ServerConn.Close()
	buf := make([]byte, 1024)
	for {
		n, addr, err := ServerConn.ReadFromUDP(buf)
		fmt.Println("Received ", string(buf[0:n]), " from ", addr)

		if err != nil {
			fmt.Println("Error: ", err)
		}
	}
}
