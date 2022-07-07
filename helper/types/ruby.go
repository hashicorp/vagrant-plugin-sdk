package types

type RawRubyValue struct {
	Source Class // Ruby source class
	Data   map[string]interface{}
}
