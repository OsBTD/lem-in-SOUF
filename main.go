package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

type farm struct {
	ants_number int
	rooms       map[string][]int
	start       map[string][]int
	end         map[string][]int
	links       map[string][]string
}

func main() {
	var myFarm farm
	myFarm.Read("test.txt")
	fmt.Printf("\nall sorted paths from start to end: %v\n", BFS(myFarm))
	fmt.Println("Place all Ants on there path: ", Ants(myFarm, BFS(myFarm)))
	MoveAnts(Ants(myFarm, BFS(myFarm)))

	// fmt.Println(Ants(myFarm, BFS(myFarm)))
	// fmt.Println("number of ants is : ", myFarm.ants_number)
	// fmt.Println("rooms are : ", myFarm.rooms)
	// fmt.Println("start is : ", myFarm.start)
	// fmt.Println("end is : ", myFarm.end)
	// fmt.Println("links are : ", myFarm.links)
	// fmt.Println("adjacent is : ", Graph(myFarm))
}

func (myFarm *farm) Read(filename string) {
	bytes, err := os.ReadFile(filename)
	if err != nil {
		log.Println("error reading", err)
	}
	content := strings.Split(string(bytes), "\n")

	myFarm.rooms = make(map[string][]int)
	myFarm.start = make(map[string][]int)
	myFarm.end = make(map[string][]int)
	myFarm.links = make(map[string][]string)

	var st, en int
	number, err := strconv.Atoi(content[0])
	if err != nil {
		log.Println("couldn't convert", err)
	}
	myFarm.ants_number = number

	for index := range content {
		if strings.TrimSpace(content[index]) == "##start" {
			st++
			if index+1 <= len(content)-1 {
				split := strings.Split(strings.TrimSpace(content[index+1]), " ")
				x, err := strconv.Atoi(split[1])
				y, err2 := strconv.Atoi(split[2])
				if err == nil && err2 == nil {
					myFarm.start[split[0]] = []int{x, y}
				}

			}

		} else if strings.TrimSpace(content[index]) == "##end" {
			en++
			if index+1 <= len(content)-1 {
				split := strings.Split(strings.TrimSpace(content[index+1]), " ")
				x, err := strconv.Atoi(split[1])
				y, err2 := strconv.Atoi(split[2])
				if err == nil && err2 == nil {
					myFarm.end[split[0]] = []int{x, y}
				}

			}
		} else if strings.Contains(content[index], "-") {
			split := strings.Split(strings.TrimSpace(content[index]), "-")
			if len(split) == 2 {
				myFarm.links[split[0]] = append(myFarm.links[split[0]], split[1])
			}
		} else if strings.Count(content[index], " ") == 2 {
			split := strings.Split(strings.TrimSpace(content[index]), " ")
			if len(split) == 3 {
				x, err := strconv.Atoi(split[1])
				y, err2 := strconv.Atoi(split[2])
				if err == nil || err2 == nil {
					myFarm.rooms[split[0]] = []int{x, y}
				}
			}
		} else if (strings.HasPrefix(strings.TrimSpace(content[index]), "#") || strings.HasPrefix(strings.TrimSpace(content[index]), "L")) && (strings.TrimSpace(content[index]) != "##start" && strings.TrimSpace(content[index]) != "##end") {
			continue
		}
	}
	if en != 1 || st != 1 {
		log.Println("rooms setup is incorrect", err)
	}
}

func Graph(farm farm) map[string][]string {
	adjacent := make(map[string][]string)
	for room := range farm.rooms {
		adjacent[room] = []string{}
	}
	for room, links := range farm.links {
		for _, link := range links {
			adjacent[room] = append(adjacent[room], link)
			adjacent[link] = append(adjacent[link], room)

		}
	}

	return adjacent
}

func BFS(myFarm farm) [][]string {
	adjacent := Graph(myFarm)
	var Queue []string
	var endd string
	start := myFarm.start
	end := myFarm.end
	var Sorted [][]string

	for key := range start {
		for _, adj := range adjacent[key] {
			Visited := make(map[string]bool)
			Parents := make(map[string]string)

			Queue = append(Queue, adj)
			Visited[adj] = true
			for key := range end {
				endd = key
			}

			for len(Queue) > 0 {
				current := Queue[0]
				Queue = Queue[1:]
				if current == endd {
					Queue = []string{}
					break
				}

				for _, link := range adjacent[current] {
					if !Visited[link] {
						Queue = append(Queue, link)
						Visited[link] = true
						Parents[link] = current
					}
				}
			}

			if !Visited[endd] {
				fmt.Printf("\n No path found to end room \n")
				return [][]string{}
			}

			path := []string{endd}
			current := endd

			for Parents[current] != "" {
				current = Parents[current]
				path = append([]string{current}, path...)
			}
			path = append([]string{key}, path...)
			Sorted = append(Sorted, path)
		}
	}
	Sorted = SortPath(Sorted)
	// fmt.Printf("\nall sorted paths from start to end: %v\n", Sorted)
	return Sorted
}

func SortPath(Paths [][]string) [][]string {
	if len(Paths) <= 1 {
		return Paths
	}
	pivot := Paths[len(Paths)-1]
	var less, greater [][]string
	for _, v := range Paths[:len(Paths)-1] {
		if len(v) <= len(pivot) {
			less = append(less, v)
		} else {
			greater = append(greater, v)
		}
	}
	return append(append(SortPath(less), pivot), SortPath(greater)...)
}

func Ants(myFarm farm, paths [][]string) [][]string {
	ants := myFarm.ants_number

	fmt.Println("num of ants is :", ants)

	k := 0
	for i := ants; i > 0; i-- {
		for j := 0; j < len(paths); j++ {
			if k < len(paths) {
				if len(paths[k]) >= len(paths[j]) {
					paths[k] = append(paths[k], "L"+strconv.Itoa(i))
					break
				}
			} else {
				k = 0
				if len(paths[k]) >= len(paths[j]) {
					paths[k] = append(paths[k], "L"+strconv.Itoa(i))
					break
				}
			}
		}
		k++
	}

	return paths
}

func MoveAnts(paths [][]string) {
	var p, a, text []string
	var lines [][]string
	var all [][][]string

	for i := 0; i < len(paths); i++ {
		for j := 1; j < len(paths[i]); j++ {
			if strings.HasPrefix(paths[i][j], "L") {
				a = append(a, paths[i][j])
			} else {
				p = append(p, paths[i][j])
			}
		}
		for l := 0; l < len(a); l++ {
			for x := 0; x < len(p); x++ {
				text = append(text, a[l]+"-"+p[x]+" ")
			}
			lines = append(lines, text)
			text = []string{}
		}
		fmt.Println("\nlines are : ", lines)
		var print [][]string
		for i := range lines {
			space := []string{}
			if i != 0 {
				for n := 0; n < i; n++ {
					space = append(space, "s")
				}
				lines[i] = append(space, lines[i]...)
				print = append(print, lines[i])
			} else {
				print = append(print, lines[i])
			}
			// fmt.Println("l0", i, lines[i][0])
			// fmt.Println("l1", i, lines[i][1])
			// fmt.Println("l2", i, line[2])
			// fmt.Println("l3", i, line[3])
		}

		all = append(all, print)
		a = []string{}
		p = []string{}
		lines = [][]string{}

	}
	fmt.Println("all is : ", all)
	for d1 := 0; d1 < len(all); d1++ { // 3 paths
		for d2 := 0; d2 < len(all[d1]); d2++ { // all ants in a path
			print(all[d1][d2][0])
		}
	}
}

// func split(s []string) []string {
// 	s = s[1:]
// 	return s
// }
