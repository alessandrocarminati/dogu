package main

import (
	"net/http"
	"strconv"
	"strings"
	"os/exec"
	"os"
	logx "log"
	"fmt"
	"runtime"
	b64 "encoding/base64"
	lt "github.com/jweslley/localtunnel"
	"github.com/sevlyar/go-daemon"
)

var (
	outbg []byte
	busy bool = false
	log []string;
)



func HelloServer(w http.ResponseWriter, req *http.Request) {
	log=append(log,"HelloServer ")
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte("Service is alive.\n"))
}

func cmd_foreground(w http.ResponseWriter, req *http.Request){
	w.Header().Set("Content-Type", "text/plain")
	cmd := req.URL.Query().Get("cmd")
	log=append(log,"cmd_foreground "+cmd)
	outs, err := exec.Command(cmd).Output()
	if err != nil {
		outs = []byte("cmd failed " + err.Error())
		}
	w.Write([]byte(outs))
}

func exec_gosub(cmd string){
	var err error
	busy=true
	log=append(log,"exec_gosub "+ cmd)
        outbg, err = exec.Command(cmd).Output()
	if err!=nil {
		log=append(log,"exec_gosub error:"+ err.Error())
		}
	busy=false
}

func cmd_background(w http.ResponseWriter, req *http.Request){
	out:="busy"
        w.Header().Set("Content-Type", "text/plain")
        cmd := req.URL.Query().Get("cmd")
	log=append(log,"cmd_background "+ cmd)
	if ! busy {
		go exec_gosub(cmd)
		out="success"
		}
        w.Write([]byte(out))
}
func cmd_background_check(w http.ResponseWriter, req *http.Request){
	log=append(log,"cmd_backgroundc ")
	out:="dummy"
        w.Header().Set("Content-Type", "text/plain")
	if busy {
		out="working"
		} else {
			out=string(outbg)
			}
        w.Write([]byte(out))
}
func upd_script(w http.ResponseWriter, req *http.Request){

	out:="ok"
	b64pl:=req.URL.Query().Get("b64pl")
	name:=req.URL.Query().Get("name")
        log=append(log,"upd_script "+ name)
	f, err := os.Create(name)
	if err != nil {
		log=append(log,"upd_script create:"+err.Error())
		out=string(err.Error())
		} else {
			defer f.Close()
			sDec, _ := b64.StdEncoding.DecodeString(b64pl)
			_, err = f.Write(sDec)
			if err != nil {
				log=append(log,"upd_script write:"+err.Error())
				out=string(err.Error())
				}else {
					var err error
					if runtime.GOOS != "windows" {
						err = os.Chmod(name, 0777);
						}
					if err != nil {
						log=append(log,"upd_script chmod:"+err.Error())
						out=string(err.Error())
						}
					}

			}
        w.Write([]byte(out))
}
func getlog(w http.ResponseWriter, req *http.Request){
	log=append(log,"getlog ")
        w.Header().Set("Content-Type", "text/plain")
        w.Write([]byte(strings.Join(log[:], "\n")))
}



func main() {
/*
	conf_host:="https://localtunnel.me"
	conf_local:="localhost"
	conf_subdomain:="antani"
	conf_port:=8080
*/
	cntxt := &daemon.Context{
		PidFileName: "dogu.pid",
		PidFilePerm: 0644,
		LogFileName: "dogu.log",
		LogFilePerm: 0640,
		WorkDir:     "./",
		Umask:       027,
		Args:        os.Args,
		}



	conf, err := args_parse(cmd_line_item_init())
	if err!=nil {
		if err.Error() != "dummy"{
			fmt.Println(err.Error())
			}
		print_help(cmd_line_item_init());
		os.Exit(-1)
		}

	d, err := cntxt.Reborn()
	if err != nil {
		logx.Fatal("Unable to run: ", err)
		}

	if d != nil {
		fmt.Println("dogu started")
		return
		}
	defer cntxt.Release()

	c := lt.NewClient(conf.Host)
	t := c.NewTunnel(conf.Target, conf.Port)
	t.OpenAs(conf.Request_dom)
	logx.Printf("your url is: %s\n", t.URL())

	mux := http.NewServeMux()
	fileServer := http.FileServer(http.Dir("./"))
	mux.Handle("/f/", http.StripPrefix("/f", fileServer))
	mux.HandleFunc("/", HelloServer)
	mux.HandleFunc("/hello", HelloServer)
	mux.HandleFunc("/cmd_fore", cmd_foreground)
	mux.HandleFunc("/cmd_back", cmd_background)
	mux.HandleFunc("/cmd_backc", cmd_background_check)
	mux.HandleFunc("/upd_script", upd_script)
	mux.HandleFunc("/getlog", getlog)

	err = http.ListenAndServe(":"+strconv.Itoa(conf.Port), mux)
//    err := http.ListenAndServeTLS(":443", "server.crt", "server.key", nil)
	if err != nil {
		panic("ListenAndServe: ")
		}
}
