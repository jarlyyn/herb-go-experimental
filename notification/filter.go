package notification

type Filter func(instance *MessageInstance, next func())

var FilterRecover = func(recovery func()) Filter {
	return func(instance *MessageInstance, next func()) {
		defer recovery()
		next()
	}
}

type RecipientConvertor func(recipient string) (string, error)

var FilterRecipient = func(convertor RecipientConvertor) Filter {
	return func(instance *MessageInstance, next func()) {
		recipient, err := instance.Message.MessageRecipient()
		if err != nil {
			panic(err)
		}
		id, err := convertor(recipient)
		if err != nil {
			panic(err)
		}
		err = instance.Message.SetMessageRecipient(id)
		if err != nil {
			panic(err)
		}
		next()
	}
}
