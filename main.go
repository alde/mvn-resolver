package main

import (
	"bufio"
	"encoding/xml"
	"fmt"
	"log"
	"os"
	"regexp"
	"sort"
	"strings"

	"github.com/hashicorp/go-version"
)

type dependency struct {
	XMLName    xml.Name `xml:"dependency"`
	GroupID    string   `xml:"groupId"`
	ArtifactID string   `xml:"artifactId"`
	Version    string   `xml:"version"`
}

func (d *dependency) String() string {
	output, _ := xml.MarshalIndent(d, "", "    ")
	return fmt.Sprint(string(output))
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	deps := scan(scanner)
	if len(deps) == 0 {
		fmt.Print("\nno dependency violations found\n")
		os.Exit(0)
	}
	fmt.Print("\nsmallest common versions:\n\n")
	for _, dep := range deps {
		fmt.Println(dep)
	}
}

func scan(scanner *bufio.Scanner) (deps []*dependency) {
	matcher := regexp.MustCompile("Require upper bound dependencies error for (.*) paths to dependency are:")
	for scanner.Scan() {
		line := scanner.Text()
		l := matcher.FindStringSubmatch(line)
		if len(l) == 0 {
			continue
		}
		upperBounds := parse(l[1])
		deps = append(deps, maxForDep(scanner, upperBounds))
	}

	if err := scanner.Err(); err != nil {
		log.Printf("error reading from stdin: %s\n", err)
	}
	return
}

func maxForDep(scanner *bufio.Scanner, dep *dependency) *dependency {
	candidates := []*dependency{}
	matcher := regexp.MustCompile(fmt.Sprintf("%s:%s:([^\\s]+)", dep.GroupID, dep.ArtifactID))
	for scanner.Scan() {
		line := scanner.Text()
		if line == "," {
			break
		}
		l := matcher.FindAllStringSubmatch(line, -1)
		if l == nil {
			continue
		}
		for _, version := range l {
			candidates = append(candidates, &dependency{
				GroupID:    dep.GroupID,
				ArtifactID: dep.ArtifactID,
				Version:    version[1],
			})
		}
	}
	sort.Slice(candidates, func(i, j int) bool {
		v1, _ := version.NewVersion(candidates[i].Version)
		v2, _ := version.NewVersion(candidates[j].Version)
		return v1.GreaterThan(v2)
	})
	return candidates[0]
}

func parse(line string) *dependency {
	parts := strings.Split(line, ":")
	dep := &dependency{
		GroupID:    parts[0],
		ArtifactID: parts[1],
		Version:    parts[2],
	}
	return dep
}
