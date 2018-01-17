package query

func conncatWith(separator string, q ...Query) *PlainQuery {
	var query = New("")
	for k := range q {
		if q == nil {
			continue
		}
		command := q[k].QueryCommand()
		if command != "" {
			query.Command += q[k].QueryCommand() + separator
		}
		query.Args = append(query.Args, q[k].QueryArgs()...)
	}
	if query.Command != "" {
		query.Command = query.Command[:len(query.Command)-len(separator)]
	}
	return query
}
func Concat(q ...Query) *PlainQuery {
	return conncatWith(" ", q...)
}

func Comma(q ...Query) *PlainQuery {
	return conncatWith(" , ", q...)
}
func Lines(q ...Query) *PlainQuery {
	return conncatWith("\n", q...)
}
func And(q ...Query) *PlainQuery {
	var query = conncatWith(" AND ", q...)
	if query.Command != "" {
		query.Command = "( " + query.Command + " )"
	}
	return query
}

func Or(q ...Query) *PlainQuery {
	var query = conncatWith(" OR ", q...)
	if query.Command != "" {
		query.Command = "( " + query.Command + " )"
	}
	return query
}
