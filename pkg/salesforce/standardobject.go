package salesforce

type StandardObject struct {
	Title string
	Value string
	Topic string
}

var StandardObjectList = []StandardObject{
	{Title: "Account", Value: "Account"},
	{Title: "Event", Value: "Event"},
	{Title: "Opportunity", Value: "Opportunity"},
	{Title: "Case", Value: "Case"},
	{Title: "Order", Value: "Order"},
}
