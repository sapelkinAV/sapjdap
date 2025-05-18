package client

/*
Id sizes in bytes
(this object should be initialized on start)
*/
type ObjectIdSizes struct {
	FieldIdSize     int
	MethodIdSize    int
	ObjectIdSize    int
	ReferenceIdSize int
	FrameIdSize     int
}

var OBJECT_ID_SIZES *ObjectIdSizes = nil
