package main

import (
	"encoding/csv"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"sort"
	"strconv"
	"strings"

	"github.com/olekukonko/tablewriter"
)

const (
	_additions = "Additions"
	_deletions = "Deletions"
	_commits   = "Commits"
	_files     = "Files"

	_outputUsage  = "output type: csv/json/table (default is table)"
	_mergeIDUsage = "merge contributor stats with same email username"
	_sortByUsage  = `sort by: commits,additions,deletions,files (default is commits)
		accepts multiple values separated by commas, sort is stable
		sort is applied in the given order left to right`
)

var (
	_headers = []string{"Contributor", "Commits", "Additions", "Deletions", "Files"}
)

// Stats represents git activity stats of a contributor.
type Stats struct {
	NameEmail string
	FileNames map[string]struct{}
	Counts    map[string]string
}

func main() {
	var outputType string
	flag.StringVar(&outputType, "output", "table", _outputUsage)
	flag.StringVar(&outputType, "o", "table", _outputUsage)

	mergeNameFlag := flag.Bool("merge-name", false, _mergeIDUsage)

	var sortBy string
	flag.StringVar(&sortBy, "sort-by", "commits", _sortByUsage)
	flag.StringVar(&sortBy, "s", "commits", _sortByUsage)

	flag.Usage = func() {
		fmt.Printf("Print the stats of all contributors of a git repository.\n")
		fmt.Printf("The command must be run in a git repository.\n\n")
		fmt.Printf("Usage:\n\n\t%v\n\n", "gitstats [options]")
		fmt.Println(`The options are:
	-output, -o
		` + _outputUsage + `
	-sort-by, -s
		` + _sortByUsage + `

The flags are:
	-merge-name
		` + _mergeIDUsage)
	}

	flag.Parse()

	if isGitRepo() == false {
		fmt.Println("gitstats must be run in a git repository. Type gitstats -h for help.")
		return
	}

	stats := findCommits()
	stats = findContributorStats(stats)

	if *mergeNameFlag {
		stats = mergeNames(stats)
	}

	stats = sortStats(sortBy, stats)

	switch outputType {
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

func isGitRepo() bool {
	args := []string{"status"}
	cmd := exec.Command("git", args...)
	_, err := cmd.CombinedOutput()
	if err != nil {
		return false
	}

	return true
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
		x := Stats{FileNames: make(map[string]struct{}, 0), Counts: make(map[string]string)}
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

		stats[i].FileNames = fileNames
		stats[i].Counts[_additions] = strconv.Itoa(additions)
		stats[i].Counts[_deletions] = strconv.Itoa(deletions)
		stats[i].Counts[_files] = strconv.Itoa(len(fileNames))
	}

	return stats
}

func mergeNames(stats []Stats) []Stats {
	m := make(map[string][]Stats)

	for _, v := range stats {
		userName := findUserName(v.NameEmail)
		if _, ok := m[userName]; !ok {
			m[userName] = make([]Stats, 0)
		}

		m[userName] = append(m[userName], v)
	}

	mergedStats := make([]Stats, 0)
	for _, v := range m {
		x := Stats{FileNames: make(map[string]struct{}), Counts: make(map[string]string)}

		// commits, additions, deletions
		c, a, d := 0, 0, 0
		// filenames
		fm := make(map[string]struct{})
		nameEmails := make([]string, 0)
		for _, w := range v {
			i, _ := strconv.Atoi(w.Counts[_commits])
			c += i
			i, _ = strconv.Atoi(w.Counts[_additions])
			a += i
			i, _ = strconv.Atoi(w.Counts[_deletions])
			d += i

			for name := range w.FileNames {
				fm[name] = struct{}{}
			}

			nameEmails = append(nameEmails, w.NameEmail)
		}

		x.NameEmail = strings.Join(nameEmails, "; ")
		x.Counts[_commits] = strconv.Itoa(c)
		x.Counts[_additions] = strconv.Itoa(a)
		x.Counts[_deletions] = strconv.Itoa(d)
		x.Counts[_files] = strconv.Itoa(len(fm))

		mergedStats = append(mergedStats, x)
	}

	return mergedStats
}

func findUserName(nameEmail string) string {
	parts := strings.Split(strings.TrimSpace(nameEmail), "<")
	parts = strings.Split(strings.TrimSpace(parts[1]), "@")
	return strings.ToLower(strings.TrimSpace(parts[0]))
}

func sortStats(sortBy string, stats []Stats) []Stats {
	sortFields := strings.Split(sortBy, ",")

	// last one wins
	for _, v := range sortFields {
		switch v {
		case "commits":
			sort.SliceStable(stats,
				func(i, j int) bool {
					a, _ := strconv.Atoi(stats[i].Counts[_commits])
					b, _ := strconv.Atoi(stats[j].Counts[_commits])
					return a > b
				},
			)
		case "additions":
			sort.SliceStable(stats,
				func(i, j int) bool {
					a, _ := strconv.Atoi(stats[i].Counts[_additions])
					b, _ := strconv.Atoi(stats[j].Counts[_additions])
					return a > b
				},
			)
		case "deletions":
			sort.SliceStable(stats,
				func(i, j int) bool {
					a, _ := strconv.Atoi(stats[i].Counts[_deletions])
					b, _ := strconv.Atoi(stats[j].Counts[_deletions])
					return a > b
				},
			)
		case "files":
			sort.SliceStable(stats,
				func(i, j int) bool {
					a, _ := strconv.Atoi(stats[i].Counts[_files])
					b, _ := strconv.Atoi(stats[j].Counts[_files])
					return a > b
				},
			)
		}
	}

	return stats
}
