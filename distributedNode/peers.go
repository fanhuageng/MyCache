package distributedNode

import pb "MyCache/cacheProtoBuf/mycachepb"

//// 根据传入的Key选择相应节点
//type PeerPicker interface {
//	PickPeer(key string) (peer PeerGetter, ok bool)
//}
//
//// 从对应group查找缓存值
//type PeerGetter interface {
//	Get(group string, key string) ([]byte, error)
//}

// 根据传入的Key选择相应节点
type PeerPicker interface {
	PickPeer(key string) (peer PeerGetter, ok bool)
}

// 从对应group查找缓存值
type PeerGetter interface {
	Get(in *pb.Request, out *pb.Response) error
}
