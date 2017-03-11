package main

import (
	"flag"
	"log"
	"os"
)

var verbose bool
func main() {
	var check, encounter string
	flag.BoolVar(&verbose, "v", false, "Verbose mode")
	flag.StringVar(&check, "c", "", "Check the XML file for unparsed XML")
	flag.StringVar(&encounter, "e", "", "Encounter YAML file")

	flag.Parse()

	if encounter != "" {
		e, err := LoadEncounter(encounter)
		if err != nil {
			log.Printf("ERROR: Could not load encounter: %s", err)
			os.Exit(1)
		}
		err = e.Print(os.Stdout)
		if err != nil {
			log.Printf("ERROR: Could not print encounter: %s", err)
			os.Exit(2)
		}
		return
	}

	if check != "" {
		c, err := LoadCompendium(check)
		if err != nil {
			log.Printf("ERROR: Could not load monsters: %s", err)
			os.Exit(3)
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

