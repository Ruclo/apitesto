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

	if mapValue, ok := rawHeader.(map[interface{}]interface{}); ok {
		for k, v := range mapValue {
			key, ok := k.(string)
			if !ok {
				return fmt.Errorf("wrong header format %s", rawHeader)
			}
			value, ok := v.(string)

			if !ok {
				return fmt.Errorf("wrong header format %s", rawHeader)
			}
			eh.Key = strings.TrimPrefix(key, "!")
			eh.Value = strings.TrimPrefix(value, "!")
			eh.Type = HeaderValueRequired

			if strings.HasPrefix(value, "!") {
				eh.Type = HeaderValueForbidden
			}

			if strings.HasPrefix(key, "!") {
				return fmt.Errorf("wrong header format %s", rawHeader)
			}
		}
			

		return nil
	}
	
	if strHeader, ok := rawHeader.(string); ok {
			eh.Key = strings.TrimPrefix(strHeader, "!")
			
			switch strings.HasPrefix(strHeader, "!") {
			case true:
				eh.Type = HeaderRequired
			case false:
				eh.Type = HeaderForbidden
			}
			
		return nil
	}
	return fmt.Errorf("unexpected header format %s", rawHeader);
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