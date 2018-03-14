package handshake

type NewSessionTicket struct {
	TicketLifeTime	uint32
	TicketAgeAdd	uint32
	TicketNonce		[]byte
	Ticket			[]byte
	Extensions		[]Extension
}