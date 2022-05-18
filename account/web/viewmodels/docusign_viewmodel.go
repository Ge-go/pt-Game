package viewmodels

//信封回调结构体
type NotifyEnvelopeXml struct {
	EnvelopeStatus EnvelopeStatus `xml:"EnvelopeStatus"`
}

type EnvelopeStatus struct {
	EnvelopeID        string                 `xml:"EnvelopeID"`
	Status            string                 `xml:"Status"`
	Declined          string                 `xml:"Declined"`
	Completed         string                 `xml:"Completed"`
	RecipientStatuses []RecipientStatusesXml `xml:"RecipientStatuses"`
}

type RecipientStatusesXml struct {
	RecipientStatus RecipientStatus `xml:"RecipientStatus"`
}

type RecipientStatus struct {
	DeclineReason string `xml:"DeclineReason"`
	Status        string `xml:"Status"`
	ClientUserId  string `xml:"ClientUserId"`
}
