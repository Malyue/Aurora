package loadbalance

type LbType int

const (
	LbRandom LbType = iota
	LbRoundRobin
	LbWeightRoundRobin
	LbConsistentHash
)

//func LoadBalanceFactory(lbType LbType) LoadBalance {
//	switch lbType {
//	case LbRandom:
//		return &RandomBalance{}
//	case LbRoundRobin:
//		return &RoundRobinBalance{}
//	case LbWeightRoundRobin:
//		return &WeightRoundRobinBalance{}
//	case LbConsistentHash:
//		return NewConsistentHashBalance(10, nil)
//	default:
//		return &RandonBalance{}
//	}
//}
