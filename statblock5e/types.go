package main

import (
	"io"
	"encoding/xml"
	"strings"
	"text/template"
)

type Compendium struct {
	XMLName xml.Name `xml:"compendium"`
	Monsters []Monster `xml:"monster"`
}

type Trait struct {
	XMLName xml.Name
	Name string `xml:"name"`
	Text []string `xml:"text"`
	Attack []string `xml:"attack"`
}

func (t Trait) FormattedText() []string {
	f := make([]string, len(t.Text))
	for i, s := range t.Text {
		f[i] = strings.Replace(s, "Hit:", "<i>Hit:</i>", -1)
	}
	return f
}

type Monster struct {
	XMLName xml.Name `xml:"monster"`
	Name string `xml:"name"`
	Size string `xml:"size"`
	Type string `xml:"type"`
	Alignment string `xml:"alignment"`
	Ac string `xml:"ac"`
	Hp string `xml:"hp"`
	Speed string `xml:"speed"`
	Str string `xml:"str"`
	Dex string `xml:"dex"`
	Con string `xml:"con"`
	Int string `xml:"int"`
	Wis string `xml:"wis"`
	Cha string `xml:"cha"`
	Save string `xml:"save"`
	Skill string `xml:"skill"`
	Vulnerabilities string `xml:"vulnerable"`
	Resistances string `xml:"resist"`
	DamageImmunity string `xml:"immune"`
	ConditionImmunity string `xml:"conditionImmune"`
	Senses string `xml:"senses"`
	Passive string `xml:"passive"`
	Languages string `xml:"languages"`
	Cr string `xml:"cr"`

	Traits []Trait `xml:"trait"`
	Actions []Trait `xml:"action"`
	Reactions []Trait `xml:"reaction"`
	Legendary []Trait `xml:"legendary"`

	Spells string `xml:"spells"` // included in text
	Slots string `xml:"slots"` // included in text

	Description string `xml:"description"`

       	Extras []struct {
       	     XMLName xml.Name
       	     Content string `xml:",innerxml"`
        } `xml:",any"`
}

func (m Monster) SizeName() (string) {
	switch m.Size {
	case "G":
		return "Gargantuan"
	case "H":
		return "Huge"
	case "L":
		return "Large"
	case "M":
		return "Medium"
	case "S":
		return "Small"
	case "T":
		return "Tiny"
	default:
		return m.Size
	}
}

func (m Monster) Subtitle() (string) {
	return m.SizeName() + " " + m.Type + ", " + m.Alignment
}

func (m Monster) EncodeStatBlock(w io.Writer) error {
	tmplTxt := `
<stat-block>
 <creature-heading>
  <h1>{{.Name}}</h1>
  <h2>{{.Subtitle}}</h2>
 </creature-heading>
 <top-stats>
  <property-line>
   <h4>Armor Class</h4>
   <p>{{.Ac}}</p>
  </property-line>
  <property-line>
   <h4>Hit Points</h4>
   <p>{{.Hp}}</p>
  </property-line>
  <property-line>
   <h4>Speed</h4>
   <p>{{.Speed}}</p>
  </property-line>
  <abilities-block data-str="{{.Str}}"
                   data-dex="{{.Dex}}"
                   data-con="{{.Con}}"
                   data-int="{{.Int}}"
                   data-wis="{{.Wis}}"
                   data-cha="{{.Cha}}"></abilities-block>
{{with .Save}}
  <property-line>
   <h4>Saving Throws</h4>
   <p>{{.}}</p>
  </property-line>
{{end}}
{{with .Skill}}
  <property-line>
   <h4>Skills</h4>
   <p>{{.}}</p>
  </property-line>
{{end}}
{{with .Vulnerabilities}}
  <property-line>
   <h4>Damage Vulnerabilities</h4>
   <p>{{.}}</p>
  </property-line>
{{end}}
{{with .Resistances}}
  <property-line>
   <h4>Damage Resistances</h4>
   <p>{{.}}</p>
  </property-line>
{{end}}
{{with .DamageImmunity}}
  <property-line>
   <h4>Damage Immunities</h4>
   <p>{{.}}</p>
  </property-line>
{{end}}
{{with .ConditionImmunity}}
  <property-line>
   <h4>Condition Immunities</h4>
   <p>{{.}}</p>
  </property-line>
{{end}}
  <property-line>
   <h4>Senses</h4>
{{if .Senses}}
   <p>{{.Senses}}, passive Perception {{.Passive}}</p>
{{else}}
   <p>passive Perception {{.Passive}}</p>
{{end}}
  </property-line>
{{with .Languages}}
  <property-line>
   <h4>Languages</h4>
   <p>{{.}}</p>
  </property-line>
{{end}}
  <property-line>
   <h4>Challenge</h4>
   <p>{{.Cr}}</p>
  </property-line>
 </top-stats>

 {{range .Traits}}
 <property-block>
  <h4>{{.Name}}.</h4>
  {{range .FormattedText}}
  <p>{{.}}</p>
  {{end}}
 </property-block>
 {{end}}

 {{if len .Actions}}
  <h3>Actions</h3>
 {{range .Actions}}
  <property-block>
    <h4>{{.Name}}.</h4>
  {{range .FormattedText}}
  <p>{{.}}</p>
  {{end}}
  </property-block>
 {{end}}
 {{end}}

 {{if len .Reactions}}
  <h3>Reactions</h3>
 {{range .Reactions}}
  <property-block>
    <h4>{{.Name}}.</h4>
  {{range .FormattedText}}
  <p>{{.}}</p>
  {{end}}
  </property-block>
 {{end}}
 {{end}}

 {{if len .Legendary}}
  <h3>Legendary Actions</h3>
 {{range .Legendary}}
  <property-block>
    <h4>{{.Name}}.</h4>
  {{range .FormattedText}}
  <p>{{.}}</p>
  {{end}}
  </property-block>
 {{end}}
 {{end}}
 {{with .Description}}
  <property-block>
    <h4></h4>
    <p>{{.}}</p>
  </property-block>
 {{end}}
</stat-block>
`
	tmpl, err := template.New("stat-block").Parse(tmplTxt)
	if err != nil { return err }
	err = tmpl.Execute(w, m)
	if err != nil { return err }
	return nil
}
