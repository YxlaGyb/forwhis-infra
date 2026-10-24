package types

type Payload struct {
	JID     string `json:"jid"`
	UID     int    `json:"uid"`
	Email   string `json:"email"`
	Expired int64  `json:"expired"`
	Create  int64  `json:"create"`
	ISS     string `json:"iss"`
	AUD     string `json:"aud"`
}
