package sgarg

import (
	"errors"
	"strings"
	"unicode/utf8"
)

var (
	ErrOptNotFound      = errors.New("option not found")
	ErrOptAlreadyExists = errors.New("option already exists")
	ErrOptFmtInvalid    = errors.New("option's format invalid")
	ErrOptNameAmbiguous = errors.New("option's name ambiguous")
)

type Parser struct {
	opts       map[optName]*opt
	nonOptArgs []string
}

func NewParser() *Parser {
	return &Parser{
		opts: make(map[optName]*opt),
	}
}

func (p *Parser) NonOptArgs() []string {
	return p.nonOptArgs
}

func (p *Parser) findLongOpt(name string) (*opt, error) {
	results := make([]*opt, 0, len(p.opts))
	for _, v := range p.opts {
		if v.abbreviatable(name) {
			results = append(results, v)
		}
	}
	switch len(results) {
	case 0:
		return nil, ErrOptNotFound
	case 1:
		return results[0], nil
	default:
		return nil, ErrOptNameAmbiguous
	}
}

func (p *Parser) setOpt(name string, values optValues) error {
	opt, err := newOpt(name, values)
	if err != nil {
		return err
	}
	if _, ok := p.opts[opt.name]; ok {
		return ErrOptAlreadyExists
	}
	p.opts[opt.name] = opt
	return nil
}

func (p *Parser) SetBoolOpt(name string, values *[]bool) error {
	return p.setOpt(name, &boolOptValues{values})
}

func (p *Parser) SetStringOpt(name string, values *[]string) error {
	return p.setOpt(name, &stringOptValues{values})
}

func (p *Parser) Parse(args []string) error {
	for idx := 0; idx < len(args); {
		switch p.argType(args[idx]) {
		case nonOptArg:
			p.setNonOptArgs(args[idx:])
			return nil
		case shortOpt:
			count, err := p.parseShortOpt(idx, args)
			if err != nil {
				return err
			}
			idx += count
		case longOpt:
			err := p.parseLongOpt(idx, args)
			if err != nil {
				return err
			}
			idx += 1
		case optTerminater:
			if idx != len(args)-1 {
				p.setNonOptArgs(args[idx+1:])
			}
			return nil
		}
	}
	return nil
}

func (p *Parser) parseShortOpt(idx int, args []string) (int, error) {
	rs := []rune(args[idx])
	parsed := 1
	for rIdx := 1; rIdx < len(rs); {
		opt, ok := p.opts[optName(rs[rIdx])]
		if !ok {
			return parsed, ErrOptNotFound
		}
		if !opt.withArg() {
			if err := opt.appendValue("true"); err != nil {
				return parsed, err
			}
			rIdx++
			continue
		}

		var value string
		if rIdx != len(rs)-1 {
			value = string(rs[rIdx+1:])
		} else if idx != len(args)-1 && p.argType(args[idx+1]) == nonOptArg {
			value = args[idx+1]
			parsed++
		} else {
			return parsed, ErrOptFmtInvalid
		}
		if err := opt.appendValue(value); err != nil {
			return parsed, err
		}
		break
	}
	return parsed, nil
}

func (p *Parser) parseLongOpt(idx int, args []string) error {
	keyVal := strings.Split(args[idx][2:], "=")
	key := keyVal[0]
	opt, err := p.findLongOpt(key)
	if err != nil {
		return err
	}
	switch len(keyVal) {
	case 1:
		if opt.withArg() {
			return ErrOptFmtInvalid
		}
		return opt.appendValue("true")
	case 2:
		val := keyVal[1]
		if !opt.withArg() {
			return ErrOptFmtInvalid
		}
		return opt.appendValue(val)
	default:
		return ErrOptFmtInvalid
	}
}

type argType int

const (
	nonOptArg argType = iota
	shortOpt
	longOpt
	optTerminater
)

func (p *Parser) argType(arg string) argType {
	if utf8.RuneCountInString(arg) < 2 || !strings.HasPrefix(arg, "-") {
		return nonOptArg
	}
	if !strings.HasPrefix(arg, "--") {
		return shortOpt
	}
	if arg == "--" {
		return optTerminater
	}
	return longOpt
}

func (p *Parser) setNonOptArgs(args []string) {
	nonOptArgs := make([]string, len(args))
	copy(nonOptArgs, args)
	p.nonOptArgs = nonOptArgs
}
