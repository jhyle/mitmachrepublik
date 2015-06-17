package mmr

import (
	"github.com/armon/go-radix"
	"strings"
)

type (
	LocationTree struct {
		*radix.Tree
	}
)

const (
	MAX_POSTCODE_LENGTH = 5
)

var (
	DistrictMap map[string][]string = map[string][]string{
		"Berlin Charlottenburg-Wilmersdorf": []string{"Berlin Charlottenburg", "Berlin Charlottenburg-Nord", "Berlin Grunewald", "Berlin Halensee", "Berlin Schmargendorf", "Berlin Westend", "Berlin Wilmersdorf"},
		"Berlin Friedrichshain-Kreuzberg":   []string{"Berlin Friedrichshain", "Berlin Kreuzberg"},
		"Berlin Lichtenberg":                []string{"Berlin Alt-Hohenschönhausen", "Berlin Falkenberg", "Berlin Fennpfuhl", "Berlin Friedrichsfelde", "Berlin Karlshorst", "Berlin Lichtenberg", "Berlin Malchow", "Berlin Neu-Hohenschönhausen", "Berlin Rummelsburg", "Berlin Wartenberg"},
		"Berlin Marzahn-Hellersdorf":        []string{"Berlin Biesdorf", "Berlin Hellersdorf", "Berlin Kaulsdorf", "Berlin Kaulsdorf", "Berlin Mahlsdorf", "Berlin Marzahn"},
		"Berlin Mitte":                      []string{"Berlin Gesundbrunnen", "Berlin Hansaviertel", "Berlin Mitte", "Berlin Moabit", "Berlin Tiergarten", "Berlin Wedding"},
		"Berlin Neukölln":                   []string{"Berlin Britz", "Berlin Buckow", "Berlin Gropiusstadt", "Berlin Neukölln", "Berlin Rudow"},
		"Berlin Pankow":                     []string{"Berlin Blankenburg", "Berlin Blankenfelde", "Berlin Buch", "Berlin Französisch Buchholz", "Berlin Heinersdorf", "Berlin Karow", "Berlin Niederschönhausen", "Berlin Pankow", "Berlin Prenzlauer Berg", "Berlin Rosenthal", "Berlin Stadtrandsiedlung Malchow", "Berlin Weißensee", "Berlin Wilhelmsruh"},
		"Berlin Reinickendorf":              []string{"Berlin Borsigwalde", "Berlin Frohnau", "Berlin Heiligensee", "Berlin Hermsdorf", "Berlin Konradshöhe", "Berlin Lübars", "Berlin Märkisches Viertel", "Berlin Reinickendorf", "Berlin Tegel", "Berlin Waidmannslust", "Berlin Wittenau"},
		"Berlin Spandau":                    []string{"Berlin Falkenhagener Feld", "Berlin Gatow", "Berlin Hakenfelde", "Berlin Haselhorst", "Berlin Kladow", "Berlin Siemensstadt", "Berlin Spandau", "Berlin Staaken", "Berlin Wilhelmstadt"},
		"Berlin Steglitz-Zehlendorf":        []string{"Berlin Dahlem", "Berlin Lankwitz", "Berlin Lichterfelde", "Berlin Nikolassee", "Berlin Steglitz", "Berlin Wannsee", "Berlin Zehlendorf"},
		"Berlin Tempelhof-Schöneberg":       []string{"Berlin Friedenau", "Berlin Lichtenrade", "Berlin Mariendorf", "Berlin Marienfelde", "Berlin Schöneberg", "Berlin Tempelhof"},
		"Berlin Treptow-Köpenick":           []string{"Berlin Adlershof", "Berlin Alt-Treptow", "Berlin Altglienicke", "Berlin Baumschulenweg", "Berlin Bohnsdorf", "Berlin Friedrichshagen", "Berlin Grünau", "Berlin Johannisthal", "Berlin Köpenick", "Berlin Müggelheim", "Berlin Niederschöneweide", "Berlin Oberschöneweide", "Berlin Plänterwald", "Berlin Rahnsdorf", "Berlin Schmöckwitz"},
	}

	PostcodeMap map[string][]string = map[string][]string{
		// Charlottenburg-Wilmersdorf
		"Berlin Charlottenburg":      []string{"10553", "10585", "10587", "10589", "10623", "10625", "10627", "10629", "10707", "10709", "10711", "10719", "10787", "10789", "13627", "14050", "14055", "14057", "14059"},
		"Berlin Charlottenburg-Nord": []string{"10589", "13353", "13627"},
		"Berlin Grunewald":           []string{"10711", "14055", "14193", "14195", "14199"},
		"Berlin Halensee":            []string{"10709", "10711", "10713"},
		"Berlin Schmargendorf":       []string{"14193", "14195", "14197", "14199"},
		"Berlin Westend":             []string{"13597", "14050", "14052", "14053", "14055", "14057", "14059"},
		"Berlin Wilmersdorf":         []string{"10707", "10709", "10711", "10713", "10715", "10717", "10719", "10777", "10779", "10789", "10825", "14195", "14197", "14199"},
		// Friedrichshain-Kreuzberg
		"Berlin Friedrichshain": []string{"10178", "10179", "10243", "10245", "10247", "10249", "10317", "10407"},
		"Berlin Kreuzberg":      []string{"10785", "10961", "10963", "10965", "10967", "10969", "10997", "10999"},
		// Lichtenberg
		"Berlin Alt-Hohenschönhausen": []string{"12681", "13051", "13053", "13055"},
		"Berlin Falkenberg":           []string{"13057"},
		"Berlin Fennpfuhl":            []string{"10367", "10369"},
		"Berlin Friedrichsfelde":      []string{"10315", "10317", "10319"},
		"Berlin Karlshorst":           []string{"10317", "10318"},
		"Berlin Lichtenberg":          []string{"10315", "10317", "10365", "10367", "10369"},
		"Berlin Malchow":              []string{"13051"},
		"Berlin Neu-Hohenschönhausen": []string{"13051", "13053", "13057", "13059"},
		"Berlin Rummelsburg":          []string{"10315", "10317", "10318", "10365"},
		"Berlin Wartenberg":           []string{"13051", "13059"},
		// Mahrzahn-Hellersdorf
		"Berlin Biesdorf":    []string{"12683", "12685"},
		"Berlin Hellersdorf": []string{"12619", "12621", "12627", "12629"},
		"Berlin Kaulsdorf":   []string{"12555", "12619", "12621", "12623"},
		"Berlin Mahlsdorf":   []string{"12621", "12623"},
		"Berlin Marzahn":     []string{"12679", "12681", "12683", "12685", "12687", "12689"},
		// Mitte
		"Berlin Gesundbrunnen": []string{"13347", "13353", "13355", "13357", "13359", "13409"},
		"Berlin Hansaviertel":  []string{"10555", "10557"},
		"Berlin Mitte":         []string{"10115", "10117", "10119", "10178", "10179", "10435"},
		"Berlin Moabit":        []string{"10551", "10553", "10555", "10557", "10559", "13353"},
		"Berlin Tiergarten":    []string{"10117", "10557", "10623", "10785", "10787", "10963"},
		"Berlin Wedding":       []string{"13347", "13349", "13351", "13353", "13359", "13405", "13407"},
		// Neukölln
		"Berlin Britz":        []string{"12051", "12057", "12099", "12347", "12349", "12351", "12359"},
		"Berlin Buckow":       []string{"12107", "12305", "12349", "12351", "12353", "12357", "12359"},
		"Berlin Gropiusstadt": []string{"12351", "12353", "12357"},
		"Berlin Neukölln":     []string{"10965", "10967", "12043", "12045", "12047", "12049", "12051", "12053", "12055", "12057", "12059", "12099"},
		"Berlin Rudow":        []string{"12353", "12355", "12357", "12359"},
		// Pankow
		"Berlin Blankenburg":               []string{"13051", "13125", "13129"},
		"Berlin Blankenfelde":              []string{"13127", "13158", "13159"},
		"Berlin Buch":                      []string{"13125", "13127"},
		"Berlin Französisch Buchholz":      []string{"13127", "13129", "13156"},
		"Berlin Heinersdorf":               []string{"13088", "13089", "13129"},
		"Berlin Karow":                     []string{"13125"},
		"Berlin Niederschönhausen":         []string{"13127", "13156", "13158", "13187"},
		"Berlin Pankow":                    []string{"10439", "13129", "13187", "13189"},
		"Berlin Prenzlauer Berg":           []string{"10119", "10247", "10369", "10405", "10407", "10409", "10435", "10437", "10439", "13187", "13189"},
		"Berlin Rosenthal":                 []string{"13156", "13158"},
		"Berlin Stadtrandsiedlung Malchow": []string{"13051", "13088", "13089"},
		"Berlin Weißensee":                 []string{"13051", "13086", "13088"},
		"Berlin Wilhelmsruh":               []string{"13156", "13158"},
		//Reinickendorf
		"Berlin Borsigwalde":        []string{"13403", "13509"},
		"Berlin Frohnau":            []string{"13465"},
		"Berlin Heiligensee":        []string{"13503", "13505"},
		"Berlin Hermsdorf":          []string{"13467"},
		"Berlin Konradshöhe":        []string{"13505"},
		"Berlin Lübars":             []string{"13435", "13469"},
		"Berlin Märkisches Viertel": []string{"13435", "13439"},
		"Berlin Reinickendorf":      []string{"13403", "13405", "13407", "13409", "13437", "13509"},
		"Berlin Tegel":              []string{"13405", "13503", "13505", "13507", "13509", "13599", "13629"},
		"Berlin Waidmannslust":      []string{"13469"},
		"Berlin Wittenau":           []string{"13403", "13407", "13435", "13437", "13469"},
		// Spandau
		"Berlin Falkenhagener Feld": []string{"13583", "13585", "13589", "13591"},
		"Berlin Gatow":              []string{"14089"},
		"Berlin Hakenfelde":         []string{"13585", "13587", "13589"},
		"Berlin Haselhorst":         []string{"13597", "13599"},
		"Berlin Kladow":             []string{"14089"},
		"Berlin Siemensstadt":       []string{"13599", "13627", "13629"},
		"Berlin Spandau":            []string{"13581", "13583", "13585", "13587", "13597", "14052"},
		"Berlin Staaken":            []string{"13581", "13589", "13491", "13593"},
		"Berlin Wilhelmstadt":       []string{"13581", "13593", "13595", "13597"},
		// Steglitz-Zehlendorf
		"Berlin Dahlem":       []string{"12205", "14169", "14195", "14199"},
		"Berlin Lankwitz":     []string{"12167", "12209", "12247", "12249", "12277"},
		"Berlin Lichterfelde": []string{"12165", "12203", "12205", "12207", "12209", "12247", "12249", "12279", "14167", "14169", "14195"},
		"Berlin Nikolassee":   []string{"14109", "14129", "14163"},
		"Berlin Steglitz":     []string{"12157", "12161", "12163", "12165", "12167", "12169", "12203", "12247", "14195", "14197"},
		"Berlin Wannsee":      []string{"14109"},
		"Berlin Zehlendorf":   []string{"14129", "14163", "14165", "14167", "15169"},
		// Tempelhof-Schöneberg
		"Berlin Friedenau":   []string{"10827", "12159", "12161", "12163", "14197"},
		"Berlin Lichtenrade": []string{"12107", "12277", "12305", "12307", "12309"},
		"Berlin Mariendorf":  []string{"12099", "12105", "12107", "12109", "12277"},
		"Berlin Marienfelde": []string{"12107", "12249", "12277", "12279", "12307"},
		"Berlin Schöneberg":  []string{"10777", "10779", "10781", "10783", "10785", "10787", "10789", "10823", "10825", "10827", "10829", "10965", "12103", "12105", "12157", "12159"},
		"Berlin Tempelhof":   []string{"10965", "12099", "12101", "12103", "12105", "12109"},
		// Treptow-Köpenick
		"Berlin Adlershof":         []string{"12439", "12487", "12489"},
		"Berlin Alt-Treptow":       []string{"12435"},
		"Berlin Altglienicke":      []string{"12524", "12526"},
		"Berlin Baumschulenweg":    []string{"12437", "12487"},
		"Berlin Bohnsdorf":         []string{"12524", "12526"},
		"Berlin Friedrichshagen":   []string{"12587"},
		"Berlin Grünau":            []string{"12526", "12527"},
		"Berlin Johannisthal":      []string{"12437", "12487", "12489"},
		"Berlin Köpenick":          []string{"12555", "12557", "12559", "12587", "12623"},
		"Berlin Müggelheim":        []string{"12559"},
		"Berlin Niederschöneweide": []string{"12437", "12439"},
		"Berlin Oberschöneweide":   []string{"12459"},
		"Berlin Schöneweide":       []string{"12459", "12437", "12439"},
		"Berlin Plänterwald":       []string{"12435", "12437"},
		"Berlin Rahnsdorf":         []string{"12587", "12589"},
		"Berlin Schmöckwitz":       []string{"12527"},
	}

	pcode2district, pcode2citypart map[string]string
)

func isPostcode(location string) bool {

	if len(location) == 0 || len(location) > MAX_POSTCODE_LENGTH {
		return false
	}

	for i := 0; i < len(location); i++ {
		if location[i] < '0' || location[i] > '9' {
			return false
		}
	}

	return true
}

func Postcodes(location string) []string {

	postcodes := make(map[string]string)

	if isPostcode(location) {
		if len(location) == MAX_POSTCODE_LENGTH  {
			postcodes[location] = location
		} else {
			for _, districtCodes := range PostcodeMap {
				for _, districtCode := range districtCodes {
					if strings.HasPrefix(districtCode, location) {
						postcodes[districtCode] = districtCode
					}
				}
			}
		}
	} else {
		districts := make([]string, 0)
		districts = append(districts, location)

		if _, found := DistrictMap[location]; found {
			for _, district := range DistrictMap[location] {
				districts = append(districts, district)
			}
		}

		for _, district := range districts {
			if _, found := PostcodeMap[district]; found {
				for _, postcode := range PostcodeMap[district] {
					postcodes[postcode] = postcode
				}
			}
		}
	}

	result := make([]string, 0, len(postcodes))
	for postcode := range postcodes {
		result = append(result, postcode)
	}

	return result
}

func (tree *LocationTree) Autocomplete(prefix string) []string {

	result := make([]string, 0)

	tree.WalkPrefix(strings.ToLower(prefix), func(word string, s interface{}) bool {
		result = append(result, s.([]string)...)
		return false
	})

	return result
}

func (tree *LocationTree) Normalize(location string) string {

	value, found := tree.Get(strings.ToLower(location))
	if !found {
		return location
	}

	return value.([]string)[0]
}

func (tree *LocationTree) insert(word, s string) {

	var places []string
	value, found := tree.Get(word)
	if !found {
		places = make([]string, 0)
	} else {
		places = value.([]string)
	}

	found = false
	for _, place := range places {
		if place == s {
			found = true
			break
		}
	}

	if !found {
		places = append(places, s)
		tree.Insert(word, places)
	}
}

func (tree *LocationTree) Add(s string) {

	for _, token := range strings.Split(s, " ") {
		for _, word := range strings.Split(token, "-") {
			tree.insert(strings.ToLower(word), s)
		}
		tree.insert(strings.ToLower(token), s)
	}
	tree.insert(strings.ToLower(s), s)
}

func NewLocationTree(locations []string) *LocationTree {

	tree := &LocationTree{radix.New()}

	for _, location := range locations {
		tree.Add(location)
	}

	for district := range DistrictMap {
		tree.Add(district)
	}

	for citypart := range PostcodeMap {
		tree.Add(citypart)
	}

	return tree
}

func init() {

	pcode2district = make(map[string]string)
	pcode2citypart = make(map[string]string)
	for district, cityparts := range DistrictMap {
		for _, citypart := range cityparts {
			for _, postcode := range PostcodeMap[citypart] {
				pcode2district[postcode] = district
				pcode2citypart[postcode] = citypart
			}
		}
	}
}
