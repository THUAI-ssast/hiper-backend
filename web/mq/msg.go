package mq

// Msg 消息
type Msg struct {
	ID        string // 消息的编号
	Topic     string // 消息的主题(Build/Run/Change)
	Body      []byte // 消息的Body
	Partition int    // 分区号
	Type      string // 消息类型(AI/Match/Game/SDK)
}
