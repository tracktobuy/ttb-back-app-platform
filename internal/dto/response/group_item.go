package response

type GroupItem struct {
	GroupUUID string `bson:"groupUUID" json:"groupUUID"`
	GroupName string `bson:"groupName" json:"groupName"`
	Items     []Item `bson:"items" json:"items"`
}
