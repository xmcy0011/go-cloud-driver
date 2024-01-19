package common

// BasicAttr 对象类型
type ObjectType int

const (
	ObjectTypeDir  ObjectType = 1 // 文件
	ObjectTypeFile ObjectType = 2 // 目录
)
