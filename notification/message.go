package notification

type ParamsIndex string

const ParamsIndexTitle = ParamsIndex("Title")
const ParamsIndexSitnName = ParamsIndex("SiteName")
const ParamsIndexHost = ParamsIndex("Host")
const ParamsIndexUserID = ParamsIndex("UserID")
const ParamsIndexUsername = ParamsIndex("Username")
const ParamsIndexObjectID = ParamsIndex("ObjectID")
const ParamsIndexObjectName = ParamsIndex("ObjectName")

type Message struct {
	MessageID  string
	TemplateID string
	Params     map[ParamsIndex][]string
}
