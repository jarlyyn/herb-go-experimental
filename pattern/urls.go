package pattern

type FixedURLs map[string]bool

type FixedURLWithHost map[string]FixedURLs

type URLPrefixs []string

type URLPrefixsWithHost map[string]URLPrefixs
