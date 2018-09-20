package internal

type StatErrors struct {
	Page map[string]StatErrorPage //.Page["https://www.example.com/errors.html"].Missing["/404/missing.jog"] = 36
}

type StatErrorPage struct {
	Page    string            //used when sorting into []StatErrorPage instead of StatErrors.Page map
	Total   uint64            //.Page["https://www.example.com/errors.html"].Total = 329687296
	Missing map[string]uint64 //.Page["https://www.example.com/errors.html"].Missing["/404/missing.jog"] = 36
}

func (se *StatErrors) Increment(page, missing string) *StatErrors {
	if se.Page == nil {
		se.Page = map[string]StatErrorPage{}
	}

	if sep, ok := se.Page[page]; ok {
		sep.Missing[missing]++
		sep.Total++
		se.Page[page] = sep
	} else {
		//initialize
		sep := StatErrorPage{}
		sep.Missing = map[string]uint64{}

		sep.Missing[missing] = 1
		sep.Total++
		se.Page[page] = sep
	}

	return se
}

func (se *StatErrors) SetVal(page, missing string, val uint64) *StatErrors {
	if se.Page == nil {
		se.Page = map[string]StatErrorPage{}
	}

	if sep, ok := se.Page[page]; ok {
		sep.Missing[missing] = val
		sep.Total += val
		se.Page[page] = sep
	} else {
		//initialize
		sep := StatErrorPage{}
		sep.Missing = map[string]uint64{}

		sep.Missing[missing] = val
		sep.Total += val
		se.Page[page] = sep
	}

	return se
}

func (se StatErrors) ToSlice(min uint64) []StatErrorPage { //still used on ghetto.go
	out := make([]StatErrorPage, 0, len(se.Page))
	for k, v := range se.Page {
		if min > 0 && v.Total < min {
			continue
		}
		out = append(out, StatErrorPage{Page: k, Total: v.Total, Missing: v.Missing})
	}
	return out[:len(out):len(out)] // trim the slice to release the unused memory
}
