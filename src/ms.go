package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	interp "mikescript/src/interp"
	parser "mikescript/src/parser"
	scanner "mikescript/src/scanner"
	"os"
	"strings"
)

// Repl commands
type Command uint8
const (
	EXIT Command = iota
	LOAD
	RUN
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

func (r MSRunner) log(input any) {
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

func (r *MSRunner) isExit(s string) bool {
	// split string and get first word
	strs := strings.Split(s, " ")

	if len(strs) == 0 {
		return false
	}

	return strs[0] == "exit"
}

func (r *MSRunner) isLoad(s string) (bool, string) {
	strs := strings.Split(s, " ")

	fmt.Println(strs)

	// Only one string, so not a command
	if len(strs) == 0 {
		return false, s
	}

	// check for load
	if strs[0] == "load" {
		return true, strings.Join(strs[1:], " ")
	}

	return false, strings.Join(strs[1:], " ")
}

func (r *MSRunner) prompt() (Command, string) {
	fmt.Print("ms> ")
	ok := r.prompter.Scan()

	// something went wrong??
	if !ok {
		fmt.Println("Something went wrong during prompting...")
		return EXIT, "exit"
	}

	// Get text
	txt := r.text()

	// Check if the text is exit, if it is we return exit
	if r.isExit(txt) {
		return EXIT, txt
	}

	// Check for load directive
	if b, txt := r.isLoad(txt); b {
		return LOAD, txt
	}

	// Not anything, normal command?
	return RUN, txt
}

func (r *MSRunner) text() string {
	return r.prompter.Text()
}

func (r *MSRunner) loadCommand(txt string) {
	// Split the txt into paths
	paths := strings.Split(txt, " ")

	// loop all paths and load the src file
	for _, p := range paths {

		src, err := readMSFile(p)

		// Something went wrong with reading the file.
		if err != nil {
			log.Println("Failed to load MikeScript file: ", p)
			continue
		}

		// Print source
		fmt.Println("#############################################")
		fmt.Println(src)
		fmt.Println("#############################################")

		// Loaded source correctly, so we exec.
		r.run(src)
	}
}

func (r *MSRunner) mainLoop() {

	promptloop: for {
		switch cmd, txt := r.prompt(); cmd {
		case LOAD: 	r.loadCommand(txt)
		case RUN:	r.run(txt)
		case EXIT: 	break promptloop
		}
	}
	fmt.Println("Goodbye!")
}

func readMSFile(path string) (string, error) {
	handle, err := os.Open(path)

	if err != nil {
		log.Println("Could not load file:", path)
		return "", err
	}

	// Close file on exit
	defer handle.Close()

	b, err := io.ReadAll(handle)
	
	if err != nil {
		return "", err
	}

	return string(b), nil
	
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
		src := strings.Join(lines, "\n")

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