package main

import (
	"encoding/csv"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/olekukonko/tablewriter"
)

const (
	_additions = "Additions"
	_deletions = "Deletions"
	_commits   = "Commits"
	_files     = "Files"
)

var (
	_headers = []string{"Contributor", "Commits", "Additions", "Deletions", "Files"}
)

// Stats represents git activity stats of a contributor.
type Stats struct {
	NameEmail string
	Counts    map[string]string
}

func main() {
	outputType := flag.String("output", "table", "output type: csv/json/table (default is table)")

	flag.Usage = func() {
		fmt.Printf("Print the stats of all contributors of a git repository.\n")
		fmt.Printf("The command must be run in a git repository.\n\n")
		fmt.Printf("Usage:\n\n\t%v\n\n", "gitstats [options]")
		fmt.Printf("The options are:\n\n")
		flag.VisitAll(func(f *flag.Flag) {
			fmt.Printf("\t-%v\n", f.Name)
			fmt.Printf("\t\t%v\n", f.Usage)
		})
	}

	flag.Parse()

	stats := findCommits()
	stats = findContributorStats(stats)

	switch *outputType {
	case "table":
		printTable(_headers, stats)
	case "csv":
		printCSV(_headers, stats)
	case "json":
		printJSON(_headers, stats)
	default:
		printTable(_headers, stats)
	}
}

func printCSV(headers []string, stats []Stats) {
	w := csv.NewWriter(os.Stdout)
	w.Write(headers)
	for _, v := range stats {
		w.Write([]string{v.NameEmail, v.Counts[_commits], v.Counts[_additions], v.Counts[_deletions], v.Counts[_files]})
	}
	w.Flush()
}

func printJSON(headers []string, stats []Stats) {
	b, _ := json.Marshal(stats)
	fmt.Println(string(b))
}

func printTable(headers []string, stats []Stats) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader(headers)
	table.SetRowLine(true)
	table.SetRowSeparator("-")
	table.SetColumnSeparator("|")

	rows := make([][]string, 0)
	for _, v := range stats {
		row := make([]string, 0)
		row = append(row, v.NameEmail, v.Counts[_commits], v.Counts[_additions], v.Counts[_deletions], v.Counts[_files])
		rows = append(rows, row)
	}

	table.AppendBulk(rows)
	table.Render()
}

func findCommits() []Stats {
	args := []string{"--no-pager", "shortlog", "-sne", "--all"}
	cmd := exec.Command("git", args...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println("unable to get commit details: git", args, ":", err)
		return nil
	}

	stats := make([]Stats, 0)
	lines := strings.Split(string(out), "\n")
	// line format: <number of commits> <Name email>
	// e.g.: 20 Peter Quill <starlord@gotg.space>
	for _, line := range lines {
		x := Stats{Counts: make(map[string]string)}
		line = strings.TrimSpace(line)
		if line == "" {
			break
		}
		// number of commits
		i := 0
		for {
			if !(line[i] >= '0' && line[i] <= '9') {
				break
			}
			x.Counts[_commits] += string(line[i])
			i++
		}

		// name email
		x.NameEmail = strings.TrimSpace(line[i:])
		stats = append(stats, x)
	}

	return stats
}

func findContributorStats(stats []Stats) []Stats {
	argTmpl := "git --no-pager log --author=\"AUTHOR\" --format=tformat: --numstat --all"

	for i := 0; i < len(stats); i++ {
		arg := strings.Replace(argTmpl, "AUTHOR", stats[i].NameEmail, 1)
		cmd := exec.Command("sh", "-c", arg)
		out, err := cmd.Output()
		if err != nil {
			fmt.Println("unable to get contributor stats: git", cmd.Args, ":", err)
			return nil
		}

		lines := strings.Split(strings.TrimSpace(string(out)), "\n")
		// record fileNames to count unique file names
		fileNames := make(map[string]struct{})
		additions, deletions := 0, 0

		// line format: <additions> <deletions> <filename>
		// e.g.: 18 3 README.md

		// sum up for all files
		for _, line := range lines {
			line = strings.TrimSpace(line)
			if line == "" {
				continue
			}
			// addition
			j, s := 0, ""
			for {
				if !(line[j] >= '0' && line[j] <= '9') {
					break
				}
				s += string(line[j])
				j++
			}
			line = strings.TrimSpace(line[j:])
			n, _ := strconv.Atoi(s)
			additions += n

			// deletion
			j, s = 0, ""
			for {
				if !(line[j] >= '0' && line[j] <= '9') {
					break
				}
				s += string(line[j])
				j++
			}
			line = strings.TrimSpace(line[j:])
			n, _ = strconv.Atoi(s)
			deletions += n

			// filename
			fileNames[line] = struct{}{}
		}

		stats[i].Counts[_additions] = strconv.Itoa(additions)
		stats[i].Counts[_deletions] = strconv.Itoa(deletions)
		stats[i].Counts[_files] = strconv.Itoa(len(fileNames))
	}

	return stats
}
