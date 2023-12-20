package main

import (
	"fmt"
	"log"
	"os"
	"slices"
	"strings"
)

func main() {
	input, err := readInput(2023, 20)
	if err != nil {
		fmt.Printf("Failed to read input file: %v\n", err)
		return
	}

	modulesToMonitor := []string{"tx", "gc", "kp", "vg"}
	monitorSignals := make(map[string][]int)
	modules := make(map[string]Module)
	lines := strings.Split(input, "\n")
	for _, line := range lines {
		parts := strings.Split(line, " -> ")
		module := parts[0]
		destinations := strings.Split(parts[1], ", ")

		m := Module{}
		if module == "broadcaster" {
			m.moduleType = "broadcaster"
			m.moduleName = "broadcaster"
		} else {
			m.moduleType = module[0:1]
			m.moduleName = module[1:]
		}
		m.destinationModuleNames = destinations
		modules[m.moduleName] = m
	}

	modules["button"] = CreateButtonModule()
	modules["output"] = CreateOutputModule()

	// Connect everything
	for _, module := range modules {
		destinations := module.destinationModuleNames
		for _, destination := range destinations {
			m := modules[destination]
			m.inputModuleNames = append(m.inputModuleNames, module.moduleName)
			modules[destination] = m
		}
	}

	low, high, bp := 0, 0, 1
	MaxPressPartOne := 1000
	MaxPressPartTwo := 10000

	unprocessed := []Action{CreateButtonPressAction()}
	for len(unprocessed) > 0 {
		action := unprocessed[0]
		unprocessed = unprocessed[1:]

		// Part One
		if action.destination != "button" && bp <= MaxPressPartOne {
			if action.signal {
				high++
			} else {
				low++
			}
		}

		//Part 2
		if (slices.Contains(modulesToMonitor, action.destination)) && !action.signal {
			monitorSignals[action.destination] = append(monitorSignals[action.destination], bp)
		}

		destinationModule := modules[action.destination]
		actions := destinationModule.processAction(action, modules)
		unprocessed = append(unprocessed, actions...)

		if len(unprocessed) == 0 && bp < MaxPressPartTwo {
			unprocessed = append(unprocessed, CreateButtonPressAction())
			bp++
		}
	}

	ansP1 := low * high

	var nums []int
	for _, v := range monitorSignals {
		nums = append(nums, v[0])
	}
	ansP2 := lcmArr(nums)

	fmt.Println("Solution Part 1:", ansP1) //681194780
	fmt.Println("Solution Part 2:", ansP2) //238593356738827
}

func CreateOutputModule() Module {
	return Module{
		moduleName:             "output",
		moduleType:             "output",
		destinationModuleNames: []string{},
	}
}

func CreateButtonModule() Module {
	return Module{
		moduleName:             "button",
		moduleType:             "button",
		destinationModuleNames: []string{"broadcaster"},
	}
}

func CreateButtonPressAction() Action {
	buttonPress := Action{
		signal:      false,
		destination: "broadcaster",
	}
	return buttonPress
}

type Action struct {
	signal      bool
	destination string
}

type Module struct {
	moduleName             string
	moduleType             string
	destinationModuleNames []string
	inputModuleNames       []string
	state                  bool
}

func (m *Module) processAction(action Action, modules map[string]Module) []Action {
	var actions []Action
	propagate := false
	switch m.moduleType {
	case "broadcaster":
		m.state = action.signal
		propagate = true
	case "%":
		// Set state to flip if input is LOW
		if action.signal == false {
			m.state = !m.state
			propagate = true
		}
	case "&":
		//conjunction modules send LOW input only if ALL inputs are HIGH
		newState := false
		for _, input := range m.inputModuleNames {
			module := modules[input]
			if module.state == false {
				newState = true
				break
			}
		}
		m.state = newState
		propagate = true

	case "output":
		m.state = action.signal
	}

	modules[m.moduleName] = *m

	if propagate {
		for _, destination := range m.destinationModuleNames {
			actions = append(actions, Action{
				signal:      m.state,
				destination: destination,
			})
		}
	}

	return actions
}

func lcmArr(nums []int) int {
	var res = 1
	for _, num := range nums {
		res = lcm(res, num)
	}
	return res
}

func lcm(a, b int) int {
	return a * b / gcd(a, b)
}

func gcd(a, b int) int {
	for b != 0 {
		a, b = b, a%b
	}
	return a
}

// Boilerplate: ReadInput reads the input file for the given year and day
func readInput(year, day int) (string, error) {
	filePath := fmt.Sprintf("aoc-%d/Day%d/input.txt", year, day)
	file, err := os.ReadFile(filePath)
	if err != nil {
		log.Fatal(err)
	}

	return string(file), nil
}
