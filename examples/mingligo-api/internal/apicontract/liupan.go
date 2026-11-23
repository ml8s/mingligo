package apicontract

// Liupan tools expose one domain method with one schema form per time layer.
// The distinct forms keep required facts exact (for example, liu_shi needs a
// longitude while xiao_yun does not) without turning one input DTO into a bag
// of mutually optional fields.
