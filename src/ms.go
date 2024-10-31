package main

import (
	"bufio"
	"fmt"
	"os"
)

type Runner interface {
	run(input string) int
}

type MSRunner struct {
	prompter *bufio.Scanner	// prompter
	scanner MSScanner		// lexer/scanner
	parser MSParser			// parser
	evaluator MSEvaluator	// evaluator
}

func (r MSRunner) run(input string) int {

	fmt.Println("--------------- Scanner ---------------------")
	// call scanner
	tokens := r.scanner.Scan(input)

	fmt.Printf("Tokens (%v): %v\n", len(tokens), tokens)
	fmt.Printf("Errors (%v): \n", len(r.scanner.errors))
	for i, err := range r.scanner.errors {
		fmt.Printf("[%v]: %v\n", i, err)
	}
	fmt.Println("")

	fmt.Println("--------------- Parser ---------------------")

	// set the source code and tokens
	r.parser.src = input
	r.parser.tokens = tokens

	// call parser
	ast, _ := r.parser.parse(tokens)

	fmt.Println("AST:")
	fmt.Println(ast)
	fmt.Println("Errors:")
	for i, err := range r.parser.Errors {
		fmt.Printf("[%v]: %v\n", i, err)
	}
	fmt.Println("")

	// Check for errors, if any return
	if len(r.parser.Errors) > 0 {
		return 1
	}

	// evaluate
	fmt.Println("--------------- Evaluator ---------------------")
	eval := r.evaluator.Eval(ast)
	fmt.Println("Eval:")
	fmt.Println(eval)

	// print the current env
	fmt.Println("Environment:")
	r.evaluator.env.printEnv()

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
		r.run(r.text())
	}
	fmt.Println("Goodbye!")
}

func main() {
	// create a new runner
	runner := MSRunner{
		prompter: bufio.NewScanner(os.Stdin),
		scanner: MSScanner{},
		parser: MSParser{},
		evaluator: *NewMSEvaluator(),
	}
	runner.mainLoop()
}