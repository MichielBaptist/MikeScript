package main

import (
	"bufio"
	"fmt"
	interp "mikescript/src/interp"
	parser "mikescript/src/parser"
	scanner "mikescript/src/scanner"
	"mikescript/src/utils"
	"os"
)

type Runner interface {
	run(input string) int
}

type MSRunner struct {
	prompter 	*bufio.Scanner			// prompter
	scanner 	scanner.MSScanner		// lexer/scanner
	parser 		parser.MSParser			// parser
	evaluator 	interp.MSEvaluator		// evaluator
	verbose 	bool
}

func (r MSRunner) log(input interface{}) {
	if r.verbose {
		fmt.Println(input)
	}
}

func (r MSRunner) run(input string) int {

	//////////////////////////////////////////////////////
	r.log("--------------- Scanner ---------------------")

	// call scanner
	tokens := r.scanner.Scan(input)

	// print the tokens
	r.log(fmt.Sprintf("Tokens (%v): %v", len(tokens), tokens))

	if len(r.scanner.Errors) > 0 {
		fmt.Printf("Scanner errors (%v): \n", len(r.scanner.Errors))
		for i, err := range r.scanner.Errors {
			fmt.Printf("[%v]: %v\n", i, err)
		}
		fmt.Println("")
		return 1
	}

	//////////////////////////////////////////////////////
	r.log("---------------- Parser ---------------------")

	// set the source code and tokens
	r.parser.SetSrc(input)
	r.parser.SetTokens(tokens)

	ast, _ := r.parser.Parse(tokens)

	r.log("AST:")
	r.log(ast)

	if len(r.parser.Errors) > 0 {
		fmt.Println("Parser errors:")
		for i, err := range r.parser.Errors {
			fmt.Printf("[%v]: %v\n", i, err)
		}
		fmt.Println("")
		return 1
	}

	//////////////////////////////////////////////////////
	r.log("--------------- Evaluator ---------------------")
	eval := r.evaluator.Eval(ast)

	// print the current env
	r.log("Environment:")
	if r.verbose {
		r.evaluator.PrintEnv()
	}

	// Print the result
	fmt.Println(&eval)

	return 0
}

func (r *MSRunner) prompt() bool {
	fmt.Print("ms> ")
	ok := r.prompter.Scan()
	return ok && 
		   r.prompter.Err() == nil &&
		   r.text() != "exit"
}

func (r *MSRunner) text() string {
	return r.prompter.Text()
}

func (r *MSRunner) mainLoop() {

	for r.prompt() {
		txt := r.text()
		r.run(txt)
	}
	fmt.Println("Goodbye!")
}

func main() {
	// create a new runner
	runner := MSRunner{
		prompter: 	bufio.NewScanner(os.Stdin),
		scanner: 	scanner.MSScanner{},
		parser: 	parser.MSParser{},
		evaluator: 	*interp.NewMSEvaluator(),
		verbose: 	true,
	}

	// Check if we have command line arguments
	if len(os.Args) > 1 {

		// read the file
		file, err := os.Open(os.Args[1])
		if err != nil {
			fmt.Println("Error reading file: ", err)
			return
		}
		defer file.Close()

		// read the entire file into a string
		scanner := bufio.NewScanner(file)
		lines := []string{}
		for scanner.Scan() {
			lines = append(lines, scanner.Text())
		}
		src := utils.StrJoin(lines, "\n")

		fmt.Println("#############################################")
		fmt.Println(src)
		fmt.Println("#############################################")

		// feed the file to the runner
		runner.run(src)

	} else {
		// No file provided, run the REPL
		fmt.Println("MikeScript 1.0 - REPL")
		fmt.Println("Type 'exit' to quit")
		runner.mainLoop()
	}
}