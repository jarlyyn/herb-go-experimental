package notification

type ParamsIndex string

const ParamsIndexTitle = ParamsIndex("Title")
const ParamsIndexSiteName = ParamsIndex("SiteName")
const ParamsIndexHost = ParamsIndex("Host")
const ParamsIndexUserID = ParamsIndex("UserID")
const ParamsIndexUsername = ParamsIndex("Username")
const ParamsIndexName = ParamsIndex("Username")
const ParamsIndexObjectID = ParamsIndex("ObjectID")
const ParamsIndexObjectName = ParamsIndex("ObjectName")
const ParamsIndexObjectList = ParamsIndex("ObjectList")
const ParamsIndexTime = ParamsIndex("Time")

type Message struct {
	MessageID string
	Params    map[ParamsIndex][]string
}
