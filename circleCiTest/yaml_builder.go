//Create a script that will build a dynamic circleci yaml file based on the parameters passed in

package main

import (
	"bytes"
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

func mergeMap(dst, src map[string]interface{}) {
	for k, v := range src {
		dst[k] = v
	}
}

func main() {
	parameters := []string{"github", "sonarqube"} // This can be set based on user input or other logic


	others := make(map[string]interface{})

	base_yaml_data, err := os.ReadFile("fragments/base.yml")
	if err != nil {
		fmt.Printf("Error reading base.yml fragment: %v\n", err)
		return
	}

	base_yaml := make(map[string]interface{})
	err = yaml.Unmarshal(base_yaml_data, &base_yaml)
	if err != nil {
		fmt.Printf("Error parsing base.yml fragment: %v\n", err)
		return
	}

	// Read and merge each fragment into the appropriate bucket
	for _, p := range parameters {
		fragmentPath := fmt.Sprintf("fragments/%s.yml", p)
		data, err := os.ReadFile(fragmentPath)
		if err != nil {
			fmt.Printf("Error reading %s fragment: %v\n", p, err)
			return
		}

		fragment := make(map[string][]interface{})
		err = yaml.Unmarshal(data, &fragment)
		if err != nil {
			fmt.Printf("Error parsing %s fragment: %v\n", p, err)
			return
		}
		
		for k, v := range fragment {
			switch k {
			case "jobs":
				base_yaml["jobs"][0] = v
				for jobKey := range v.(map[string]interface{}) {
					base_yaml["workflows"][0] = jobKey
					break //only need the first job key to add to workflows
				}
			default:
				others[k] = v
			}
		}
	}

	var buf bytes.Buffer
	// Write version first
	buf.WriteString("version: 2.1\n\n")

	b, err := yaml.Marshal(base_yaml)
	if err != nil {
		fmt.Println("Error marshaling jobs:", err)
		return
	}
	buf.Write(b)
	buf.WriteString("\n")

	// Marshal any other top-level keys
	if len(others) > 0 {
		b, err := yaml.Marshal(others)
		if err != nil {
			fmt.Println("Error marshaling other keys:", err)
			return
		}
		buf.Write(b)
		buf.WriteString("\n")
	}

	// Write to file
	err_fi := os.WriteFile("config.yml", buf.Bytes(), 0644)
	if err_fi != nil {
		fmt.Println("Error writing to file:", err)
		return
	}

	fmt.Println("CircleCI config.yml generated successfully!")
}
