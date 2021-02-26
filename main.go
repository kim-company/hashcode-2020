package main

import (
	"os"
	"fmt"
	"log"
	"io"
	"sort"
)

func logf(format string, args ...interface{}) {
	log.Printf(format, args...)
}

func fatale(err error) {
	logf("Error: %v", err)
	os.Exit(1)
}

func must(n int, err error) {
	if n == 0 {
		fatale(fmt.Errorf("empty scan"))
	}
	if err != nil {
		fatale(err)
	}
}

type Car struct {
	Roadmap []string
	Tripd int // Trip duration
}

type Intersection struct {
	Idx int
	In map[string]*Street
}

func (i *Intersection) Incoming() []*Street {
	streets := make([]*Street, 0, len(i.In))
	for _, v := range i.In {
		streets = append(streets, v)
	}
	sort.SliceStable(streets, func(i, j int) bool {
		return streets[i].HottestAt() < streets[j].HottestAt()
	})
	return streets
}

type Street struct {
	Name string
	T int // Time the semaphore is green on this street
	N int // Number of cars passing trough this point potentially
	S *Intersection
	E *Intersection
	L int // Length of the street (t required to cross it)
	First int // When the first car hits this point ever.
	Hotness []int

	score float64
}

func (s *Street) HottestAt() int {
	var max, maxi, rise, risei int
	for i, v := range s.Hotness {
		if v > max {
			maxi, max = i, v
		}
	}
	rise, risei = max, maxi
	for i := maxi; i > 0; i-- {
		if s.Hotness[i] > rise {
			break
		}
		risei--
	}
	return risei
}

func Encode(w io.Writer, ii []*Intersection) {
	fmt.Fprintf(w, "%d\n", len(ii))
	for _, i := range ii {
		fmt.Fprintf(w, "%d\n", i.Idx)
		fmt.Fprintf(w, "%d\n", len(i.In))
		for _, s := range i.Incoming() {
			fmt.Fprintf(w, "%s %d\n", s.Name, s.T)
		}
	}
}

func main() {
	log.SetFlags(0)
	filename := os.Args[1]
	in, err := os.Open(filename)
	if err != nil {
		fatale(err)
	}
	defer in.Close()

	var (
		d int
		ni int
		ns int
		nc int
		bonus int
	)

	must(fmt.Fscanln(in, &d, &ni, &ns, &nc, &bonus))

	inter := make([]*Intersection, ni)
	for i := 0; i < ni; i++ {
		inter[i] = &Intersection{
			Idx: i,
			In: make(map[string]*Street),
		}
	}

	streets := make(map[string]*Street)
	for i := 0; i < ns; i++ {
		var name string
		var start, end, l int
		must(fmt.Fscanln(in, &start, &end, &name, &l))
		streets[name] = &Street{
			Name: name,
			S: inter[start],
			E: inter[end],
			L: l,
			Hotness: make([]int, d),
		}
		inter[end].In[name] = streets[name]
	}

	cars := make([]*Car, nc)
	for i := 0; i < nc; i++ {
		var n int
		must(fmt.Fscan(in, &n))
		c := Car{}
		for j := 0; j < n; j++ {
			var street string
			must(fmt.Fscan(in, &street))
			c.Roadmap = append(c.Roadmap, street)
		}
		cars[i] = &c
	}

	logf("T: %d", d)
	logf("Cars: %d", nc)
	logf("Streets: %d", ns)
	logf("Intersections: %d", ni)

	for _, c := range cars {
		t := 0
		for _, n := range c.Roadmap {
			s := streets[n]
			s.N++
			for i := t; i < t+s.L; i++ {
				if i >= len(s.Hotness) {
					break
				}
				s.Hotness[i]++
			}
			t += s.L
		}
	}
	for _, i := range inter {
		remove := make([]string, 0, len(i.In))
		for _, s := range i.In {
			if s.N <= 0 {
				remove = append(remove, s.Name)
			}
		}
		for _, v := range remove {
			delete(i.In, v)
		}

		var tot float64
		for _, s := range i.In {
			tot += float64(s.N)
		}
		var min float64
		for _, s := range i.In {
			s.score = float64(s.N)/tot
			if s.score < min || min == 0 {
				min = s.score
			}
		}
		for _, s := range i.In {
			s.T = int(s.score / min )
		}
	}
	valid := make([]*Intersection, 0, len(inter))
	for _, i := range inter {
		if len(i.In) > 0 {
			valid = append(valid, i)
		}
	}

	Encode(os.Stdout, valid)
}
