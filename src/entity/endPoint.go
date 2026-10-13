package entity

type Endpoint struct {
	nom               string
	uri               string
	method            string
	model             []Model
	params            []string
	role              string
	returnContent     []Model
	logic             map[string]string
	returnContentType string
	returnPage        string
	redirectUri       string
	bodyType          string
}

func (e *Endpoint) ToJSON() EndpointJSON {
	var model []ModelJSON
	for _, val := range e.model {
		model = append(model, val.ToJSON())
	}
	var retunModel []ModelJSON
	for _, val := range e.returnContent {
		retunModel = append(retunModel, val.ToJSON())
	}
	return EndpointJSON{
		Nom:               e.nom,
		Uri:               e.uri,
		Method:            e.method,
		Model:             model,
		Params:            e.params,
		Role:              e.role,
		Logic:             e.logic,
		ReturnContent:     retunModel,
		ReturnContentType: e.returnContentType,
		ReturnPage:        e.returnPage,
		RedirectUri:       e.redirectUri,
		BodyType:          e.bodyType,
	}
}

func (e *Endpoint) GetNom() string {
	return e.nom
}

func (e *Endpoint) GetUri() string {
	return e.uri
}

func (e *Endpoint) GetMethod() string {
	return e.method
}

func (e *Endpoint) GetModel() []Model {
	return e.model
}

func (e *Endpoint) GetParams() []string {
	return e.params
}

func (e *Endpoint) GetRole() string {
	return e.role
}

func (e *Endpoint) SetNom(nom string) {
	e.nom = nom
}

func (e *Endpoint) SetUri(uri string) {
	e.uri = uri
}

func (e *Endpoint) SetMethod(method string) {
	e.method = method
}

func (e *Endpoint) SetModel(Model []Model) {
	e.model = Model
}

func (e *Endpoint) SetParams(Params []string) {
	e.params = Params
}

func (e *Endpoint) SetRole(Role string) {
	e.role = Role
}

func (e *Endpoint) GetLogic() map[string]string {
	return e.logic
}

func (e *Endpoint) SetLogic(logic map[string]string) {
	e.logic = logic
}
