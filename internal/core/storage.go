package core

var zsetStore map[string]*ZSet

// var setStore map[string]Set
// var dictStore *Dict
// var sbStore map[string]*SBChain
// var cmsStore map[string]*CMS

func init() {
	zsetStore = make(map[string]*ZSet)
	// setStore = make(map[string]Set)
	// dictStore = CreateDict()
	// sbStore = make(map[string]*SBChain)
	// cmsStore = make(map[string]*CMS)
}
