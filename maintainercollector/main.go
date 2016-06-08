package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/Sirupsen/logrus"
)

const (
	ghRawUri = "https://raw.githubusercontent.com"
	head     = `#
# THIS FILE IS AUTOGENERATED; SEE "./maintainercollector"!
#
# Docker projects maintainers file
#
# This file describes who runs the Docker project and how.
# This is a living document - if you see something out of date or missing,
# speak up!
#
# It is structured to be consumable by both humans and programs.
# To extract its contents programmatically, use any TOML-compliant
# parser.
`
)

var (
	org      = "docker"
	projects = []string{
		"boot2docker",
		"compose",
		"containerd",
		"datakit",
		"distribution",
		"docker",
		"docker-bench-security",
		"docker-credential-helpers",
		"docker-py",
		"dockercraft",
		"engine-api",
		"go-connections",
		"go-events",
		"go-p9p",
		"go-plugins-helpers",
		"go-units",
		"hyperkit",
		"kitematic",
		"leadership",
		"leeroy",
		"libchan",
		"libcompose",
		"libkv",
		"libnetwork",
		"machine",
		"migrator",
		"notary",
		"spdystream",
		"swarm",
		"swarmkit",
		"swarm-frontends",
		"toolbox",
		"v1.10-migrator",
		"vpnkit",
	}
)

//go:generate go run generate.go

func main() {
	// initialize the project MAINTAINERS file
	projectMaintainers := Maintainers{
		Org:    map[string]*Org{},
		People: map[string]Person{},
	}

	// initialize Curators
	projectMaintainers.Org["Curators"] = &Org{}
	projectMaintainers.Org["Docs maintainers"] = &Org{}

	// parse the MAINTAINERS file for each repo
	for _, project := range projects {
		maintainers, err := getMaintainers(project)
		if err != nil {
			logrus.Errorf("%s: parsing MAINTAINERS file failed: %v", project, err)
			continue
		}

		// create the Org object for the project
		p := &Org{
			// Repo: fmt.Sprintf("https://github.com/%s/%s", org, project),
			// TODO: change this to:
			// People: maintainers.Org["Core maintainers"].People,
			// once MaintainersDepreciated is removed.
			People: maintainers.Organization.CoreMaintainers.People,
		}

		// lowercase all maintainers nicks for consistency
		for i, n := range p.People {
			p.People[i] = strings.ToLower(n)
		}

		projectMaintainers.Org[project] = p

		if maintainers.Organization.DocsMaintainers != nil {
			projectMaintainers.Org["Docs maintainers"].People = append(projectMaintainers.Org["Docs maintainers"].People, maintainers.Organization.DocsMaintainers.People...)
		}

		if maintainers.Organization.Curators != nil {
			projectMaintainers.Org["Curators"].People = append(projectMaintainers.Org["Curators"].People, maintainers.Organization.Curators.People...)
		}

		// iterate through the people and add them to compiled list
		for nick, person := range maintainers.People {
			projectMaintainers.People[strings.ToLower(nick)] = person
		}
	}

	projectMaintainers.Org["Curators"].People = removeDuplicates(projectMaintainers.Org["Curators"].People)
	projectMaintainers.Org["Docs maintainers"].People = removeDuplicates(projectMaintainers.Org["Docs maintainers"].People)

	// encode the result to a file
	buf := new(bytes.Buffer)
	t := toml.NewEncoder(buf)
	t.Indent = "    "
	if err := t.Encode(projectMaintainers); err != nil {
		logrus.Fatalf("TOML encoding error: %v", err)
	}

	file := append([]byte(head), []byte(rules)...)
	file = append(file, []byte(roles)...)
	file = append(file, buf.Bytes()...)

	if err := ioutil.WriteFile("MAINTAINERS", file, 0755); err != nil {
		logrus.Fatal(err)
	}

	logrus.Infof("Successfully wrote new combined MAINTAINERS file.")
}

func removeDuplicates(slice []string) []string {
	seens := map[string]bool{}
	uniqs := []string{}
	for _, element := range slice {
		if _, seen := seens[element]; !seen {
			uniqs = append(uniqs, element)
			seens[element] = true
		}
	}
	return uniqs
}

func getMaintainers(project string) (maintainers MaintainersDepreciated, err error) {
	fileUrl := fmt.Sprintf("%s/%s/%s/master/MAINTAINERS", ghRawUri, org, project)
	resp, err := http.Get(fileUrl)
	if err != nil {
		return maintainers, fmt.Errorf("%s: %v", project, err)
	}
	defer resp.Body.Close()

	file, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return maintainers, fmt.Errorf("%s: %v", project, err)
	}

	if _, err := toml.Decode(string(file), &maintainers); err != nil {
		return maintainers, fmt.Errorf("%s: parsing MAINTAINERS file failed: %v", project, err)
	}

	return maintainers, nil
}
