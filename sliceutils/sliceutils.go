package sliceutils

import "fmt"

// CompareList diffs the list lines and returns values only found in one list, intersect values
// and values, which are not unique in between one list
func CompareLists(l1, l2 []string) (only1 []string, only2 []string, inter []string, nu1 map[string]int, nu2 map[string]int) {
	nu1 = make(map[string]int)
	m1 := make(map[string]bool)

	for _, e := range l1 {
		_, ok := m1[e]
		if ok {
			_, found := nu1[e]
			if found {
				nu1[e]++
			} else {
				// it's 2 not 1: because we have found it once
				// and placed it to m1; for the second finding
				// we write it to nu!
				nu1[e] = 2
			}
		} else {
			m1[e] = true
		}
	}
	for k := range nu1 {
		delete(m1, k)
	}

	nu2 = make(map[string]int)
	m2 := make(map[string]bool)

	for _, e := range l2 {
		_, ok := m2[e]
		if ok {
			_, found := nu2[e]
			if found {
				nu2[e]++
			} else {
				// it's 2 not 1: because we have found it once
				// and placed it to m2; for the second finding
				// we write it to nu!
				nu2[e] = 2
			}
		} else {
			m2[e] = true
		}
	}
	for k := range nu2 {
		delete(m2, k)
	}

	// do the compare job...
	for k := range m1 {
		_, ok := m2[k]
		if ok {
			inter = append(inter, k)
			delete(m2, k)
		} else {
			only1 = append(only1, k)
		}
	}

	for k := range m2 {
		only2 = append(only2, k)
	}

	return only1, only2, inter, nu1, nu2
}

func MakeListFromMapkeys(m map[string]int) []string {
	lst := []string{}
	for k := range m {
		lst = append(lst, k)
	}
	return lst
}

func MakeListFromStringMapkeys(m map[string]string) []string {
	lst := []string{}
	for k := range m {
		lst = append(lst, k)
	}
	return lst
}

func TraverseSlice(rec [][]string) [][]string {
	fmt.Println(len(rec))
	fmt.Println(len(rec[0]))
	res := make([][]string, len(rec[0]))
	for i := 0; i < len(res); i++ {
		res[i] = make([]string, len(rec))
	}

	for i, row := range rec {
		for k, val := range row {
			res[k][i] = val
		}
	}
	return res
}
