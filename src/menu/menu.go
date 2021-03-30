package menu

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"

	rice "github.com/GeertJohan/go.rice"
	"github.com/manifoldco/promptui"

	utils "nc-shell/src/utils"
)

type revshell struct {
	Type    string `json:"type"`
	Note    string `json:"note"`
	Payload string `json:"payload"`
}

func interfaceMenu() string {
	interfaces := utils.ListInterfaces()

	templates := &promptui.SelectTemplates{
		Label:    "{{ . }}",
		Active:   "\U0001F336 {{ .Name | cyan }} ({{ .IP | red }})",
		Inactive: "  {{ .Name | cyan }} ({{ .IP | red }})",
		Selected: "\U0001F336 {{ .Name | green }}",
	}

	prompt := promptui.Select{
		Label:     "Interface: ",
		Items:     interfaces,
		Templates: templates,
		Size:      4,
	}

	i, _, err := prompt.Run()

	if err != nil {
		os.Exit(1)
	}
	return interfaces[i].IP
}

func portMenu() int {
	var port int
	templates := &promptui.SelectTemplates{
		Label:    "{{ . }}",
		Active:   "\U0001F336 {{ . | cyan }}",
		Inactive: "  {{ . | cyan }}",
		Selected: "\U0001F336 {{ . | green }}",
	}

	prompt := promptui.Select{
		Label:     "Port:",
		Items:     []string{"HTTP (80)", "HTTPS (443)", "DNS (53)", "Custom"},
		Templates: templates,
	}

	_, result, err := prompt.Run()

	if err != nil {
		os.Exit(1)
	}

	if result == "Custom" {
		port = customPortMenu()
	} else {
		r := regexp.MustCompile(`\w+(\d+)`)
		port, _ = strconv.Atoi(r.FindStringSubmatch(result)[0])
	}
	return port
}

func customPortMenu() int {
	var port int
	validate := func(input string) error {
		_, err := strconv.Atoi(input)
		return err
	}

	templates := &promptui.PromptTemplates{
		Prompt:  "{{ . }} ",
		Valid:   "{{ . | green }} ",
		Invalid: "{{ . | red }} ",
		Success: "{{ . | green }} ",
	}

	prompt := promptui.Prompt{
		Label:     "Port:",
		Templates: templates,
		Validate:  validate,
	}

	result, err := prompt.Run()

	if err != nil {
		os.Exit(1)
	}
	port, _ = strconv.Atoi(result)
	return port
}

func revshellMenu() string {
	//Tricks to use the rice box static multiple times
	templateBox, err := rice.FindBox(".././static/")
	if err != nil {
		log.Fatal(err)
	}
	// get file contents as string
	jsonFile, err := templateBox.Open("data.json")
	if err != nil {
		log.Fatal(err)
	}
	defer jsonFile.Close()
	byteValue, _ := ioutil.ReadAll(jsonFile)
	var revshells []revshell
	var typeRevShell []string
	json.Unmarshal(byteValue, &revshells)

	for _, revshell := range revshells {
		if !utils.SliceStringContains(typeRevShell, revshell.Type) {
			typeRevShell = append(typeRevShell, revshell.Type)
		}

	}

	templates := &promptui.SelectTemplates{
		Label:    "{{ . }}",
		Active:   "\U0001F336 {{ . | cyan }}",
		Inactive: "  {{ . | cyan }}",
		Selected: "\U0001F336 {{ . | green }}",
	}

	prompt := promptui.Select{
		Label:     "Reverse shell: ",
		Items:     typeRevShell,
		Templates: templates,
		Size:      20,
	}

	_, shell, err := prompt.Run()

	if err != nil {
		os.Exit(1)
	}

	return shell
}

func cliRevshellMenu(ip string, port int, shell string) {
	//Tricks to use the rice box static multiple times
	templateBox, err := rice.FindBox(".././static/")
	if err != nil {
		log.Fatal(err)
	}
	// get file contents as string
	jsonFile, err := templateBox.Open("data.json")
	if err != nil {
		log.Fatal(err)
	}
	defer jsonFile.Close()
	byteValue, _ := ioutil.ReadAll(jsonFile)
	var revshells []revshell
	var onelinersRevshell []revshell
	json.Unmarshal(byteValue, &revshells)
	i := 0
	for _, revshell := range revshells {
		if revshell.Type == shell {
			revshell.Payload = strings.ReplaceAll(revshell.Payload, "{LHOST}", ip)
			revshell.Payload = strings.ReplaceAll(revshell.Payload, "{LPORT}", strconv.Itoa(port))
			onelinersRevshell = append(onelinersRevshell, revshell)
			fmt.Println("[" + strconv.Itoa(i) + "] " + revshell.Payload + "\n")
			i++
		}
	}

}

// Menu to choose the interface and the kind of shell. Oneliners will be print
func Menu(port int) {
	ip := interfaceMenu()
	shell := revshellMenu()
	cliRevshellMenu(ip, port, shell)
}
