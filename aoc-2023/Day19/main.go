package main

import (
	"fmt"
	"log"
	"os"
	"reflect"
	"strconv"
	"strings"
)

/*
*
x: Extremely cool looking
m: Musical (it makes a noise when you hit it)
a: Aerodynamic
s: Shiny
*/
type Workflow struct {
	id    string
	rules []Rule
}

type Rule struct {
	category     string
	operator     string //<, >
	limit        int
	nextWorkflow string //workflow id
}

type Rating struct {
	x int
	m int
	a int
	s int
}

func main() {
	input, err := readInput(2023, 19)
	if err != nil {
		fmt.Printf("Failed to read input file: %v\n", err)
		return
	}
	parts := strings.Split(input, "\n\n")
	workflows := parseWorkflows(parts[0])
	ratings := parseRatings(parts[1])

	println("workflows", len(workflows))
	println("ratings", len(ratings))

	accepted := []Rating{}
	rejected := []Rating{}

	unprocessed := append([]Rating{}, ratings...)
	for len(unprocessed) > 0 {
		rating := unprocessed[0]
		unprocessed = unprocessed[1:]

		workflow := workflows["in"]
		for workflow.id != "" {
			for _, rule := range workflow.rules {
				if rule.category == "FINAL RULE" {
					workflow = Workflow{}
					if rule.nextWorkflow == "A" {
						accepted = append(accepted, rating)
					} else if rule.nextWorkflow == "R" {
						rejected = append(rejected, rating)
					} else {
						workflow = workflows[rule.nextWorkflow]
					}
					break
				}
				fieldValue, err := getFieldValue(&rating, rule.category)
				if err != nil {
					panic(err)
				}
				if rule.operator == "<" && fieldValue < rule.limit {
					workflow = Workflow{}
					if rule.nextWorkflow == "A" {
						accepted = append(accepted, rating)
					} else if rule.nextWorkflow == "R" {
						rejected = append(rejected, rating)
					} else {
						workflow = workflows[rule.nextWorkflow]
					}
					break
				}
				if rule.operator == ">" && fieldValue > rule.limit {
					workflow = Workflow{}
					if rule.nextWorkflow == "A" {
						accepted = append(accepted, rating)
					} else if rule.nextWorkflow == "R" {
						rejected = append(rejected, rating)
					} else {
						workflow = workflows[rule.nextWorkflow]
					}
					break
				}
			}
		}
	}

	ansP1 := 0
	for _, r := range accepted {
		ansP1 += r.x + r.m + r.a + r.s
	}
	// Solution here
	fmt.Println("Solution:", ansP1)
}

func parseRatings(ratingLines string) []Rating {
	var ratings []Rating
	for _, ratingLine := range strings.Split(ratingLines, "\n") {
		var r Rating
		input := strings.Trim(ratingLine, "{}")
		pairs := strings.Split(input, ",")

		for _, pair := range pairs {
			kv := strings.Split(pair, "=")
			if len(kv) != 2 {
				panic("Invalid Format: " + pair)
			}

			key := kv[0]
			value := toInt(kv[1])

			switch key {
			case "x":
				r.x = value
			case "m":
				r.m = value
			case "a":
				r.a = value
			case "s":
				r.s = value
			default:
				panic("Invalid key: " + key)
			}
		}
		ratings = append(ratings, r)
	}
	return ratings
}

func parseWorkflows(workflowLines string) map[string]Workflow {
	workflows := make(map[string]Workflow)
	for _, wf := range strings.Split(workflowLines, "\n") {
		//px{a<2006:qkq,m>2090:A,rfg}
		id := strings.Split(wf, "{")[0]
		ruleSection := strings.Split(strings.Split(wf, "{")[1], "}")[0]
		ruleTokens := strings.Split(ruleSection, ",")
		rules := make([]Rule, 0)
		for _, rt := range ruleTokens {
			condition := strings.Split(rt, ":")[0]
			operatorIndex := strings.IndexAny(condition, "<>")
			if operatorIndex == -1 || operatorIndex == 0 || operatorIndex == len(condition)-1 {
				//panic("Invalid operator")
				// THIS IS A FINAL RULE
				rules = append(rules, Rule{
					category:     "FINAL RULE",
					operator:     "",
					limit:        0,
					nextWorkflow: condition,
				})
				continue
			}
			category := condition[operatorIndex-1 : operatorIndex]
			operator := condition[operatorIndex : operatorIndex+1]
			limit := condition[operatorIndex+1:]

			nextWorkflowId := strings.Split(rt, ":")[1]
			//fmt.Println(category, operator)
			println("id", id, "category", category, "operator", operator, "limit", limit, "nextWorkflowId", nextWorkflowId)
			rules = append(rules, Rule{
				category:     category,
				operator:     operator,
				limit:        toInt(limit),
				nextWorkflow: nextWorkflowId,
			})
		}
		workflows[id] = Workflow{
			id:    id,
			rules: rules,
		}
	}
	return workflows
}

func getFieldValue(r *Rating, field string) (int, error) {
	val := reflect.ValueOf(r).Elem()
	fieldValue := val.FieldByName(field)

	if !fieldValue.IsValid() {
		return 0, fmt.Errorf("no such field: %s in obj", field)
	}

	if !fieldValue.CanInt() {
		return 0, fmt.Errorf("cannot convert field: %s to int", field)
	}

	return int(fieldValue.Int()), nil
}

func toInt(s string) int {
	i, err := strconv.Atoi(s)
	if err != nil {
		panic(err)
	}
	return i
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
