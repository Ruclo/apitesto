package config
import (
	"fmt"
	"io/ioutil"
	"gopkg.in/yaml.v2"
	"strings"
)
type HTTPMethod string
type HeaderType int
type HeadersMap map[string]string

func (hm *HeadersMap) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var rawHeaders []string
	if err := unmarshal(&rawHeaders); err != nil {
		return err
	}

	*hm = make(HeadersMap)

	for _, headerStr := range rawHeaders {
		if !strings.Contains(headerStr, ":") {
			return fmt.Errorf("invalid header format: %s", headerStr)
		}
		parts := strings.SplitN(headerStr, ":", 2)
		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])
		(*hm)[key] = value
	}

	return nil
}

const (
    GET     HTTPMethod = "GET"
    POST    HTTPMethod = "POST"
    PUT     HTTPMethod = "PUT"
    DELETE  HTTPMethod = "DELETE"
    PATCH   HTTPMethod = "PATCH"
    OPTIONS HTTPMethod = "OPTIONS"
    HEAD    HTTPMethod = "HEAD"
)

type Config struct {
	BaseURLs			map[string]string		`yaml:"base-urls,omitempty"`	
	Groups				map[string]GroupConfig	`yaml:"groups,omitempty"`
	FunctionalTests		[]FunctionalTest 		`yaml:"functional-tests,omitempty"`
	LoadTests			[]LoadTest				`yaml:"load-tests,omitempty"`
}

type BaseTest struct {
	Name			string				`yaml:"name"`
	URL				string				`yaml:"url"`
	Method			HTTPMethod 			`yaml:"method"`
	Headers			HeadersMap 			`yaml:"request-headers,omitempty"`
	Body			string				`yaml:"body,omitempty"`
	QueryParams 	map[string]string 	`yaml:"query-parameters,omitempty"`
	PathParams 		map[string]string 	`yaml:"path-parameters,omitempty"`
	Groups			[]string			`yaml:"groups,omitempty"`
	Timeout			Timeout				`yaml:"timeout,omitempty"`
}

type GroupConfig struct {
	/*URL				string				`yaml:"url,omitempty"`
	Method			HTTPMethod 			`yaml:"method,omitempty"`
	Headers			HeadersMap 			`yaml:"request-headers,omitempty"`
	Body			string				`yaml:"body,omitempty"`
	QueryParams 	map[string]string 	`yaml:"query-parameters,omitempty"`
	PathParams 		map[string]string 	`yaml:"path-parameters,omitempty"`
	Timeout			Timeout				`yaml:"timeout,omitempty"`
	Assertions		[]string			`yaml:"assertions,omitempty"`
	SuccessCriteria []string			`yaml:"success-criteria,omitempty`
	Expected		ExpectedResponse	`yaml:"expected,omitempty"`
	Phases			[]Phase				`yaml:"phases,omitempty"`
	ThinkTime		ThinkTime			`yaml:"think-time,omitempty`*/
	FunctionalTest
	LoadTest
}

type FunctionalTest struct {
	BaseTest
	Expected		ExpectedResponse	`yaml:"expected"`
}

type LoadTest struct {
	BaseTest
	Phases			[]Phase		`yaml:"phases"`
	ThinkTime		ThinkTime	`yaml:"think-time,omitempty`
	SuccessCriteria []string	`yaml:"success-criteria,omitempty`
}

type Phase struct {
	Name 				string 		`yaml:"name"`
	Duration 			string 		`yaml:"duration"`
	RequestsPerSecond	int			`yaml:"rps"`
	ConcurrentUsers		int 		`yaml:"ccs"`
}

type ThinkTime struct {
	Min int `yaml:"min"`
	Max int `yaml:"max"`

}

type Timeout struct {
	Connect	int
	Request int
}


const (
	HeaderRequired	HeaderType = iota
	HeaderForbidden
	HeaderValueRequired
	HeaderValueForbidden
)

type ExpectedHeader struct {
	Key		string
	Type	HeaderType
	Value	string
}

func (eh *ExpectedHeader) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var rawHeader interface{}
	if err := unmarshal(&rawHeader); err != nil {
		return err
	}

	mapValue, ok := rawHeader.(map[interface{}]interface{})
	if (!ok) {
		return fmt.Errorf("unexpected header format %s", rawHeader);
	}

	for k,v := range mapValue {
		key, ok := k.(string)
		if !ok {
			return fmt.Errorf("wrong header format %s", k)
		}

		value, ok := v.(string)
		
		if !ok {
			return fmt.Errorf("wrong header format %s", v)
		}


	}
	return nil
	/*if strings.Contains(headerStr, ":") {
		if strings.HasPrefix(headerStr, "!") {
			return fmt.Errorf("expected header format wrong, starts with ! and contains =, %s", headerStr)
		}

		parts := strings.SplitN(headerStr, ":", 2)
		if len(parts) != 2 {
			return fmt.Errorf("invalid expected header format: %s", headerStr)
		}
		
		eh.Key = parts[0]
		eh.Value = strings.TrimPrefix(parts[1], "!")

		if len(eh.Value) == 0 {
			return fmt.Errorf("invalid expected header format: %s", headerStr)
		}
		
		if strings.HasPrefix(parts[1], "!") {
			eh.Type = HeaderValueForbidden 
		} else {
			eh.Type = HeaderValueRequired
		}
		
		
	} else {
		eh.Key = strings.TrimPrefix(headerStr, "!")
		if len(eh.Key) == 0 {
			return fmt.Errorf("invalid expected header format: %s", headerStr)
		}
		
		if strings.HasPrefix(headerStr, "!") {
			eh.Type = HeaderForbidden
		} else {
			eh.Type = HeaderRequired 
		}
	}
	return nil*/
	
}

type ExpectedResponse struct {
	Status			int 					`yaml:"status"`
	Headers			[]ExpectedHeader 		`yaml:"headers,omitempty"`
	ResponseSchema 	map[string]interface{} 	`yaml:"response-schema,omitempty"`
	Assertions		[]string				`yaml:"assertions,omitempty"`
}

func LoadConfigFromYAML(filePath string) (*Config, error) {
	yamlData, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	var config Config
	err = yaml.Unmarshal(yamlData, &config)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return &config, nil

}