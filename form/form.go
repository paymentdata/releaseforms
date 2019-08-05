package form

//ReleaseTemplateData is the encapulative struct for software release forms.
type ReleaseTemplateData struct {

	//Release Date
	Date string `json:"Date"`
	//Product
	Product string `json:"Product"`

	//Included changes
	Commit string `json:"Commit"`
	Author string `json:"Author"`
}

//Commit is the Change Item primitive
type Commit struct {
	Text        string
	RequestedBy string
	SummaryOfChanges,
	Notes,
	Developer,
	TestedBy,
	CodeReviewAndTesting,
	CodeReviewAndTestingNotes,
	ApprovedBy string
}
