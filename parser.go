package dolly

import (
	"errors"
	"fmt"
	"strings"
)

var (
	ErrMissingArguments = errors.New("missing arguments")
	ErrUnknownCommand   = errors.New("unknown command")
	ErrUnknownOptions   = errors.New("unknown options")
)

const (
	commentPrefix = "#"
	optionsPrefix = "@"
)

// Parse takes a string as input and returns the commands to be executed.
func Parse(s string) ([]Command, []error) {
	var commands []Command
	var errs []error

	lines := strings.Split(s, "\n")

	for i, line := range lines {
		lineNumber := i + 1

		if shouldSkip(line) {
			continue
		}

		valid := false
		for _, command := range CommandTypes {
			if strings.HasPrefix(line, command.String()) {
				valid = true
				options, args, err := parseArgs(command, line)
				if err != nil {
					errs = append(errs, fmt.Errorf("%s\n%d | %s", err, lineNumber, line))
					break
				}
				commands = append(commands, Command{command, options, args})
				break
			}
		}
		if !valid {
			errs = append(errs, fmt.Errorf("%s\n%d | %s", ErrUnknownCommand, lineNumber, line))
			continue
		}
	}

	return commands, errs
}

func parseArgs(command CommandType, line string) (string, string, error) {
	rawArgs := strings.TrimPrefix(line[len(command):], " ")

	// Set command
	// Set <Option> <Value>
	if command == Set {
		splitIndex := strings.Index(rawArgs, " ")
		if splitIndex == -1 {
			return "", "", ErrMissingArguments
		}

		options := rawArgs[:splitIndex]
		args := rawArgs[splitIndex+1:]
		_, ok := SetCommands[options]
		if !ok {
			return "", "", ErrUnknownOptions
		}

		return options, args, nil
	}

	// No @ options, return rawArgs as args
	if !strings.HasPrefix(rawArgs, optionsPrefix) {
		if command == Type && rawArgs == "" {
			return "", "", ErrMissingArguments
		}
		return "", rawArgs, nil
	}

	// Has @ options, parse options and arguments
	var options, args string
	splitIndex := strings.Index(rawArgs, " ")

	if splitIndex < 0 || splitIndex == len(rawArgs)-1 {
		return "", "", ErrMissingArguments
	}

	options = strings.TrimPrefix(rawArgs[:splitIndex], optionsPrefix)
	args = rawArgs[splitIndex+1:]
	return options, args, nil
}

func shouldSkip(line string) bool {
	return strings.HasPrefix(line, commentPrefix) || strings.TrimSpace(line) == ""
}
