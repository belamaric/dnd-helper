package main

import (
	"flag"
	"log"
	"os"
)

var verbose bool
func main() {
	var check, encounter, addr, root string
	flag.BoolVar(&verbose, "v", false, "Verbose mode")
	flag.StringVar(&check, "c", "", "Check the XML file for unparsed XML")
	flag.StringVar(&encounter, "e", "", "Encounter YAML file")
	flag.StringVar(&addr, "s", "", "Start server on specified address")
	flag.StringVar(&root, "d", "", "root directory that contains data and html subdirs")

	flag.Parse()

	if addr != "" {
		es, err := NewEncounterServer(addr, root)
		if err != nil {
			log.Printf("ERROR: Could not create server: %s", err)
			os.Exit(1)
		}
		err = es.Serve()
		if err != nil {
			log.Printf("ERROR: Could not start http server: %s", err)
			os.Exit(1)
		}
	}

	if encounter != "" {
		f, err := os.Open(encounter)
		if err != nil {
			log.Printf("ERROR: Could not open encounter file: %s", err)
			os.Exit(1)
		}
		e, err := NewEncounterFromYaml(f)
		if err != nil {
			log.Printf("ERROR: Could not load encounter: %s", err)
			os.Exit(1)
		}
		err = e.Load()
		if err != nil {
			log.Printf("ERROR: Could not load encounter: %s", err)
			os.Exit(1)
		}
		err = e.Print(os.Stdout)
		if err != nil {
			log.Printf("ERROR: Could not print encounter: %s", err)
			os.Exit(1)
		}
		return
	}

	if check != "" {
		c, err := LoadCompendium(check)
		if err != nil {
			log.Printf("ERROR: Could not load monsters: %s", err)
			os.Exit(1)
		}

		checkXml(c)
		return
	}

	flag.Usage()
}


func checkXml(c *Compendium) {
	for _, m := range c.Monsters {
		if len(m.Extras) > 0 {
			log.Printf("Monster %q has unparsed data: %v", m.Name, m.Extras)
		}
	}
}

