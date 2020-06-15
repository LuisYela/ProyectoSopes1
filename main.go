package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"io/ioutil"
	"strconv"
	"strings"

	socketio "github.com/googollee/go-socket.io"
)

type objcpuUsage struct {
	Uso   string `json:"uso"`
	Busy  string `json:"busy"`
	Total string `json:"total"`
}
type objRamStats struct {
	Libre      string `json:"libre"`
	Disponible string `json:"disponible"`
	Total      string `json:"total"`
}
type objProcesosStats struct {
	Ejecutandose string `json:"ejecutandose"`
	Suspendidos  string `json:"suspendidos"`
	Detenidos    string `json:"detenidos"`
	Zombies      string `json:"zombies"`
	Total        string `json:"total"`
}
type objProcesosTodos struct {
	Procid      string `json:"procid"`
	ProcNombre  string `json:"procnombre"`
	ProcUsuario string `json:"procusuario"`
	ProcPid     string `json:"procpid"`
	ProcEstado  string `json:"procestado"`
	ProcRam     string `json:"procram"`
}

func main() {

	server, err := socketio.NewServer(nil)
	if err != nil {
		log.Fatal(err)
	}
	//sockets
	server.OnConnect("/", func(s socketio.Conn) error {
		s.SetContext("")
		fmt.Println("connected:", s.ID())
		return nil
	})
	server.OnEvent("/", "statscpu", func(s socketio.Conn) {
		go calculocpu(s)
	})
	server.OnEvent("/", "statsram", func(s socketio.Conn) {
		//
		go calculoram(s)
	})
	server.OnEvent("/", "statsproc", func(s socketio.Conn) {
		//
		go calculoNumProc(s)
	})
	server.OnEvent("/", "TextoEnviado", func(s socketio.Conn, msg string) {
		s.SetContext(msg)
		fmt.Println("Proceso que se decea borrar:", msg)
		s.Emit("statsMemory", msg)
	})
	server.OnEvent("/", "bye", func(s socketio.Conn) string {
		last := s.Context().(string)
		s.Emit("bye", last)
		s.Close()
		return last
	})
	server.OnError("/", func(s socketio.Conn, e error) {
		fmt.Println("meet error:", e)
	})
	server.OnDisconnect("/", func(s socketio.Conn, reason string) {
		fmt.Println("closed", reason)
	})
	go server.Serve()
	defer server.Close()

	http.Handle("/socket.io/", server)
	http.Handle("/", http.FileServer(http.Dir("./public")))
	log.Println("Serving at localhost:8000...")
	log.Fatal(http.ListenAndServe(":8000", nil))
}

func calculocpu(s socketio.Conn) {
	for {
		//
		idle0, total0 := getCPUSample()
		time.Sleep(3 * time.Second)
		idle1, total1 := getCPUSample()

		idleTicks := float64(idle1 - idle0)
		totalTicks := float64(total1 - total0)
		cpuUsage := 100 * (totalTicks - idleTicks) / totalTicks

		//fmt.Printf("CPU usage is %f%% [busy: %f, total: %f]\n", cpuUsage, totalTicks-idleTicks, totalTicks)
		///
		objcpu := objcpuUsage{Uso: strconv.FormatFloat(cpuUsage, 'g', 10, 64), Busy: strconv.FormatFloat(totalTicks-idleTicks, 'g', 10, 64), Total: strconv.FormatFloat(totalTicks, 'g', 10, 64)}
		s.Emit("statsMemory", objcpu)
	}
}
func calculoram(s socketio.Conn) {
	for {
		time.Sleep(5 * time.Second)
		//
		file_data, err := ioutil.ReadFile("/proc/meminfo")
		if err != nil {
			fmt.Println("Hubo un error")
		}
		lines := strings.Split(string(file_data), "\n")
		/*fmt.Println(string(file_data))
		for i := 0; i < 3; i++ {
			fmt.Println(lines[i])
		}*/
		lineauxT := strings.Split(lines[0], ":")
		lineT := strings.ReplaceAll(lineauxT[1], " ", "")
		lineT = strings.ReplaceAll(lineT, "kB", "")
		lineauxF := strings.Split(lines[1], ":")
		lineF := strings.ReplaceAll(lineauxF[1], " ", "")
		lineF = strings.ReplaceAll(lineF, "kB", "")
		lineauxA := strings.Split(lines[2], ":")
		lineA := strings.ReplaceAll(lineauxA[1], " ", "")
		lineA = strings.ReplaceAll(lineA, "kB", "")
		objram := objRamStats{Libre: lineF, Disponible: lineA, Total: lineT}
		//fmt.Println(lineF, lineA, lineT)
		s.Emit("statsram", objram)
	}
}
func calculoNumProc(s socketio.Conn) {
	for {
		//
		files, err := ioutil.ReadDir("/proc")
		if err != nil {
			fmt.Println("Hubo un error")
		}
		//fmt.Println(len(files))
		//fmt.Println(len(files) - 63)
		var objsprocesos []objProcesosTodos
		var (
			enejec   int = 0
			enidle   int = 0
			ensusp   int = 0
			enzombie int = 0
			tot      int = 0
		)
		for i := 0; i < len(files)-64; i++ {
			fmt.Println(files[i].Name())
			file_data, err := ioutil.ReadFile("/proc/" + files[i].Name() + "/status")
			if err != nil {
				fmt.Println("Hubo un error")
			}
			var (
				nameactual string
				pid        string
				ppid       string
				state      string
				uid        string
				ramact     string = "0"
			)
			lines := strings.Split(string(file_data), "\n")
			//fmt.Println(string(file_data))
			for j := 0; j < 11; j++ {
				actual := strings.Split(string(lines[j]), ":")
				if actual[0] == "Name" || actual[0] == "Pid" || actual[0] == "State" || actual[0] == "Uid" || actual[0] == "PPid" || actual[0] == "FDsize" {
					if actual[0] == "Name" {
						//fmt.Print(strings.ReplaceAll(actual[0], " ", ""))
						//fmt.Println(strings.ReplaceAll((strings.ReplaceAll(actual[1], " ", "")), "	", ""))
						nameactual = (strings.ReplaceAll((strings.ReplaceAll(actual[1], " ", "")), "	", ""))
					}
					if actual[0] == "Pid" {
						//fmt.Print(actual[0])
						//fmt.Println(strings.ReplaceAll((strings.ReplaceAll(actual[1], " ", "")), "	", ""))
						pid = (strings.ReplaceAll((strings.ReplaceAll(actual[1], " ", "")), "	", ""))
					}
					if actual[0] == "State" {
						//fmt.Print(actual[0])
						//fmt.Println(strings.ReplaceAll((strings.ReplaceAll(actual[1], " ", "")), "	", ""))
						state = (strings.ReplaceAll((strings.ReplaceAll(actual[1], "	", "")), " ", ""))
						if state == "S(sleeping)" {
							ensusp = ensusp + 1
						} else if state == "R(running)" {
							enejec = enejec + 1
						} else if state == "I(idle)" {
							enidle = enidle + 1
						} else {
							enzombie = enzombie + 1
						}
					}
					if actual[0] == "Uid" {
						//fmt.Println(actual[0])
						//fmt.Println(strings.ReplaceAll((strings.ReplaceAll(actual[1], " ", "")), "	", ""))
						uid = (strings.ReplaceAll((strings.ReplaceAll(actual[1], " ", "")), "	", ""))
						uid = string(uid[0])
						//fmt.Println(uid)
					}
					if actual[0] == "PPid" {
						//fmt.Print(actual[0])
						//fmt.Println(strings.ReplaceAll((strings.ReplaceAll(actual[1], " ", "")), "	", ""))
						ppid = (strings.ReplaceAll((strings.ReplaceAll(actual[1], " ", "")), "	", ""))
					}
					if actual[0] == "FDsize" {
						//fmt.Print(actual[0])
						//fmt.Println(strings.ReplaceAll((strings.ReplaceAll(actual[1], " ", "")), "	", ""))
						ramact = (strings.ReplaceAll((strings.ReplaceAll(actual[1], " ", "")), "	", ""))
					}
				}
			}
			objproc := objProcesosTodos{Procid: pid, ProcNombre: nameactual, ProcUsuario: uid, ProcPid: ppid, ProcEstado: state, ProcRam: ramact}
			objsprocesos = append(objsprocesos, objproc)
		}
		fmt.Println(objsprocesos)
		tot = ensusp + enejec + enidle + enzombie
		objprocstat := objProcesosStats{Ejecutandose: strconv.Itoa(enejec), Suspendidos: strconv.Itoa(ensusp), Detenidos: strconv.Itoa(enidle), Zombies: strconv.Itoa(enzombie), Total: strconv.Itoa(tot)}
		fmt.Println(objprocstat)
		/*lines := strings.Split(string(file_data), "\n")
		fmt.Println(string(file_data))
		for i := 0; i < 3; i++ {
			fmt.Println(lines[i])
		}
		lineauxT := strings.Split(lines[0], ":")
		lineT := strings.ReplaceAll(lineauxT[1], " ", "")
		lineT = strings.ReplaceAll(lineT, "kB", "")
		lineauxF := strings.Split(lines[1], ":")
		lineF := strings.ReplaceAll(lineauxF[1], " ", "")
		lineF = strings.ReplaceAll(lineF, "kB", "")
		lineauxA := strings.Split(lines[2], ":")
		lineA := strings.ReplaceAll(lineauxA[1], " ", "")
		lineA = strings.ReplaceAll(lineA, "kB", "")
		objram := objRamStats{Libre: lineF, Disponible: lineA, Total: lineT}
		//fmt.Println(lineF, lineA, lineT)*/
		s.Emit("statsproc", objprocstat)
		s.Emit("proclistado", objsprocesos)
		time.Sleep(10 * time.Second)
	}
}

func getCPUSample() (idle, total uint64) {
	contents, err := ioutil.ReadFile("/proc/stat")
	if err != nil {
		return
	}
	lines := strings.Split(string(contents), "\n")
	for _, line := range lines {
		fields := strings.Fields(line)
		if fields[0] == "cpu" {
			numFields := len(fields)
			for i := 1; i < numFields; i++ {
				val, err := strconv.ParseUint(fields[i], 10, 64)
				if err != nil {
					fmt.Println("Error: ", i, fields[i], err)
				}
				total += val // tally up all the numbers to get total ticks
				if i == 4 {  // idle is the 5th field in the cpu line
					idle = val
				}
			}
			return
		}
	}
	return
}
