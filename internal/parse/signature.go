package parse

import (
	"reflect"
	"regexp"
	"strings"
)

type (
	Validator struct {
		Definition string
		Validate   func(args []any, ctx any) ([]any, error)
	}

	Param struct {
		Regex        *regexp.Regexp
		Type         string
		Subtype      string
		Array        bool
		Context      bool
		ContextRegex *regexp.Regexp
	}
)

var (
	arraySignatureMapping = map[string]string{
		"a": "arrays",
		"b": "booleans",
		"f": "functions",
		"n": "numbers",
		"o": "objects",
		"s": "strings",
	}
)

// ParseSignature Parses a function signature definition and returns a validation function
func ParseSignature(signature string) (*Validator, error) {
	signatureRunes := []rune(signature)

	// create a Regex that represents this signature and return a function that when invoked,
	// returns the validated (possibly fixed-up) arguments, or throws a validation error
	// step through the signature, one symbol at a time
	var (
		position        = 1
		params          = []Param{}
		param           = Param{}
		prevParam Param = param
	)
	for position < len(signatureRunes) {
		symbol := signatureRunes[position]
		if symbol == ':' {
			// TODO figure out what to do with the return type
			// ignore it for now
			break
		}

		next := func() {
			params = append(params, param)
			prevParam = param
			param = Param{}
		}

		findClosingBracket := func(str []rune, start int, openSymbol rune, closeSymbol rune) int {
			// returns the position of the closing symbol (e.g. bracket) in a string
			// that balances the opening symbol at position start
			depth := 1
			position := start
			for position < len(str) {
				position++
				symbol := str[position]
				if symbol == closeSymbol {
					depth--
					if depth == 0 {
						// we're done
						break // out of while loop
					}
				} else if symbol == openSymbol {
					depth++
				}
			}
			return position
		}

		switch symbol {
		case 's', 'n', 'b', 'l', 'o': // string, number, boolean, not so sure about expecting null?, object
			if regex, err := regexp.Compile("[" + string(symbol) + "m]"); err != nil {
				return nil, err
			} else {
				param.Regex = regex
			}
			param.Type = string(symbol)
			next()
		case 'a': // array
			//  normally treat any value as singleton array
			if regex, err := regexp.Compile("[asnblfom]"); err != nil {
				return nil, err
			} else {
				param.Regex = regex
			}
			param.Type = string(symbol)
			param.Array = true
			next()
		case 'f': // function
			if regex, err := regexp.Compile("f"); err != nil {
				return nil, err
			} else {
				param.Regex = regex
			}
			param.Type = string(symbol)
			next()
		case 'j': // any JSON type
			if regex, err := regexp.Compile("[asnblom]"); err != nil {
				return nil, err
			} else {
				param.Regex = regex
			}
			param.Type = string(symbol)
			next()
		case 'x': // any type
			if regex, err := regexp.Compile("[asnblfom]"); err != nil {
				return nil, err
			} else {
				param.Regex = regex
			}
			param.Type = string(symbol)
			next()
		case '-': // use context if param not supplied
			prevParam.Context = true
			prevParam.ContextRegex = prevParam.Regex // pre-compiled to test the context type at runtime
			if regex, err := regexp.Compile(prevParam.Regex.String() + "?"); err != nil {
				return nil, err
			} else {
				prevParam.Regex = regex
			}
		case '?', '+': // optional param, one or more
			if regex, err := regexp.Compile(prevParam.Regex.String() + string(symbol)); err != nil {
				return nil, err
			} else {
				prevParam.Regex = regex
			}
		case '(': // choice of types
			// search forward for matching ')'
			endParen := findClosingBracket(signatureRunes, position, '(', ')')
			choice := signatureRunes[position+1 : endParen]
			if !strings.ContainsRune(string(choice), '<') {
				// no parameterized types, simple regex
				if regex, err := regexp.Compile("[" + string(choice) + "m]"); err != nil {
					return nil, err
				} else {
					param.Regex = regex
				}
			} else {
				// TODO harder
				return nil, &Error{
					Code:   "S0402",
					Value:  choice,
					Offset: position,
				}
			}
			param.Type = "(" + string(choice) + ")"
			position = endParen
			next()
		case '<': // type parameter - can only be applied to 'a' and 'f'
			if prevParam.Type == "a" || prevParam.Type == "f" {
				// search forward for matching '>'
				endPos := findClosingBracket(signatureRunes, position, '<', '>')
				prevParam.Subtype = string(signatureRunes[position+1 : endPos])
				position = endPos
			} else {
				return nil, &Error{
					Code:   "S0401",
					Value:  prevParam.Type,
					Offset: position,
				}
			}
		}
		position++
	}
	regexStrBuilder := &strings.Builder{}
	regexStrBuilder.WriteString("^")
	for _, param := range params {
		regexStrBuilder.WriteString("(" + param.Regex.String() + ")")
	}
	regexStrBuilder.WriteString("$")

	regex, err := regexp.Compile(regexStrBuilder.String())
	if err != nil {
		return nil, err
	}

	getSymbol := func(value any) rune {
		if IsNil(value) {
			return 'l'
		}

		t := reflect.TypeOf(value)
		for t.Kind() == reflect.Pointer {
			t = t.Elem()
		}

		switch t.Kind() {
		case reflect.String:
			return 's'
		case reflect.Int, reflect.Uint, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Float32, reflect.Float64:
			return 'n'
		case reflect.Bool:
			return 'b'
		case reflect.Func:
			return 'f'
		case reflect.Array, reflect.Slice:
			return 'a'
		case reflect.Map, reflect.Struct:
			return 'o'
		default:
			// any value can be undefined, but should be allowed to match
			return 'm' // m for missing
		}
	}

	throwValidationError := func(badArgs []any, badSig string) error {
		// to figure out where this went wrong we need apply each component of the
		// regex to each argument until we get to the one that fails to match
		partialPattern := "^"
		goodTo := 0
		for index := 0; index < len(params); index++ {
			partialPattern += params[index].Regex.String()
			partialRegex, err := regexp.Compile(partialPattern)
			if err != nil {
				return err
			}
			match := partialRegex.FindAllString(badSig, 1)
			if match == nil {
				// failed here
				return &Error{
					Code:  "T0410",
					Value: badArgs[goodTo],
					Index: goodTo + 1,
				}
			}
			goodTo = len([]rune(match[0]))
		}
		// if it got this far, it's probably because of extraneous arguments (we
		// haven't added the trailing '$' in the regex yet.
		return &Error{
			Code:  "T0410",
			Value: badArgs[goodTo],
			Index: goodTo + 1,
		}
	}

	return &Validator{
		Definition: signature,
		Validate: func(args []any, context any) ([]any, error) {
			suppliedSig := &strings.Builder{}
			for _, arg := range args {
				suppliedSig.WriteRune(getSymbol(arg))
			}
			if isValid := regex.FindAllString(suppliedSig.String(), -1); isValid != nil {
				var validatedArgs []any
				argIndex := 0
				for index, param := range params {
					arg := args[argIndex]
					if len(isValid) < index {
						if param.Context && param.ContextRegex != nil {
							// substitute context value for missing arg
							// first check that the context value is the right type
							contextType := getSymbol(context)
							// test contextType against the regex for this arg (without the trailing ?)
							if param.ContextRegex.MatchString(string(contextType)) {
								validatedArgs = append(validatedArgs, context)
							} else {
								// context value not compatible with this argument
								return nil, &Error{
									Code:  "T0411",
									Value: context,
									Index: argIndex + 1,
								}
							}
						} else {
							validatedArgs = append(validatedArgs, arg)
							argIndex++
						}
					} else {
						match := isValid[index+1]

						// may have matched multiple args (if the regex ends with a '+'
						// split into single tokens
						for _, single := range match {
							if param.Type == "a" {
								if single == 'm' {
									// missing (undefined)
									arg = nil
								} else {
									arg = args[argIndex]
									arrayOK := true
									// is there type information on the contents of the array?
									if param.Subtype != "" {
										if single != 'a' && match != param.Subtype {
											arrayOK = false
										} else if single == 'a' {
											var itemType rune
											var differentItems []any

											ForEach(arg, func(k, v any) bool {
												if itemType == 0 {
													itemType = getSymbol(v)
													if itemType != []rune(param.Subtype)[0] { // TODO recurse further
														arrayOK = false
														return false
													}
													return true
												}

												// make sure every item in the array is this type
												if getSymbol(v) != itemType {
													differentItems = append(differentItems, v)
												}

												return true
											})

											arrayOK = len(differentItems) == 0
										}
									}
									if !arrayOK {
										return nil, &Error{
											Code:  "T0412",
											Value: arg,
											Index: argIndex + 1,
											Type:  arraySignatureMapping[param.Subtype],
										}
									}
									// the function expects an array. If it's not one, make it so
									if single != 'a' {
										arg = []any{arg}
									}
								}
								validatedArgs = append(validatedArgs, arg)
								argIndex++
							} else {
								validatedArgs = append(validatedArgs, arg)
								argIndex++
							}
						}
					}
				}
				return validatedArgs, nil
			}
			return nil, throwValidationError(args, suppliedSig.String())
		},
	}, nil
}
