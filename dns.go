package dns_toy

const (
	TYPE_A   uint16 = 1
	TYPE_NS  uint16 = 2
	CLASS_IN uint16 = 1

	RECURSION_DESIRED uint16 = 1 << 8
)
