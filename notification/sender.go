package notification

const SendResultSuccess = 0
const SendResultTemplateNotFound = 1
const SendResultRecipientIncorrect = 2

const QueryResultComplete = 0
const QueryResultPendding = 1
const QueryResultFail = 2
const QueryResultError = 3
const QueryResultNotFound = 4

type Sender interface {
	Send(message Message, templateid string, recipient string) (result int, err error)
	Query(mesageID string) (result int, msg string, err error)
}
