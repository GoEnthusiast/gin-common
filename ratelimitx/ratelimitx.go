package ratelimitx

import (
	"github.com/juju/ratelimit"
	"net/http"
	"strings"
	"time"
)

type LimiterIface interface {
	Key(r *http.Request) string                         // 获取对应的限流器的键值对名称
	GetBucket(key string) (*ratelimit.Bucket, bool)     // 获取令牌桶
	AddBuckets(rules ...LimiterBucketRule) LimiterIface // 新增多个令牌桶
}

type Limiter struct {
	limiterBuckets map[string]*ratelimit.Bucket
}

// LimiterBucketRule
/*
	Key:          "/api",           // 自定义键值对名称
	FillInterval: 60 * time.Second, // 间隔多久时间释放一次令牌桶
	Capacity:     150,              // 令牌捅的容量
	Quantum:      150,              // 每次到达间隔时间后所释放的具体令牌数量
*/
type LimiterBucketRule struct {
	Key          string
	FillInterval time.Duration
	Capacity     int64
	Quantum      int64
}

type MethodLimiter struct {
	*Limiter
}

func (l *MethodLimiter) Key(r *http.Request) string {
	uri := r.RequestURI
	index := strings.Index(uri, "?")
	if index != -1 {
		uri = uri[:index]
	}
	for k, _ := range l.Limiter.limiterBuckets {
		if strings.HasPrefix(uri, k) {
			return k
		}
	}
	return uri
}

func (l *MethodLimiter) GetBucket(key string) (*ratelimit.Bucket, bool) {
	bucket, ok := l.limiterBuckets[key]
	return bucket, ok
}

func (l *MethodLimiter) AddBuckets(rules ...LimiterBucketRule) LimiterIface {
	for _, rule := range rules {
		if _, ok := l.limiterBuckets[rule.Key]; !ok {
			bucket := ratelimit.NewBucketWithQuantum(
				rule.FillInterval,
				rule.Capacity,
				rule.Quantum,
			)
			l.limiterBuckets[rule.Key] = bucket
		}
	}
	return l
}

// NewMethodLimiter
/*
	Key:          "/api",           // 自定义加入该限流的路径
	FillInterval: 60 * time.Second, // 间隔多久时间释放一次令牌桶
	Capacity:     150,              // 令牌捅的容量
	Quantum:      150,              // 每次到达间隔时间后所释放的具体令牌数量
*/
func NewMethodLimiter(limiterInfos []LimiterInfo) LimiterIface {
	l := &Limiter{limiterBuckets: make(map[string]*ratelimit.Bucket)}
	lim := &MethodLimiter{Limiter: l}

	var LimiterBucketRules []LimiterBucketRule
	for _, limiterInfo := range limiterInfos {
		LimiterBucketRules = append(LimiterBucketRules, LimiterBucketRule{
			Key:          limiterInfo.Key,
			FillInterval: time.Duration(limiterInfo.FillInterval) * time.Second,
			Capacity:     limiterInfo.Capacity,
			Quantum:      limiterInfo.Quantum,
		})
	}

	lim.AddBuckets(LimiterBucketRules...)
	return lim
}

type LimiterInfo struct {
	Key          string `json:"key"`
	FillInterval int64  `json:"fillInterval"`
	Capacity     int64  `json:"capacity"`
	Quantum      int64  `json:"quantum"`
}
