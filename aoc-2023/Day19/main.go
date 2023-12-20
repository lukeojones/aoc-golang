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

type RatingRange struct {
	destWorkflowId                                 string
	Xmin, Xmax, Mmin, Mmax, Amin, Amax, Smin, Smax int
}

func initialRatingRange() RatingRange {
	return RatingRange{
		destWorkflowId: "in",
		Xmin:           1,
		Xmax:           4000,
		Mmin:           1,
		Mmax:           4000,
		Amin:           1,
		Amax:           4000,
		Smin:           1,
		Smax:           4000,
	}
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

	ansP1 := solveP1(ratings, workflows)
	fmt.Println("Solution to Part 1:", ansP1)

	//Part 2
	workflows = parseWorkflows(parts[0])
	ansP2 := solveP2(workflows)

	println("Solution to Part 2", ansP2)
	if ansP2 != 125317461667458 {
		panic("Wrong answer")
	}
}

func solveP2(workflows map[string]Workflow) int {
	var accepted []RatingRange

	rr := initialRatingRange()
	unprocessed := []RatingRange{rr}
	for len(unprocessed) > 0 {
		ratingRange := unprocessed[0]
		unprocessed = unprocessed[1:]

		workflow := workflows[ratingRange.destWorkflowId]
		for workflow.id != "A" && workflow.id != "R" {
			for _, rule := range workflow.rules {
				if rule.category == "FINAL RULE" {
					workflow = workflows[rule.nextWorkflow]
					ratingRange.destWorkflowId = rule.nextWorkflow
					break
				}

				fieldMinValue, fieldMaxValue, _ := getFieldBounds(&ratingRange, strings.ToUpper(rule.category))
				if rule.operator == "<" {
					if fieldMinValue < rule.limit && fieldMaxValue < rule.limit {
						// Everything is in range, so pass the entire set to the next workflow
						rrn := ratingRange
						rrn.destWorkflowId = rule.nextWorkflow
						unprocessed = append(unprocessed, rrn)
					} else if fieldMinValue < rule.limit && fieldMaxValue > rule.limit {
						//Split the set into two - pass and fail

						// Passes go to the next workflow via unprocessed
						rrn1 := ratingRange
						setFieldValue(&rrn1, strings.ToUpper(rule.category)+"max", rule.limit-1)
						rrn1.destWorkflowId = rule.nextWorkflow
						unprocessed = append(unprocessed, rrn1)

						//Failures go to the next rule via staying loop
						rrn2 := ratingRange
						setFieldValue(&rrn2, strings.ToUpper(rule.category)+"min", rule.limit)
						ratingRange = rrn2
					} else if fieldMinValue > rule.limit && fieldMaxValue > rule.limit {
						//Everything is out of range, so pass the entire set to the next rule via ?
						//no op - just leave ratingRange as it is
					}
				} else if rule.operator == ">" {
					if fieldMinValue > rule.limit && fieldMaxValue > rule.limit {
						// Everything is in range, so pass the entire set to the next workflow
						rrn := ratingRange
						rrn.destWorkflowId = rule.nextWorkflow
						unprocessed = append(unprocessed, rrn)
					} else if fieldMinValue < rule.limit && fieldMaxValue > rule.limit {
						//Split the set into two - pass and fail

						// Passes go to the next workflow via unprocessed
						rrn1 := ratingRange
						setFieldValue(&rrn1, strings.ToUpper(rule.category)+"min", rule.limit+1)
						rrn1.destWorkflowId = rule.nextWorkflow
						unprocessed = append(unprocessed, rrn1)

						//Failures go to the next rule via staying loop
						rrn2 := ratingRange
						setFieldValue(&rrn2, strings.ToUpper(rule.category)+"max", rule.limit)
						ratingRange = rrn2
					} else if fieldMinValue > rule.limit && fieldMaxValue > rule.limit {
						//Everything is out of range, so pass the entire set to the next rule via ?
						//no op - just leave ratingRange as it is
					}
				}
			}
		}

		if workflow.id == "A" {
			accepted = append(accepted, ratingRange)
		}
	}

	ans := 0
	for _, ac := range accepted {
		ans += (ac.Xmax - ac.Xmin + 1) * (ac.Mmax - ac.Mmin + 1) * (ac.Amax - ac.Amin + 1) * (ac.Smax - ac.Smin + 1)
	}
	return ans
}

func solveP1(ratings []Rating, workflows map[string]Workflow) int {
	accepted := []Rating{}
	rejected := []Rating{}

	unprocessed := append([]Rating{}, ratings...)
	for len(unprocessed) > 0 {
		rating := unprocessed[0]
		unprocessed = unprocessed[1:]

		workflow := workflows["in"]
		for workflow.id != "A" && workflow.id != "R" {
			for _, rule := range workflow.rules {
				if rule.category == "FINAL RULE" {
					workflow = workflows[rule.nextWorkflow]
					break
				}
				fieldValue, err := getFieldValue(&rating, rule.category)
				if err != nil {
					panic(err)
				}
				if rule.operator == "<" && fieldValue < rule.limit {
					workflow = workflows[rule.nextWorkflow]
					break
				}
				if rule.operator == ">" && fieldValue > rule.limit {
					workflow = workflows[rule.nextWorkflow]
					break
				}
			}
		}
		if workflow.id == "A" {
			accepted = append(accepted, rating)
		} else if workflow.id == "R" {
			rejected = append(rejected, rating)
		}
	}

	ansP1 := 0
	for _, r := range accepted {
		ansP1 += r.x + r.m + r.a + r.s
	}

	return ansP1
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
		id := strings.Split(wf, "{")[0]
		ruleSection := strings.Split(strings.Split(wf, "{")[1], "}")[0]
		ruleTokens := strings.Split(ruleSection, ",")
		var rules []Rule
		for _, rt := range ruleTokens {
			condition := strings.Split(rt, ":")[0]
			operatorIndex := strings.IndexAny(condition, "<>")
			if operatorIndex == -1 || operatorIndex == 0 || operatorIndex == len(condition)-1 {
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

	workflows["A"] = Workflow{
		id:    "A",
		rules: []Rule{},
	}
	workflows["R"] = Workflow{
		id:    "R",
		rules: []Rule{},
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

func setFieldValue(r *RatingRange, field string, value int) error {
	v := reflect.ValueOf(r).Elem()
	fieldValue := v.FieldByName(field)

	if !fieldValue.IsValid() {
		return fmt.Errorf("no such field: %s in obj", field)
	}

	if !fieldValue.CanSet() {
		return fmt.Errorf("cannot set field: %s", field)
	}

	if fieldValue.Kind() != reflect.Int {
		return fmt.Errorf("field: %s is not an int", field)
	}

	fieldValue.SetInt(int64(value))
	return nil
}

func getFieldBounds(r *RatingRange, field string) (int, int, error) {
	val := reflect.ValueOf(r).Elem()
	fieldMinValue := val.FieldByName(field + "min")
	fieldMaxValue := val.FieldByName(field + "max")

	if !fieldMinValue.IsValid() {
		return 0, 0, fmt.Errorf("no such field: %s in obj", field)
	}

	if !fieldMaxValue.CanInt() {
		return 0, 0, fmt.Errorf("cannot convert field: %s to int", field)
	}

	return int(fieldMinValue.Int()), int(fieldMaxValue.Int()), nil
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
