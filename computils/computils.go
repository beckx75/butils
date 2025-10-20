package computils

func CompareStringList(l1, l2 []string) (only1 []string, only2 []string, intersect []string, nuL1 map[string]int, nuL2 map[string]int) {
	nuL1 = make(map[string]int)
	nuL2 = make(map[string]int)
	m1 := make(map[string]bool)
	m2 := make(map[string]bool)

	for _, e := range l1 {
		_, ok := nuL1[e]
		if ok {
			nuL1[e]++
			continue
		}
		_, ok = m1[e]
		if ok {
			nuL1[e] = 2
			delete(m1, e)
		} else {
			m1[e] = true
		}
	}
	for _, e := range l2 {
		_, ok := nuL2[e]
		if ok {
			nuL2[e]++
			continue
		}
		_, ok = m2[e]
		if ok {
			nuL2[e] = 2
			delete(m2, e)
		} else {
			m2[e] = true
		}
	}

	o1, o2, inter := CompareMapKeys(m1, m2)
	return o1, o2, inter, nuL1, nuL2
}

func CompareMapKeys[V any](m1, m2 map[string]V) (only1 []string, only2 []string, intersect []string) {
	o1 := make(map[string]bool)
	o2 := make(map[string]bool)
	for k := range m2 {
		o2[k] = true
	}

	for k := range m1 {
		_, ok := m2[k]
		if ok {
			intersect = append(intersect, k)
			delete(o2, k)
		} else {
			o1[k] = true
		}
	}
	for k := range o1 {
		only1 = append(only1, k)
	}
	for k := range o2 {
		only2 = append(only2, k)
	}

	return only1, only2, intersect
}

func ChangeStringMapKeyVal(m1 map[string]string) (map[string]string, map[string][]string) {
	changed := make(map[string]string)
	nu := make(map[string][]string)

	for k, v := range m1 {
		_, found := nu[v]
		if found {
			nu[v] = append(nu[v], k)
		} else {
			nuKey, ok := changed[v]
			if ok {
				nu[v] = []string{nuKey, k}
			} else {
				changed[v] = k
			}
		}
	}
	return changed, nu
}
