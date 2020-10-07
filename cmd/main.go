package main

import (
	"fmt"
	"strconv"

	"github.com/manifoldco/promptui"
)

func main() {
	typ := selection("Type of the app", []string{
		"Web application (e.g. Node.js server, Golang server)",
		"Long-running worker",
		"AI/ML based application",
		"Other",
	})

	var cc int
	var tbc int
	switch typ {
	case 0:
	case 1:
		tbc = -1
		cc = promptInt("How many items should be handled in parallel by each worker instance?")
	case 2:
		cc = promptInt("How many items should be handled in parallel by each instance? Consider values < 5 if the model is very CPU/memory intensive")
	case 3:
		isResourceBound := selection("Is the application strictly constrained by resources like CPU and Memory?", []string{
			"No", "Yes",
		}) == 1
		if isResourceBound {
			cc = promptInt("How many items should be handled in parallel by each instance? Consider values < 10")
		}
	}

	var min int
	// Skip sensitivity question for long running workers.
	if typ != 1 {
		isSensitive := selection("Is the application very sensitive to latency, i.e. should it always respond in less than a second?", []string{
			"No", "Yes",
		}) == 1
		if isSensitive {
			min = promptInt("How many instances do you want to keep around to avoid the cold-start latency?")
		}
	}

	max := promptInt("Do you want to cap the amount of pods deployed for this service at tops (i.e. to avoid the service getting too expensive)?")

	avoidInitialScale := selection("Do you want to deploy a pod initially when creating the service to make sure it works appropriately?", []string{
		"Yes", "No",
	}) == 1

	fmt.Println()
	fmt.Println("Recommended configuration:")

	fmt.Println("containerConcurrency:", cc)
	if tbc != 0 {
		fmt.Println("autoscaling.knative.dev/targetBurstCapacity:", tbc)
	}
	if min != 0 {
		fmt.Println("autoscaling.knative.dev/minScale:", min)
	}
	if max != 0 {
		fmt.Println("autoscaling.knative.dev/maxScale:", max)
	}
	if avoidInitialScale {
		fmt.Println("autoscaling.knative.dev/initialScale: 0")
	}
}

func prompt(label string) string {
	p := promptui.Prompt{
		Label: label,
	}
	result, err := p.Run()
	if err != nil {
		panic(err)
	}
	return result
}

func promptInt(label string) int {
	resultInt, err := strconv.Atoi(prompt(label))
	if err != nil {
		panic(err)
	}
	return resultInt
}

func selection(label string, items []string) int {
	p := promptui.Select{
		Label: label,
		Items: items,
	}
	i, _, err := p.Run()
	if err != nil {
		panic(err)
	}
	return i
}
