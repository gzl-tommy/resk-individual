package lb

import "math/rand"

var _ Balancer = new(RandomBalancer)

// 随机负载均衡算法
type RandomBalancer struct {
}

func (r *RandomBalancer) Next(key string, hosts []*ServerInstance) *ServerInstance {
	if len(hosts) == 0 {
		return nil
	}
	//随机数
	count := rand.Uint32()
	//取模计算索引
	index := int(count) % len(hosts)
	//按照索引取出实例
	instance := hosts[index]
	return instance
}