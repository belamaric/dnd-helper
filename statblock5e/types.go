package main

import (
	"io"
	"io/ioutil"
	"encoding/xml"
	"fmt"
	"log"
	"strings"
	"text/template"

	yaml "gopkg.in/yaml.v2"
)

type Encounter struct {
	Name string `yaml:name`
	Source string `yaml:source`
	Monsters []*struct {
		Source string `yaml:source`
		Name string `yaml:name`
		Quantity int `yaml:quantity`
		Monster *Monster
	} `yaml:monsters`
}

func LoadEncounter(path string) (*Encounter, error) {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	e := &Encounter{}
	err = yaml.Unmarshal(b, e)
	if err != nil {
		return nil, err
	}

	log.Println("Encounter: ", e)

	if e.Source == "" {
		return nil, fmt.Errorf("No source specified for encounter.")
	}

	c, err := LoadCompendium(e.Source)
	if err != nil {
		return nil, err
	}

	sources := make(map[string]*Compendium)
	sources[e.Source] = c
	for _, m := range e.Monsters {
		s := e.Source
		if m.Source != "" {
			s = m.Source
		}
		if _, ok := sources[s]; !ok {
			cc, err := LoadCompendium(m.Source)
			if err != nil {
				return nil, err
			}
			sources[s] = cc
		}
		m.Monster = sources[s].FindMonster(m.Name)
		if m.Monster == nil {
			return nil, fmt.Errorf("Could not find %q in %q.", m.Name, s)
		}
	}
	return e, nil
}

type Compendium struct {
	XMLName xml.Name `xml:"compendium"`
	Monsters []Monster `xml:"monster"`
}

func LoadCompendium(path string) (*Compendium, error) {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("Could not load compendium from file %q: %s", path, err)
	}

	c := &Compendium{}
	err = xml.Unmarshal(b, c)
	if err != nil {
		return nil, err
	}

	return c, nil
}

func (e *Encounter) Print(w io.Writer) error {
	tmpl := template.New("page")
	tmpl.Funcs(template.FuncMap{"add": func(i, j int) int { return i+j }})
	tmpl.Funcs(template.FuncMap{"breakrow": func(i, j int) bool { return (i+1) %j == 0 }})
	tmpl.Funcs(template.FuncMap{"intarray": func(i, j int) []int {
		a := []int{}
		inc := 1
		if i >= j {
			inc = -1
		}
		for k := i; k != j; k = k + inc {
			a = append(a, k)
		}
		return a
		}})
	tmpl, err := tmpl.Parse(page)
	if err != nil { return err }
	err = tmpl.Execute(w, e)
	if err != nil { return err }
	return nil
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

func (m *Monster) SizeName() (string) {
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

func (m *Monster) ShortAc() (string) {
	s := strings.SplitN(m.Ac, " ", 2)
	if len(s) > 0 {
		return s[0]
	}
	return m.Ac
}

func (m *Monster) ShortHp() (string) {
	s := strings.SplitN(m.Hp, " ", 2)
	if len(s) > 0 {
		return s[0]
	}
	return m.Hp
}

func (m *Monster) Subtitle() (string) {
	return m.SizeName() + " " + m.Type + ", " + m.Alignment
}

func (c *Compendium) FindMonster(name string) *Monster {
	for _, v := range c.Monsters {
		if v.Name == name {
			return &v
		}
	}
	return nil
}

const page = `
{{define "MONSTER"}}
<div style="width: 50px; border-bottom: 1px solid black;"/>
{{.Name}} {{.}}
{{end}}
{{define "STATBLOCK"}}
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
{{end}}
<!DOCTYPE html>
<html>
<head>
 <link href="https://fonts.googleapis.com/css?family=Libre+Baskerville:700" rel="stylesheet" type="text/css"/>
 <link href="http://fonts.googleapis.com/css?family=Noto+Sans:400,700,400italic,700italic" rel="stylesheet" type="text/css"/>
 <meta charset="utf-8"/>
 <title>{{.Name}}</title>
 <style>
      h1.encounter {
      		font-family: 'Libre Baskerville', 'Lora', 'Calisto MT',
                   'Bookman Old Style', Bookman, 'Goudy Old Style',
                   Garamond, 'Hoefler Text', 'Bitstream Charter',
                   Georgia, serif;
      		color: #7A200D;
		margin: 5pt;
		text-decoration: underline;
	}
	tr.header {
      		font-family: 'Noto Sans', 'Myriad Pro', Calibri, Helvetica, Arial,
                    sans-serif;
      		font-weight: bold;
		font-size: 12pt;
      		color: #7A200D;
	}
	tr.header td {
		padding: 2px;
	}
	tr.content {
      		font-family: 'Noto Sans', 'Myriad Pro', Calibri, Helvetica, Arial,
                    sans-serif;
		font-size: 8pt;
      		color: #7A200D;
	}
	tr.content td {
		vertical-align: bottom;
		padding-right: 4px;
	}
      body {
        margin: 0;
      }

      stat-block {
        /* A bit of margin for presentation purposes, to show off the drop
        shadow. */
        margin-left: 20px;
        margin-top: 20px;
	vertical-align: top;
      }
 </style>
</head>
<body>
<template id="tapered-rule">
  <style>
    svg {
      fill: #922610;
      /* Stroke is necessary for good antialiasing in Chrome. */
      stroke: #922610;
      margin-top: 0.6em;
      margin-bottom: 0.35em;
    }
  </style>
  <svg height="5" width="400">
    <polyline points="0,0 400,2.5 0,5"></polyline>
  </svg>
</template>
<script>
(function(window, document) {
  var elemName = 'tapered-rule';
  var thatDoc = document;
  var thisDoc = (thatDoc.currentScript || thatDoc._currentScript).ownerDocument;
  var proto = Object.create(HTMLElement.prototype, {
    createdCallback: {
      value: function() {
        var template = thisDoc.getElementById(elemName);
        var clone = thatDoc.importNode(template.content, true);
        this.createShadowRoot().appendChild(clone);
      }
    }
  });
  thatDoc.registerElement(elemName, {prototype: proto});
})(window, document);
</script>
<template id="top-stats">
  <style>
    ::content * {
      color: #7A200D;
    }
  </style>

  <tapered-rule></tapered-rule>
  <content></content>
  <tapered-rule></tapered-rule>
</template>
<script>
(function(window, document) {
  var elemName = 'top-stats';
  var thatDoc = document;
  var thisDoc = (thatDoc.currentScript || thatDoc._currentScript).ownerDocument;
  var proto = Object.create(HTMLElement.prototype, {
    createdCallback: {
      value: function() {
        var template = thisDoc.getElementById(elemName);
        var clone = thatDoc.importNode(template.content, true);
        this.createShadowRoot().appendChild(clone);
      }
    }
  });
  thatDoc.registerElement(elemName, {prototype: proto});
})(window, document);
</script>
<template id="creature-heading">
  <style>
    ::content > h1 {
      font-family: 'Libre Baskerville', 'Lora', 'Calisto MT',
                   'Bookman Old Style', Bookman, 'Goudy Old Style',
                   Garamond, 'Hoefler Text', 'Bitstream Charter',
                   Georgia, serif;
      color: #7A200D;
      font-weight: 700;
      margin: 0px;
      font-size: 23px;
      letter-spacing: 1px;
      font-variant: small-caps;
    }

    ::content > h2 {
      font-weight: normal;
      font-style: italic;
      font-size: 12px;
      margin: 0;
    }
  </style>
  <content select="h1"></content>
  <content select="h2"></content>
</template>
<script>
(function(window, document) {
  var elemName = 'creature-heading';
  var thatDoc = document;
  var thisDoc = (thatDoc.currentScript || thatDoc._currentScript).ownerDocument;
  var proto = Object.create(HTMLElement.prototype, {
    createdCallback: {
      value: function() {
        var template = thisDoc.getElementById(elemName);
        var clone = thatDoc.importNode(template.content, true);
        this.createShadowRoot().appendChild(clone);
      }
    }
  });
  thatDoc.registerElement(elemName, {prototype: proto});
})(window, document);
</script>
<template id="abilities-block">
  <style>
    table {
      width: 100%;
      border: 0px;
      border-collapse: collapse;
    }
    th, td {
      width: 50px;
      text-align: center;
    }
  </style>
  <tapered-rule></tapered-rule>
  <table>
   <tbody>
    <tr>
      <th>STR</th>
      <th>DEX</th>
      <th>CON</th>
      <th>INT</th>
      <th>WIS</th>
      <th>CHA</th>
    </tr>
    <tr>
      <td id="str"></td>
      <td id="dex"></td>
      <td id="con"></td>
      <td id="int"></td>
      <td id="wis"></td>
      <td id="cha"></td>
    </tr>
   </tbody>
  </table>
  <tapered-rule></tapered-rule>
</template><script>
(function(window, document) {
  function abilityModifier(abilityScore) {
    var score = parseInt(abilityScore, 10);
    return Math.floor((score - 10) / 2);
  }

  function formattedModifier(abilityModifier) {
    if (abilityModifier >= 0) {
      return '+' + abilityModifier;
    }
    // This is an en dash, NOT a "normal" dash. The minus sign needs to be more
    // visible.
    return '–' + Math.abs(abilityModifier);
  }

  function abilityText(abilityScore) {
    return [String(abilityScore),
            ' (',
            formattedModifier(abilityModifier(abilityScore)),
            ')'].join('');
  }

  var elemName = 'abilities-block';
  var thatDoc = document;
  var thisDoc = (thatDoc.currentScript || thatDoc._currentScript).ownerDocument;
  var proto = Object.create(HTMLElement.prototype, {
    createdCallback: {
      value: function() {
        var template = thisDoc.getElementById(elemName);
        var clone = thatDoc.importNode(template.content, true);
        var root = this.createShadowRoot().appendChild(clone);
      }
    },
    attachedCallback: {
      value: function() {
        var root = this.shadowRoot;
        for (var i = 0; i < this.attributes.length; i++) {
          var attribute = this.attributes[i];
          var abilityShortName = attribute.name.split('-')[1];
          root.getElementById(abilityShortName).textContent =
             abilityText(attribute.value);
        }

      }
    }
  });
  thatDoc.registerElement(elemName, {prototype: proto});
})(window, document);
</script><template id="property-block">
  <style>
    :host {
      margin-top: 0.3em;
      margin-bottom: 0.9em;
      line-height: 1.5;
      display: block;
    }

    ::content > h4 {
      margin: 0;
      display: inline;
      font-weight: bold;
      font-style: italic;
    }

    ::content > p:first-of-type {
      display: inline;
      text-indent: 0;
    }

    ::content > p {
      text-indent: 1em;
      margin: 0;
    }
  </style>
  <content></content>
</template><script>
(function(window, document) {
  var elemName = 'property-block';
  var thatDoc = document;
  var thisDoc = (thatDoc.currentScript || thatDoc._currentScript).ownerDocument;
  var proto = Object.create(HTMLElement.prototype, {
    createdCallback: {
      value: function() {
        var template = thisDoc.getElementById(elemName);
        var clone = thatDoc.importNode(template.content, true);
        this.createShadowRoot().appendChild(clone);
      }
    }
  });
  thatDoc.registerElement(elemName, {prototype: proto});
})(window, document);
</script><template id="property-line">
  <style>
    :host {
      line-height: 1.4;
      display: block;
      text-indent: -1em;
      padding-left: 1em;
    }

    ::content > h4 {
      margin: 0;
      display: inline;
      font-weight: bold;
    }

    ::content > p:first-of-type {
      display: inline;
      text-indent: 0;
    }

    ::content > p {
      text-indent: 1em;
      margin: 0;
    }
  </style>
  <content></content>
</template><script>
(function(window, document) {
  var elemName = 'property-line';
  var thatDoc = document;
  var thisDoc = (thatDoc.currentScript || thatDoc._currentScript).ownerDocument;
  var proto = Object.create(HTMLElement.prototype, {
    createdCallback: {
      value: function() {
        var template = thisDoc.getElementById(elemName);
        var clone = thatDoc.importNode(template.content, true);
        this.createShadowRoot().appendChild(clone);
      }
    }
  });
  thatDoc.registerElement(elemName, {prototype: proto});
})(window, document);
</script><template id="stat-block">
  <style>
    .bar {
      height: 5px;
      background: #E69A28;
      border: 1px solid #000;
      position: relative;
      z-index: 1;
    }

    :host {
      display: inline-block;
    }

    #content-wrap {
      font-family: 'Noto Sans', 'Myriad Pro', Calibri, Helvetica, Arial,
                    sans-serif;
      font-size: 13.5px;
      background: #FDF1DC;
      padding: 0.6em;
      padding-bottom: 0.5em;
      border: 1px #DDD solid;
      box-shadow: 0 0 1.5em #867453;

      /* We don't want the box-shadow in front of the bar divs. */
      position: relative;
      z-index: 0;

      /* Leaving room for the two bars to protrude outwards */
      margin-left: 2px;
      margin-right: 2px;

      /* This is possibly overriden by next CSS rule. */
      width: 400px;

      -webkit-columns: 400px;
         -moz-columns: 400px;
              columns: 400px;
      -webkit-column-gap: 40px;
         -moz-column-gap: 40px;
              column-gap: 40px;

      /* When height is constrained, we want sequential filling of columns. */
      -webkit-column-fill: auto;
         -moz-column-fill: auto;
              column-fill: auto;
    }

    :host([data-two-column]) #content-wrap {
      /* One column is 400px and the gap between them is 40px. */
      width: 840px;
    }

    ::content > h3 {
      border-bottom: 1px solid #7A200D;
      color: #7A200D;
      font-size: 21px;
      font-variant: small-caps;
      font-weight: normal;
      letter-spacing: 1px;
      margin: 0;
      margin-bottom: 0.3em;

      break-inside: avoid-column;
      break-after: avoid-column;
    }

    /* For user-level p elems. */
    ::content > p {
      margin-top: 0.3em;
      margin-bottom: 0.9em;
      line-height: 1.5;
    }

    /* Last child shouldn't have bottom margin, too much white space. */
    ::content > *:last-child {
      margin-bottom: 0;
    }
  </style>
  <div class="bar"></div>
  <div id="content-wrap">
    <content></content>
  </div>
  <div class="bar"></div>
</template><script>
(function(window, document) {
  var elemName = 'stat-block';
  var thatDoc = document;
  var thisDoc = (thatDoc.currentScript || thatDoc._currentScript).ownerDocument;
  var proto = Object.create(HTMLElement.prototype, {
    createdCallback: {
      value: function() {
        var template = thisDoc.getElementById(elemName);
        // If the attr() CSS3 function were properly implemented, we wouldn't
        // need this hack...
        if (this.hasAttribute('data-content-height')) {
          var wrap = template.content.getElementById('content-wrap');
          wrap.style.height = this.getAttribute('data-content-height') + 'px';
        }
        var clone = thatDoc.importNode(template.content, true);
        this.createShadowRoot().appendChild(clone);
      }
    }
  });
  thatDoc.registerElement(elemName, {prototype: proto});
})(window, document);
</script>

<h1 class="encounter">{{.Name}}</h1>
<table>
<tr>
<td colspan="2">
<table>
<tr>
<td style="vertical-align: top">
<div style="margin-left: 1em; border: 1px solid black; padding: 4px">
 <table>
  <tr class="header"><td>Initiative</td></tr>
{{range intarray 22 4}}
  <tr class="content"><td width="120px" style="border-bottom: 1px solid black; margin-right: 1em">{{.}}</td></tr>
{{end}}
 </table>
</div>
</td>
<td style="vertical-align: top">
<div style="margin-left: 1em; border: 1px solid black; padding: 4px">
 <table>
  <tr class="header">
   <td>Monster</td>
   <td>AC</td>
   <td>Conditions</td>
   <td>Current HP</td>
  </tr>
{{range $i, $m := .Monsters}}
{{range $num := intarray 0 $m.Quantity}}
  <tr class="content">
   <td>{{$m.Monster.Name}} {{add $num 1}}</td><td>{{$m.Monster.ShortAc}}</td>
   <td style="width: 40px; border-bottom: 1px solid black"/>
   <td style="width: 300px; border-bottom: 1px solid black">{{$m.Monster.ShortHp}}</td>
  </tr>
{{end}}
{{end}}
 </table>
</div>
</td></tr></table>
</td>
</tr>
<tr>
{{range $i, $m := .Monsters}}
<td valign="top">{{template "STATBLOCK" $m.Monster}}</td>
{{if breakrow $i 2}}</tr><tr>{{end}}
{{end}}
</tr>
</table>
</body></html>
`
