// Package spellcheck provides markdown spellchecking capabilities.
package spellcheck

import (
	"fmt"
	"strings"
	"unicode"

	"github.com/EdgarOrtegaRamirez/markdownforge/internal/parser"
)

// Issue represents a spelling issue found in the document.
type Issue struct {
	Line       int
	Word       string
	Suggestion string
	Message    string
}

// Checker checks markdown documents for common spelling issues.
type Checker struct {
	rules map[string]string
}

// NewChecker creates a new spell checker with common English misspellings.
func NewChecker() *Checker {
	c := &Checker{
		rules: make(map[string]string),
	}
	c.AddDefaultRules()
	return c
}

// AddDefaultRules adds common English misspellings.
func (c *Checker) AddDefaultRules() {
	rules := map[string]string{
		"acheive":        "achieve",
		"accomodate":     "accommodate",
		"acurately":      "accurately",
		"agressive":      "aggressive",
		"allign":         "align",
		"alot":           "a lot",
		"alright":        "all right",
		"analyse":        "analyze",
		"aparent":        "apparent",
		"aplied":         "applied",
		"arguement":      "argument",
		"arythmatic":     "arithmetic",
		"attemp":         "attempt",
		"attened":        "attended",
		"becuase":        "because",
		"beleive":        "believe",
		"buisness":       "business",
		"calender":       "calendar",
		"catagory":       "category",
		"certian":        "certain",
		"challange":      "challenge",
		"changeble":      "changeable",
		"cheif":          "chief",
		"choosen":        "chosen",
		"civilisation":   "civilization",
		"collegue":       "colleague",
		"commited":       "committed",
		"complitely":     "completely",
		"comprehenssion": "comprehension",
		"concious":       "conscious",
		"confusioned":    "confused",
		"connnection":    "connection",
		"contians":       "contains",
		"continous":      "continuous",
		"contructor":     "constructor",
		"corisponding":    "corresponding",
		"corupted":       "corrupted",
		"coustomer":      "customer",
		"develope":       "develop",
		"devided":        "divided",
		"diffrence":      "difference",
		"dissapear":      "disappear",
		"divison":        "division",
		"embarass":       "embarrass",
		"embarressed":    "embarrassed",
		"enviroment":     "environment",
		"exagerate":      "exaggerate",
		"exellence":      "excellence",
		"existance":      "existence",
		"exmaple":        "example",
		"experiance":     "experience",
		"exsist":         "exist",
		"faield":         "failed",
		"familar":        "familiar",
		"faviour":        "favourite",
		"fianlly":        "finally",
		"finacial":       "financial",
		"foriegn":        "foreign",
		"fortuantely":    "fortunately",
		"fourty":         "forty",
		"freind":         "friend",
		"functon":        "function",
		"functons":       "functions",
		"gaurentee":      "guarantee",
		"generaly":       "generally",
		"gettin":         "getting",
		"goverment":      "government",
		"grammer":        "grammar",
		"happend":        "happened",
		"heros":          "heroes",
		"humourous":      "humorous",
		"hypocracy":      "hypocrisy",
		"identifcation":  "identification",
		"imaginery":      "imaginary",
		"imediately":     "immediately",
		"incomming":      "incoming",
		"independant":    "independent",
		"infomation":     "information",
		"initialy":       "initially",
		"inteligent":     "intelligent",
		"intelllectual":  "intellectual",
		"interuption":    "interruption",
		"intresting":     "interesting",
		"knowlege":       "knowledge",
		"lenghth":        "length",
		"liason":         "liaison",
		"libary":         "library",
		"licence":        "license",
		"lightining":     "lightning",
		"lisence":        "license",
		"loosing":        "losing",
		"maintenence":    "maintenance",
		"managable":      "manageable",
		"manuever":       "maneuver",
		"mateiral":       "material",
		"memorandum":     "memorandum",
		"mispell":        "misspell",
		"mispelled":      "misspelled",
		"mispelling":     "misspelling",
		"mornig":         "morning",
		"mulitple":       "multiple",
		"neccessary":     "necessary",
		"neighbour":      "neighbor",
		"necesary":       "necessary",
		"noticable":      "noticeable",
		"occurance":      "occurrence",
		"occurrs":        "occurs",
		"offical":        "official",
		"ocassion":       "occasion",
		"ocurr":          "occur",
		"ocurrred":       "occurred",
		"ocurrrence":     "occurrence",
		"ocurrrences":    "occurrences",
		"ommit":          "omit",
		"ommitted":       "omitted",
		"ommiting":       "omitting",
		"oposite":        "opposite",
		"organise":       "organize",
		"orginally":      "originally",
		"ouput":          "output",
		"palatte":        "palette",
		"particulary":    "particularly",
		"particuarly":    "particularly",
		"passionnate":    "passionate",
		"pavement":       "pavement",
		"peice":          "piece",
		"peformed":       "performed",
		"peformance":     "performance",
		"performence":    "performance",
		"permenant":      "permanent",
		"perphas":        "perhaps",
		"persue":         "pursue",
		"posession":      "possession",
		"posessions":     "possessions",
		"posses":         "possess",
		"practise":       "practice",
		"predominately":  "predominantly",
		"premissions":    "permissions",
		"preogative":     "prerogative",
		"privelege":      "privilege",
		"privalege":      "privilege",
		"professer":      "professor",
		"profil":         "profile",
		"programable":    "programmable",
		"programer":      "programmer",
		"progrm":         "program",
		"pronounciation": "pronunciation",
		"promary":        "primary",
		"propogate":      "propagate",
		"proposeable":    "proposable",
		"propotion":      "proportion",
		"pubilc":         "public",
		"publically":     "publicly",
		"pupose":         "purpose",
		"recieve":        "receive",
		"recievers":      "receivers",
		"recomend":       "recommend",
		"recomended":     "recommended",
		"recomends":      "recommends",
		"reffered":       "referred",
		"reffering":      "referred",
		"refrences":      "references",
		"reight":         "right",
		"relevent":       "relevant",
		"religous":       "religious",
		"remeber":        "remember",
		"removables":     "removable",
		"reoccurrence":   "recurrence",
		"repatition":     "repetition",
		"repetead":       "repeated",
		"repeteadly":     "repeatedly",
		"repitivly":      "repetitively",
		"represtnation":  "representation",
		"representitive": "representative",
		"reponsible":     "responsible",
		"resistence":     "resistance",
		"resistent":      "resistant",
		"resoultion":     "resolution",
		"resourses":      "resources",
		"respose":        "response",
		"restaraunt":     "restaurant",
		"retreive":       "retrieve",
		"rythem":         "rhythm",
		"rythm":          "rhythm",
		"sacriligious":   "sacrilegious",
		"salendar":       "calendar",
		"satifaction":    "satisfaction",
		"satify":         "satisfy",
		"satistfying":    "satisfying",
		"seach":          "search",
		"securty":        "security",
		"sentance":       "sentence",
		"seperate":       "separate",
		"serch":          "search",
		"shouldnt":       "shouldn't",
		"similiar":       "similar",
		"simpatico":      "simpatia",
		"speach":         "speech",
		"spelling":       "spelling",
		"spilt":          "spill",
		"staion":         "station",
		"statment":       "statement",
		"stauts":         "status",
		"stictly":        "strictly",
		"strengh":        "strength",
		"strenght":       "strength",
		"stregth":        "strength",
		"strucure":       "structure",
		"subconsiously":  "subconsciously",
		"succeded":       "succeeded",
		"succesful":      "successful",
		"sucess":         "success",
		"sucessful":      "successful",
		"suceed":         "succeed",
		"supercede":      "supersede",
		"suprise":        "surprise",
		"suported":       "supported",
		"surround":       "surround",
		"surounding":     "surrounding",
		"surounded":      "surrounded",
		"surroudings":    "surroundings",
		"taht":           "that",
		"teh":            "the",
		"temperture":     "temperature",
		"therefor":       "therefore",
		"threshhold":     "threshold",
		"tommorrow":      "tomorrow",
		"tonite":         "tonight",
		"tommorow":       "tomorrow",
		"tomorow":        "tomorrow",
		"tounge":         "tongue",
		"touble":         "trouble",
		"truely":         "truly",
		"unforseen":      "unforeseen",
		"unfortunatly":   "unfortunately",
		"unforutnately":  "unfortunately",
		"uninterrupted":  "uninterrupted",
		"unneccessary":   "unnecessary",
		"unnecesary":     "unnecessary",
		"unnessary":      "unnecessary",
		"unsuccessfull":  "unsuccessful",
		"unsuccesssful":  "unsuccessful",
		"upcomming":      "upcoming",
		"usful":          "useful",
		"usefull":        "useful",
		"vaccuum":        "vacuum",
		"vehical":        "vehicle",
		"vehsicle":       "vehicle",
		"versus":         "versus",
		"verison":        "version",
		"viably":         "viable",
		"visious":        "vicious",
		"visting":        "visiting",
		"visted":         "visited",
		"wether":         "whether",
		"wich":           "which",
		"wierd":          "weird",
		"woudl":          "would",
		"writeable":      "writable",
		"wrtie":          "write",
		"yesterda":       "yesterday",
		"yesteday":       "yesterday",
		"yestearday":     "yesterday",
	}
	for k, v := range rules {
		c.rules[k] = v
	}
}

// Check runs spell check on a parsed document.
func (c *Checker) Check(doc *parser.Document) []Issue {
	var issues []Issue
	for lineNum, line := range doc.Lines {
		lineNum++
		words := tokenize(line)
		for _, word := range words {
			lower := strings.ToLower(word)
			if suggestion, ok := c.rules[lower]; ok {
				issues = append(issues, Issue{
					Line:       lineNum,
					Word:       word,
					Suggestion: suggestion,
					Message:    fmt.Sprintf("Possible misspelling: %s → %s", word, suggestion),
				})
			}
		}
	}
	return issues
}

// Count returns the number of issues found.
func (c *Checker) Count(issues []Issue) int {
	return len(issues)
}

// tokenize splits a line into words for spell checking.
func tokenize(line string) []string {
	var words []string
	var current strings.Builder

	for _, r := range line {
		if unicode.IsLetter(r) || unicode.IsDigit(r) || r == '\'' {
			current.WriteRune(r)
		} else {
			if current.Len() > 0 {
				words = append(words, current.String())
				current.Reset()
			}
		}
	}
	if current.Len() > 0 {
		words = append(words, current.String())
	}
	return words
}
