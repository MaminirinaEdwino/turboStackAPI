package entity

func (aj *AbonnementJSON) ToModel() Abonnement {
	return Abonnement{
		type_abonnement: aj.Type_abonnement,
		debut:           aj.Debut,
		fin:             aj.Fin,
	}
}

func (cj *ChampsJSON) ToModel() Champs {
	return Champs{
		nom:           cj.Nom,
		type_champs:   cj.Type,
		default_value: cj.DefaultValue,
		constraint:    cj.Constraint,
	}
}

func (cj *ComposantJSON) ToModel() Composant {
	var composant []pageContent
	for _, val := range cj.Contenu {
		composant = append(composant, val.ToModel())
	}
	return Composant{
		nom:     cj.Nom,
		contenu: composant,
		params:  cj.Params,
	}
}

func (ej *EndpointJSON) ToModel() Endpoint {
	var model []Model
	for _, val := range ej.Model {
		model = append(model, val.ToModel())
	}
	var returnModel []Model
	for _, val := range ej.ReturnContent {
		returnModel = append(returnModel, val.ToModel())
	}
	return Endpoint{
		nom:               ej.Nom,
		uri:               ej.Uri,
		method:            ej.Method,
		model:             model,
		params:            ej.Params,
		role:              ej.Role,
		logic:             ej.Logic,
		returnContent:     returnModel,
		returnContentType: ej.ReturnContentType,
		returnPage:        ej.ReturnPage,
		redirectUri:       ej.RedirectUri,
		bodyType:          ej.BodyType,
	}
}

func (mj *ModelJSON) ToModel() Model {
	var attr []Champs
	for _, val := range mj.Attributs {
		attr = append(attr, val.ToModel())
	}
	return Model{
		nom:       mj.Nom,
		attributs: attr,
	}
}

func (pj *PageJSON) ToModel() Page {
	var pcontent []pageContent
	for _, val := range pj.Contenu {
		pcontent = append(pcontent, val.ToModel())
	}
	return Page{
		nom:     pj.Nom,
		contenu: pcontent,
		uri:     pj.Uri,
		styles:  pj.Styles,
	}
}

func (bddjson *BDDJSON) ToModel() BDD {
	var dbModel []Model
	for _, val := range bddjson.Models {
		dbModel = append(dbModel, val.ToModel())
	}
	return BDD{
		sgbd:    bddjson.Sgbd,
		models:  dbModel,
		scripts: bddjson.Scripts,
	}
}

func (pj *ProjectJSON) ToModel() Project {
	var assets AssetList
	for _, asset := range pj.Assets {
		assets = append(assets, asset.ToModel())
	}
	return Project{
		nom:           pj.Nom,
		type_project:  pj.Type,
		description:   pj.Description,
		bdd:           pj.BDD.ToModel(),
		rest_api:      pj.RestApi.ToModel(),
		web_app:       pj.WebApp.ToModel(),
		site_statique: pj.SiteStatique.ToModel(),
		updateAt:      pj.UpdateAt,
		createdAt:     pj.CreatedAt,
		assets:        assets,
	}
}

func (rj *RestApiJSON) ToModel() RestApi {
	var endpoints []Endpoint

	for _, val := range rj.Endpoints {
		endpoints = append(endpoints, val.ToModel())
	}
	return RestApi{
		endpoints: endpoints,
		role:      rj.Role,
		bdd:       rj.BDD.ToModel(),
	}
}

func (sj *SiteStatiqueJSON) ToModel() SiteStatique {
	var pages []Page
	var composants []Composant
	for _, val := range sj.Pages {
		pages = append(pages, val.ToModel())
	}
	for _, val := range sj.Composants {
		composants = append(composants, val.ToModel())
	}
	return SiteStatique{
		pages:      pages,
		composants: composants,
	}
}

func (wj *WebAppJSON) ToModel() WebApp {
	var pages []Page
	var composants []Composant
	var controllers []Controller
	for _, val := range wj.Pages {
		pages = append(pages, val.ToModel())
	}
	for _, val := range wj.Composant {
		composants = append(composants, val.ToModel())
	}
	for _, val := range wj.Controllers {
		controllers = append(controllers, val.ToModel())
	}
	return WebApp{
		pages:       pages,
		role:        wj.Role,
		composant:   composants,
		controllers: controllers,
		bdd:         wj.BDD.ToModel(),
	}
}

func (cj *ControllerJSON) ToModel() Controller {
	return Controller{
		Nom:       cj.Nom,
		PageCible: cj.PageCible,
		Bindings:  cj.Bindings,
	}
}

func (pc *PageContentJSON) ToModel() pageContent {
	var children []pageContent
	for _, val := range pc.Children {
		children = append(children, val.ToModel())
	}

	return pageContent{
		id:                  pc.Id,
		tag:                 pc.Tag,
		content:             pc.Content,
		className:           pc.ClassName,
		href:                pc.Href,
		styles:              pc.Styles,
		children:            children,
		inputType:           pc.InputType,
		placeholder:         pc.Placeholder,
		formType:            pc.FormType,
		value:               pc.Value,
		name:                pc.Name,
		checked:             pc.Checked,
		formAction:          pc.FormAction,
		passwordViewToggler: pc.PasswordViewToggler,
		targetedPassword:    pc.TargetedPassword,
		htmlId:              pc.HtmlId,
	}
}

type AssetJSON struct {
	FileName    string `json:"file_name"`
	Base64Image string `json:"base_64_image"`
}

func (assetJSON *AssetJSON) ToModel() Asset {
	return Asset{
		fileName:    assetJSON.FileName,
		base64Image: assetJSON.Base64Image,
	}
}

func (asset *Asset) ToJSON() AssetJSON {
	return AssetJSON{
		FileName:    asset.fileName,
		Base64Image: asset.base64Image,
	}
}
