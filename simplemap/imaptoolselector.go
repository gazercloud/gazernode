package simplemap

type IMapToolSelector interface {
	CurrentTool() string
	ResetCurrentTool()
}
