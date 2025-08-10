package enum

type AgentType uint

const (
	AgentTypeGeneral AgentType = iota + 1
	AgentTypeCustomer
	AgentTypeCorporation
	AgentTypeAdmin
)
