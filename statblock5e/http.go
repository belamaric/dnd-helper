package main

import (
	"io"
	"log"
	"net"
	"encoding/json"
	"net/http"
	"path/filepath"
	"strings"
)

type EncounterServer struct {
	addr string
	dir string
	compendiums map[string]*Compendium
	monsters map[string]*Monster
	server *http.ServeMux
}

func NewEncounterServer(addr, dir string) (*EncounterServer, error) {
	if addr == "" {
		addr = ":80"
	}
	es := &EncounterServer{addr: addr, dir: dir}
	es.monsters = make(map[string]*Monster)
	es.compendiums = make(map[string]*Compendium)

	files, err := filepath.Glob(dir + "/data/*.xml")
	if err != nil {
		return nil, err
	}
	for _, file := range files {
		c, err := LoadCompendium(file)
		if err != nil {
			log.Printf("ERROR: Skipping file %q because it failed loading: %q", file, err)
			continue
		}
		es.compendiums[c.Name] = c
		for _, m := range c.Monsters {
			m.Source = c.File
			es.monsters[m.Name + " (" + c.Name + ")"] = m
		}
	}

	return es, nil
}

func (es *EncounterServer) Serve() error {

	ln, err := net.Listen("tcp", es.addr)
	if err != nil {
		return err
	}

	es.server = http.NewServeMux()
	es.server.HandleFunc("/api/encounter/statblock5e", func(w http.ResponseWriter, r *http.Request) {
		es.handleEncounterStatBlock5e(w,r)
	})
	es.server.HandleFunc("/api/monsters", func(w http.ResponseWriter, r *http.Request) {
		es.handleMonsterList(w,r)
	})
	es.server.Handle("/", http.FileServer(http.Dir(es.dir + "/html")))

	return http.Serve(ln, es.server)
}

func (es *EncounterServer) handleEncounterStatBlock5e(w http.ResponseWriter, r *http.Request) {
	e, err := NewEncounterFromJson(r.Body)
	if err != nil {
		io.WriteString(w, err.Error())
		return
	}

	e.Fill(es.monsters)
	if err != nil {
		io.WriteString(w, err.Error())
		return
	}
	err = e.Print(w)
	if err != nil {
		io.WriteString(w, "\n\n" + err.Error())
	}
	return
}

func (es *EncounterServer) handleMonsterList(w http.ResponseWriter, r *http.Request) {
	compendium := strings.ToLower(r.FormValue("compendium"))
	search := strings.ToLower(r.FormValue("search"))
	var monsters []*Monster
	for _, c := range es.compendiums {
		if !strings.Contains(strings.ToLower(c.Name), compendium) {
			continue
		}
		for _, m := range c.Monsters {
			if strings.Contains(strings.ToLower(m.Name), search) {
				monsters = append(monsters, m)
			}
		}
	}
	str, err := json.Marshal(monsters)
	if err != nil {
		io.WriteString(w, err.Error())
		return
	}
	w.Header().Add(`Content-type`, `application/json`)
	io.WriteString(w, `{"monsters":`)
	w.Write(str)
	io.WriteString(w, "}")
}
