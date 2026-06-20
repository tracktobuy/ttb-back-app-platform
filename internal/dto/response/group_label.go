package response

type GroupLabel struct {
	GroupUUID string   `bson:"groupUUID" json:"groupUUID"`
	GroupName string   `bson:"groupName" json:"groupName"`
	Labels    []string `bson:"labels" json:"labels"`
}
