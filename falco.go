package main

import (
	"./builtins"
	"bufio"
	"flag"
	"fmt"
	"io/ioutil"
	"math"
	"os"
	"strconv"
	"strings"
)

func push_stack(stack *[]float64, value float64) {
	*stack = append(*stack, value)
}

func pop_stack(stack *[]float64) (float64, bool) {
	if len(*stack) == 0 {
		return -1, false
	}

	value, ok := return_stack_top_value(*stack)

	if ok {
		*stack = (*stack)[0 : len(*stack)-1]
	}
	return value, ok
}

func return_stack_top_value(stack []float64) (float64, bool) {
	if len(stack) == 0 {
		return -1, false
	}
	return stack[len(stack)-1], true
}

func find_label(code_list []string, label string) int {
	for i := 0; i < len(code_list); i++ {
		if code_list[i] == label {
			return i
		}
	}

	return -1
}

func falco(code string, debug bool) {
	var code_list []string

	code_list = strings.Split(builtins.Builtin_func()+code, "\n")

	var stacks [][]float64
	var memory float64 = 0
	var stack_n int = 0
	var line_n int = 0

	var accuracy float64 = math.Pow(10, 64)

	stacks = append(stacks, []float64{})

	dict := make(map[string]float64)

	var arg string

	var builtin_func map[string]int = builtins.Builtin_func_arg_info()

	for line_n < len(code_list) {
		line := code_list[line_n]
		line_n++

		if line == "" {
			continue
		}
		if line[0] == ':' || line[0] == ';' {
			continue
		}

		order := (line[:3])

		if len(line) > 3 {
			arg = (line[4:])
		} else {
			arg = ""
		}

		if debug {
			fmt.Println()
			fmt.Println(strings.Repeat("-", 50))
			fmt.Println("|stacks        :", stacks)
			fmt.Println("|stack number  :", stack_n)
			fmt.Println("|stack         :", stacks[stack_n])
			fmt.Println("|map           :", dict)
			fmt.Println("|memory        :", memory)
			fmt.Println("|")
			fmt.Println("|order         :", order)
			fmt.Println("|argument      :", arg)
			fmt.Println(strings.Repeat("-", 50))
		}

		switch order {
		case "psh":
			value, err := strconv.ParseFloat(arg, 64)
			if err == nil {
				push_stack(&stacks[stack_n], value)
			} else {
				if arg[0] == '"' && arg[len(arg)-1] == '"' {
					arg = arg[1 : len(arg)-1]
				}
				str := []rune(arg)

				for i := len(str) - 1; i >= 0; i-- {
					push_stack(&stacks[stack_n], float64(str[i]))
				}
			}

		case "pop":
			value, ok := pop_stack(&stacks[stack_n])

			if ok {
				memory = value
			} else {
				fmt.Println("Error")
			}

		case "del":
			pop_stack(&stacks[stack_n])

		case "mem":
			push_stack(&stacks[stack_n], memory)

		case "mks":
			stacks = append(stacks, []float64{})

		case "mov":
			value, ok := pop_stack(&stacks[stack_n])

			if ok {
				if arg == "" {
					n, ok := pop_stack(&stacks[stack_n])
					if ok {
						push_stack(&stacks[int(n)], value)
					}
				} else {
					n, err := strconv.Atoi(arg)
					if err == nil {
						push_stack(&stacks[n], value)
					}
				}
			}

		case "cng":
			if arg == "" {
				if len(stacks[stack_n]) > 0 {
					n, ok := pop_stack(&stacks[stack_n])
					if ok {
						stack_n = int(n)
					} else {
						fmt.Println("Error")
					}
				} else {
					fmt.Println("Error")
				}
				continue
			}

			n, err := strconv.Atoi(arg)

			if err == nil {
				if n < len(stacks) {
					stack_n = n
				}
			}

		case "cpy":
			value, ok := return_stack_top_value(stacks[stack_n])
			if ok {
				push_stack(&stacks[stack_n], value)
			}

		case "len":
			l := len(stacks[stack_n])

			push_stack(&stacks[stack_n], float64(l))

		case "clr":
			stacks[stack_n] = []float64{}

		case "stn":
			push_stack(&stacks[stack_n], float64(stack_n))

		case "nst":
			push_stack(&stacks[stack_n], float64(len(stacks)))

		case "lnm":
			push_stack(&stacks[stack_n], float64(line_n-1))

		case "jmp":
			if len(stacks[stack_n]) <= 0 {
				fmt.Println("Error")
				continue
			}

			var setting_n int = 0

			if arg == "" {
				n, ok := pop_stack(&stacks[stack_n])
				if ok {
					setting_n = int(n)
				}
			} else {
				if arg[0] == ':' {
					n := find_label(code_list, arg)

					if n < 0 {
						fmt.Println("Error")
						continue
					} else {
						setting_n = n
					}

				} else {
					n, err := strconv.Atoi(arg)

					if err != nil {
						fmt.Println("Error")
						continue
					} else {
						setting_n = n
					}

				}
			}

			b, ok := pop_stack(&stacks[stack_n])
			if ok {
				if b == 0 {
					line_n = setting_n
				}
			} else {
				fmt.Println("Error")
				continue
			}

		case "cal":
			if arg == "" {
				fmt.Println("Error")
				continue
			}

			number_of_args, ok := builtin_func[arg]

			if !ok {
				fmt.Println("Error")
				continue
			}
			if len(stacks[stack_n]) < number_of_args {
				fmt.Println("Error")
				continue
			}

			push_stack(&stacks[stack_n], float64(line_n))
			line_n = find_label(code_list, ":"+arg)

		case "ins":
			if len(stacks[stack_n]) <= 0 {
				fmt.Println("Error")
				continue
			}

			name := arg
			value, ok := pop_stack(&stacks[stack_n])
			if ok {
				dict[name] = value
			}

		case "get":
			name := arg

			value, ok := dict[name]

			if ok {
				push_stack(&stacks[stack_n], value)
			} else {
				fmt.Println("Error")
			}

		case "dlk":
			name := arg

			_, ok := dict[name]

			if ok {
				delete(dict, name)
			} else {
				fmt.Println("Error")
			}

		case "inp":
			stdin := bufio.NewScanner(os.Stdin)
			stdin.Scan()
			text := stdin.Text()

			str := []rune(text)

			for i := len(str) - 1; i >= 0; i-- {
				push_stack(&stacks[stack_n], float64(str[i]))
			}

		case "chr":
			if len(stacks[stack_n]) <= 0 {
				fmt.Println("Error")
				continue
			}

			value, ok := pop_stack(&stacks[stack_n])
			if ok {
				fmt.Print(string(int(value)))
			}

		case "prt":
			if len(stacks[stack_n]) <= 0 {
				fmt.Println("Error")
				continue
			}

			value, ok := pop_stack(&stacks[stack_n])

			if ok {
				fmt.Print(value)
			}

		case "cmp", "add", "sub", "mul", "div", "pow", "mod":
			if len(stacks[stack_n]) < 2 {
				fmt.Println("Error")
				continue
			}

			right, ok := pop_stack(&stacks[stack_n])
			left, ok := pop_stack(&stacks[stack_n])

			if ok {
				var value float64

				switch order {
				case "cmp":
					if left == right {
						value = 0
					} else {
						if left < right {
							value = -1
						} else {
							value = 1
						}
					}

				case "add":
					value = left + right
				case "sub":
					value = left - right
				case "mul":
					value = left * right
				case "div":
					value = left / right
				case "pow":
					value = math.Pow(left, right)
				case "mod":
					value = ((left / right) - float64(int(left/right))) * right
					value = math.Round(value*accuracy) / accuracy

				default:
					value = 0
				}

				push_stack(&stacks[stack_n], value)
			}

		default:
			fmt.Println("Error")
		}

	}
	if debug {
		fmt.Println("\n|END|")
		fmt.Println(strings.Repeat("*", 50))
		fmt.Println("|stacks  :", stacks)
		fmt.Println("|map     :", dict)
		fmt.Println("|memory  :", memory)
		fmt.Println(strings.Repeat("*", 50))
	}

}

func main() {
	flag.Parse()
	args := flag.Args()

	var debug bool = false

	if args[0] == "debug" {
		debug = true
	}

	name := args[len(args)-1]

	if len(name) <= 3 {
		fmt.Println("Error")
		return
	}

	if name[len(name)-4:] != ".flc" {
		fmt.Println("Error")
		return
	}

	f, err := os.Open(name)
	if err != nil {
		fmt.Println("Error")
		return
	}
	defer f.Close()

	// 一気に全部読み取り
	b, err := ioutil.ReadAll(f)
	if err != nil {
		fmt.Println("Error")
		return
	}

	falco(string(b), debug)
}
