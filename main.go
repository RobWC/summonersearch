package main

import (
	"html/template"
	"os"
	"time"

	"github.com/robwc/riotapi"

	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

const (
	summonerTmpl = `Name: {{ .Name }}
Level: {{ .Level }}
Last Seen: {{ lastSeen .RevisionDate }}
`
)

var (
	server = kingpin.Flag("server", "Specify server").Default("na").String()

	search       = kingpin.Command("search", "Search for summoner by name")
	summonerName = search.Arg("name", "Name to search for").Required().String()
)

func lastSeen(last int) string {
	return time.Unix(int64(last/1000), 0).Format("Mon Jan 2 15:04:05 -0700 MST 2006")
}

func main() {
	switch kingpin.Parse() {
	case "search":
		apiKeyEnv := os.Getenv("RIOTKEY")

		if apiKeyEnv == "" {
			panic("API Key Not Specified")
		}

		c := riotapi.NewAPIClient(*server, apiKeyEnv)
		s, err := c.SummonerByName(*summonerName)
		if err != nil && s != nil {
			panic(err)
		}

		t := template.Must(template.New("champ").Funcs(template.FuncMap{"lastSeen": lastSeen}).Parse(summonerTmpl))

		err = t.Execute(os.Stdout, s)
		if err != nil {
			panic(err)
		}

	}

	os.Exit(0)
}
