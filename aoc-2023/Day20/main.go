package main

import (
	"fmt"
	"log"
	"os"
	"strings"
)

func main() {
	input, err := readInput(2023, 20)
	if err != nil {
		fmt.Printf("Failed to read input file: %v\n", err)
		return
	}

	/*
	   broadcaster -> a, b, c
	   %a -> b
	   %b -> c
	   %c -> inv
	   &inv -> a
	*/

	periods := make(map[string][]int)

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
		//fmt.Println(parts)
		//fmt.Println(m)
		modules[m.moduleName] = m
	}

	modules["button"] = Module{
		moduleName:             "button",
		moduleType:             "button",
		destinationModuleNames: []string{"broadcaster"},
	}
	modules["output"] = Module{
		moduleName:             "output",
		moduleType:             "output",
		destinationModuleNames: []string{},
	}

	// Connect everything
	for _, module := range modules {
		destinations := module.destinationModuleNames
		for _, destination := range destinations {
			m := modules[destination]
			m.inputModuleNames = append(m.inputModuleNames, module.moduleName)
			modules[destination] = m
		}
	}

	low, high, buttonPresses := 0, 0, 1
	//MAX_BUTTON_PRESSES := 1000
	MAX_BUTTON_PRESSES := 1000000

	buttonPress := Action{
		signal:      false,
		destination: "broadcaster",
	}

	unprocessed := []Action{buttonPress}
	for len(unprocessed) > 0 {
		action := unprocessed[0]
		unprocessed = unprocessed[1:]

		if action.destination != "button" {
			if action.signal == true {
				high++
			} else {
				low++
			}
		}

		if action.destination == "rx" && action.signal == false {
			println("Button Presses Reached:", buttonPresses)
			break
		}

		if (action.destination == "tx" || action.destination == "gc" || action.destination == "kp" || action.destination == "vg") && action.signal == false {
			periods[action.destination] = append(periods[action.destination], buttonPresses)
		}

		destinationModule := modules[action.destination]
		actions := destinationModule.processAction(action, modules)
		unprocessed = append(unprocessed, actions...)

		if len(unprocessed) == 0 && buttonPresses < MAX_BUTTON_PRESSES {
			buttonPresses++
			unprocessed = append(unprocessed, buttonPress)
			//println("Button pressed:", buttonPresses)
		}
	}

	//println("Low:", low, "High:", high)
	ansP1 := low * high

	nums := []int{}
	for _, v := range periods {
		nums = append(nums, v[0])
	}

	ansP2 := lcmArr(nums)

	// Solution here
	fmt.Println("Solution Part 1:", ansP1)
	fmt.Println("Solution Part 2:", ansP2)
}

func lcmArr(nums []int) int {
	var res int = 1
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

type Module struct {
	moduleName             string
	moduleType             string
	destinationModuleNames []string
	inputModuleNames       []string
	state                  bool
}

type Action struct {
	signal      bool
	destination string
}

func (m *Module) processAction(action Action, modules map[string]Module) []Action {
	var actions []Action
	if m.moduleType == "broadcaster" {
		// Send signal to all destinations
		for _, destination := range m.destinationModuleNames {
			//println("Sending signal to", destination, " = ", action.signal, "from", m.moduleName)
			actions = append(actions, Action{
				signal:      action.signal,
				destination: destination,
			})
		}
		modules[m.moduleName] = *m
		return actions
	} else if m.moduleType == "%" {
		// Set state to flip if input it low
		if action.signal == false {
			m.state = !m.state
			// Send signal to all destinations
			for _, destination := range m.destinationModuleNames {
				//println("Sending signal to", destination, " = ", m.state, "from", m.moduleName)
				actions = append(actions, Action{
					signal:      m.state,
					destination: destination,
				})
			}
		}
		modules[m.moduleName] = *m
		return actions
	} else if m.moduleType == "&" {
		// Set state to true initially
		state := false

		// Check if all inputs are high, otherwise set state to false
		for _, input := range m.inputModuleNames {
			module := modules[input]
			if module.state == false {
				state = true
				break
			}
		}

		// Actually set the state
		m.state = state

		// Send signal to all destinations
		for _, destination := range m.destinationModuleNames {
			//println("Sending signal to", destination, " = ", m.state, "from", m.moduleName)
			actions = append(actions, Action{
				signal:      m.state,
				destination: destination,
			})
		}
		modules[m.moduleName] = *m
		return actions
	} else if m.moduleType == "output" {
		m.state = action.signal
		//println("Output:", m.state)
		modules[m.moduleName] = *m
		return actions
	} else {
		//do nothing
		return actions
	}
	panic("Unknown module type" + m.moduleType)
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
