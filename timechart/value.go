package timechart

type Value struct {
	DatetimeFirst int64
	DatetimeLast  int64
	FirstValue    float64
	LastValue     float64
	MinValue      float64
	MaxValue      float64
	AvgValue      float64
	CountOfValues int
	Qualities     []int64
	Loaded        bool
	UOM           string
}

func (c *Value) hasGoodQuality() bool {
	for _, q := range c.Qualities {
		if q == 192 {
			return true
		}
	}
	return false
}

func (c *Value) hasBadQuality() bool {
	for _, q := range c.Qualities {
		if q == 0 {
			return true
		}
	}
	return false
}
