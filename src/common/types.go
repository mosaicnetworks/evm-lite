package common

type Genesis struct {
	Alloc AccountMap
	Poa   PoaMap
}

type AccountMap map[string]struct {
	Code        string
	Storage     map[string]string
	Balance     string
	Authorising bool
}

type PoaMap struct {
	Address string
	Balance string
	Abi     string
	Code    string
}
