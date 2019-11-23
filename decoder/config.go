package decoder

type Config struct {
	Checkers []TypeChecker
	Unifiers Unifiers
	Tagname  string
}
