package uom

const (
	EntityInformation = "information"
	BYTES             = "bytes"
	KB                = "KB"
	MB                = "MB"
	GB                = "GB"
	TB                = "TB"
)

const (
	EntityTime = "time"
	MS         = "ms"
	SEC        = "sec"
	MIN        = "min"
	HOUR       = "hour"
	DAY        = "day"
)

const (
	ERROR    = "error"
	EVENT    = "-"
	STARTED  = "started"
	STOPPED  = "stopped"
	NONE     = ""
	PERCENTS = "%"
)

type UOM struct {
	Name   string
	Entity string
	IsBase bool
	A      float64
	B      float64
}

type Provider struct {
	unitsOfMeasure map[string]*UOM
}

var provider Provider

func GetProvider() *Provider {
	return &provider
}

func init() {
	provider.unitsOfMeasure = make(map[string]*UOM)

	// Information
	provider.unitsOfMeasure[BYTES] = &UOM{Name: BYTES, Entity: EntityInformation, IsBase: true}
	provider.unitsOfMeasure[KB] = &UOM{Name: KB, Entity: EntityInformation, A: 1024, B: 0}
	provider.unitsOfMeasure[MB] = &UOM{Name: MB, Entity: EntityInformation, A: 1024 * 1024, B: 0}
	provider.unitsOfMeasure[GB] = &UOM{Name: GB, Entity: EntityInformation, A: 1024 * 1024 * 1024, B: 0}
	provider.unitsOfMeasure[TB] = &UOM{Name: TB, Entity: EntityInformation, A: 1024 * 1024 * 1024 * 1024, B: 0}

	// Time
	provider.unitsOfMeasure[SEC] = &UOM{Name: SEC, Entity: EntityTime, IsBase: true}
	provider.unitsOfMeasure[MS] = &UOM{Name: MS, Entity: EntityTime, A: 0.001, B: 0}
	provider.unitsOfMeasure[MIN] = &UOM{Name: MIN, Entity: EntityTime, A: 60, B: 0}
	provider.unitsOfMeasure[HOUR] = &UOM{Name: HOUR, Entity: EntityTime, A: 3600, B: 0}
	provider.unitsOfMeasure[DAY] = &UOM{Name: DAY, Entity: EntityTime, A: 86400, B: 0}
}
