package dotenv

import (
	"bufio"
	"os"
	"strconv"
	"strings"
	"time"
)

type Env struct {
	Variables map[string][]string
	loaded    bool
}

func (e *Env) GetAll(key string) []string {
	if !e.loaded {
		panic("env not loaded")
	}
	return e.Variables[key]
}

func (e *Env) Get(key string) string {
	if !e.loaded {
		panic("env not loaded")
	}
	return e.Variables[key][0]
}

func (e *Env) GetDefault(key, def string) string {
	if !e.loaded {
		panic("env not loaded")
	}
	if val, ok := e.Variables[key]; ok {
		return val[0]
	}
	return def
}

func (e *Env) GetBool(key string) bool {
	if !e.loaded {
		panic("env not loaded")
	}
	var b, err = strconv.ParseBool(e.Variables[key][0])
	if err != nil {
		panic(err)
	}
	return b
}

func (e *Env) GetInt(key string) int {
	if !e.loaded {
		panic("env not loaded")
	}
	var i, err = strconv.Atoi(e.Variables[key][0])
	if err != nil {
		panic(err)
	}
	return i
}

var env Env = Env{
	Variables: make(map[string][]string),
}

func GetAll(key string, def ...string) []string {
	if !env.loaded {
		env.Load(".env")
	}
	if val, ok := env.Variables[key]; ok {
		return val
	}
	if len(def) > 0 {
		return def
	}
	return nil
}

func Get(key string, def ...string) string {
	if !env.loaded {
		env.Load(".env")
	}
	if val, ok := env.Variables[key]; ok {
		return val[0]
	}
	if len(def) > 0 {
		return def[0]
	}
	return ""
}

func GetBool(key string, def ...bool) bool {
	if !env.loaded {
		env.Load(".env")
	}
	if val, ok := env.Variables[key]; ok {
		var b, err = strconv.ParseBool(val[0])
		if err != nil {
			panic(err)
		}
		return b
	}
	if len(def) > 0 {
		return def[0]
	}
	return false
}

func GetInt(key string, def ...int) int {
	if !env.loaded {
		env.Load(".env")
	}
	if val, ok := env.Variables[key]; ok {
		var i, err = strconv.Atoi(val[0])
		if err != nil {
			panic(err)
		}
		return i
	}
	if len(def) > 0 {
		return def[0]
	}
	return 0
}

func GetTimeDuration(key string, def ...time.Duration) time.Duration {
	if !env.loaded {
		env.Load(".env")
	}
	if val, ok := env.Variables[key]; ok {
		var d, err = ParseDuration(val[0])
		if err != nil {
			goto retDef
		}
		return d
	}
retDef:
	if len(def) > 0 {
		return def[0]
	}
	return 0
}

func Load(path string) {
	env.Load(path)
}

func LoadString(s string) {
	env.LoadString(s)
}

func Unmarshal(s ...interface{}) {
	env.Unmarshal(s...)
}

func (e *Env) Unmarshal(s ...interface{}) error {
	if !e.loaded {
		panic("env not loaded")
	}
	return e.unmarshal(s...)
}

func (e *Env) Load(path string) error {
	e.loaded = true
	e.Variables = make(map[string][]string)
	var lines, err = readLines(path)
	if err != nil {
		return err
	}
	e.loadLines(lines)
	return nil
}

func (e *Env) LoadString(s string) {
	e.loaded = true
	e.Variables = make(map[string][]string)
	e.loadLines(strings.Split(s, "\n"))
}

//	// Capture variables in the form of:
//	// $VAR
//	// $var[1]
//	// $PARAM[1:3]
//	// $VAR[1:]
//	// $DATA[:3]
//	var varRegex = regexp.MustCompile(`\$\w+(\[\d+\]|\[\d+:\d+\]|\[\d+:\]|\[:\d+\]|)`)

func (e *Env) loadLines(lines []string) {
	for _, line := range lines {
		line = stripComments(line)
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		var parts = strings.SplitN(line, "=", 2)
		var key = strings.TrimSpace(parts[0])
		var val = strings.Split(parts[1], ",")
		for i, v := range val {
			// Check if the value is a string, if so remove the quotes
			// and replace escaped quotes with unescaped quotes
			v = strings.TrimSpace(v)
			if v != "" {
				if v[0] == '"' && v[len(v)-1] == '"' ||
					v[0] == '\'' && v[len(v)-1] == '\'' ||
					v[0] == '`' && v[len(v)-1] == '`' {

					v = v[1 : len(v)-1]
				}
				v = strings.ReplaceAll(v, "\"", "")
				v = strings.ReplaceAll(v, "'", "")
				v = strings.ReplaceAll(v, "`", "")

				v = strings.TrimSpace(v)
			}
			if v == "" {
				// Remove the key if the value is empty
				val[i] = ""
			}

			if newv := strings.ToLower(v); newv == "null" || newv == "nil" || newv == "none" {
				v = ""
			}
			if strings.HasPrefix(v, "$") {
				v = e.Get(v[1:])
			}
			val[i] = v
		}
		if len(val) == 1 && val[0] == "" {
			continue
		}
		e.Variables[key] = append(e.Variables[key], val...)
	}
}

func readLines(path string) ([]string, error) {
	var file, err = os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	var lines []string
	var scanner = bufio.NewScanner(file)
	for scanner.Scan() {
		var line = scanner.Text()
		lines = append(lines, line)
	}
	return lines, nil
}

func stripComments(line string) string {
	var commentOpen = false
	var commentChars []rune = []rune{'#', ';'}
	var quoteChars []rune = []rune{'"', '\'', '`'}
	var quoteChar rune
	var quoteOpen = false
	var quoteEscaped = false
	var newLineChars []rune
	for _, char := range line {
		if quoteOpen {
			if quoteEscaped {
				quoteEscaped = false
			} else if char == '\\' {
				quoteEscaped = true
			} else if char == quoteChar {
				quoteOpen = false
			}
			newLineChars = append(newLineChars, char)
			continue
		}
		if commentOpen {
			if char == '\r' || char == '\n' {
				commentOpen = false
			}
			continue
		}
		if char == '\r' || char == '\n' {
			break
		}
		for _, quoteChar = range quoteChars {
			if char == quoteChar {
				quoteOpen = true
				break
			}
		}
		if quoteOpen {
			newLineChars = append(newLineChars, char)
			continue
		}
		for _, commentChar := range commentChars {
			if char == commentChar {
				commentOpen = true
				break
			}
		}
		if commentOpen {
			continue
		}
		newLineChars = append(newLineChars, char)
	}
	return string(newLineChars)
}
