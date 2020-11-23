package builtins

const abs_code string = `
cpy
psh 0
cmp
psh 1
sub
jmp :_positive_number@abs
psh -1
mul
:_positive_number@abs
`

const float_to_int_code string = `
cpy
psh 1
mod
sub
`

const str_to_int_code string = `
psh 0
ins _result@s2i

:_change@s2i
cpy
cpy

psh 48
cmp
psh 1
add
jmp :_end@s2i

psh 57
cmp
psh 1
sub
jmp :_end@s2i

psh 48
sub

psh 10
len
psh 2
sub
pow
mul

get _result@s2i
add
ins _result@s2i

len
psh 0
cmp
psh 1
sub
jmp :_change@s2i

:_end@s2i
clr
get _result@s2i

dlk _result@s2i
`

const round_code string = `
psh 0.5
add
mem
ins _line_n@round
cal float_to_int
get _line_n@round
pop

dlk _line_n@round
`

const lcg_code string = `
; (A * X + B) mod M
;A
psh 48271
mul
;B
psh 0
add
;M
psh 2
psh 31
pow
psh 1
sub
mod
mem
ins _line_n@lcg
cal float_to_int
get _line_n@lcg
pop

dlk _line_n@lcg
`

var builtin_funcs map[string]string = map[string]string{
	"abs":          abs_code,
	"float_to_int": float_to_int_code,
	"str_to_int":   str_to_int_code,
	"round":        round_code,
	"lcg":          lcg_code,
}

var builtin_funcs_arg map[string]int = map[string]int{
	"abs":          1,
	"float_to_int": 1,
	"str_to_int":   1,
	"round":        1,
	"lcg":          1,
}

func Builtin_func() string {
	var funcs_code string = "psh 0\njmp :_funcs_end\n"

	for k := range builtin_funcs {
		funcs_code += ":" + k + "\npop\n" + builtin_funcs[k] + "\npsh 0\nmem\njmp\n"
	}

	return funcs_code + ":_funcs_end\n"
}

func Builtin_func_arg_info() map[string]int {
	return builtin_funcs_arg
}
