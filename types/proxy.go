package types

type ProxyDefinition struct {
	Delegate  AccountID
	ProxyType ProxyType
	Delay     BlockNumber
}

type ProxyStorageEntry struct {
	ProxyDefinitions []ProxyDefinition
	Balance          U128
}
