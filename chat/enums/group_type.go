package enums

type GroupType int64

const (
	Group GroupType = iota + 1
	Personal
)

var GroupTypeMap = map[GroupType]string{
	Group:    "Group",
	Personal: "Personal",
}

func (gt GroupType) String() string {
	return GroupTypeMap[gt]
}
