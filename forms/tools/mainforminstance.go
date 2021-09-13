package tools

type MainFormInterface interface {
	ShowFullScreenValue(show bool, itemId string)
	ShowChartGroup(resId string)
}

var MainFormInstance MainFormInterface
