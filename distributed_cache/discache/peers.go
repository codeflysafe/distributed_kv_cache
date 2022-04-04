package discache

// 抽象出， 选择group 然后返回
type PeerGetter interface {
	Get(group string, key string) ([]byte, error)
}

// 根据传输过来的缓存键值，来获的对应的PeerGetter
type PeerPicker interface {
	PickPeer(key string) (peer PeerGetter, ok bool)
}

