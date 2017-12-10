package notification

const SendResultSuccess = 0
const SendResultTemplateNotFound = 1
const SendResultRecipientFormatIncorrect = 2

const QueryResultComplete = 0
const QueryResultPendding = 1
const QueryResultFail = 2
const QueryResultError = 3

type Sender interface {
	Send(message Message, recipient string) (result int, err error)
	Query(mesageID string) (result int, msg string, err error)
}
