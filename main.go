package main

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/urfave/cli/v2"
	"os"
	"path"
	"strings"
)

const (
	TAB   = 9
	LINE  = 10
	SPACE = 32
	COMMA = 44
	DOT   = 46
)

func main() {
	app := &cli.App{
		Name:   path.Base(os.Args[0]),
		Usage:  "input stream to bracket",
		Action: Action,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "mode",
				Aliases: []string{"M"},
				Usage:   "separate mode line/comma/dot/tab/space",
				Value:   "line",
			},
			&cli.StringFlag{
				Name:    "format",
				Aliases: []string{"F"},
				Usage:   "bracket format printf",
				Value:   `"%s"`,
			},
			&cli.StringFlag{
				Name:    "separate",
				Aliases: []string{"S"},
				Usage:   "separate output args",
				Value:   "\n",
			},
		},
	}
	err := app.Run(os.Args)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	os.Exit(0)
}

func Action(context *cli.Context) error {
	scanner := bufio.NewScanner(os.Stdin)
	switch context.String("mode") {
	case "line":
		scanner.Split(ScanWords(LINE))
		break
	case "comma":
		scanner.Split(ScanWords(COMMA))
		break
	case "dot":
		scanner.Split(ScanWords(DOT))
		break
	case "tab":
		scanner.Split(ScanWords(TAB))
		break
	case "space":
		scanner.Split(ScanWords(SPACE))
		break
	default:
		return errors.New("unsupport separate mode: " + context.String("mode"))
	}

	for scanner.Scan() {
		text := TrimAll(scanner.Text())
		if text == "" {
			continue
		}
		fmt.Printf("%s%s", fmt.Sprintf(context.String("format"), text), context.String("separate"))
	}
	if err := scanner.Err(); err != nil {
		return err
	}
	return nil
}

func TrimAll(line string) string {
	line = strings.TrimRight(line, "\n")
	if strings.HasSuffix(line, "\r") {
		line = strings.TrimRight(line, "\r")
	}
	line = strings.TrimLeft(line, "\"")
	line = strings.TrimRight(line, "\"")
	line = strings.TrimLeft(line, "'")
	line = strings.TrimRight(line, "'")
	return line
}

func ScanWords(separate byte) func(data []byte, atEOF bool) (advance int, token []byte, err error) {
	return func(data []byte, atEOF bool) (advance int, token []byte, err error) {
		for i := 0; i < len(data); i++ {
			if data[i] == separate {
				return i + 1, data[:i], nil
			}
		}
		return 0, data, bufio.ErrFinalToken
	}
}
