具体发送见msg.go的描述，消息的body中一般为ID，在传输aiid中为一个小端存储的id数组。

msg.go定义发送消息，match.go负责match相关部分的worker调用，stream.go定义了发送消息，transfer.go定义了具体发送和回传的消息，callscript.go负责调用赛事脚本部分。

请在match结束后，调用func CallOnMatchFinished(matchID uint, replay string) error函数回传。