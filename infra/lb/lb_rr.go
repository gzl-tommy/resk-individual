package lb

import "sync/atomic"

var _ Balancer = new(RoundRobinBalancer)

// round robin:简单轮算法
type RoundRobinBalancer struct {
	ct uint32 // 计数器
}

func (r *RoundRobinBalancer) Next(key string, hosts []*ServerInstance) *ServerInstance {
	if len(hosts) == 0 {
		return nil
	}

	//自增
	count := atomic.AddUint32(&r.ct, 1)
	//取模计算索引
	index := int(count) % len(hosts)
	//按照索引取出实例
	instance := hosts[index]
	return instance
}
