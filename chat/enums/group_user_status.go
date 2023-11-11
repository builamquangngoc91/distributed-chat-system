package enums

type GroupUserStatus int64

const (
	Joined GroupUserStatus = iota + 1
	Leaved
)

var GroupUserStatusMap = map[GroupUserStatus]string{
	Joined: "Joined",
	Leaved: "Leaved",
}

func (gst GroupUserStatus) String() string {
	return GroupUserStatusMap[gst]
}
