package entity

type Composant struct {
	nom     string
	contenu []pageContent
	params  []any
}

func (c *Composant) ToJSON() ComposantJSON  {
	var composant []PageContentJSON
	for _, val := range c.contenu{
		composant = append(composant, val.ToJSON())
	}
	return ComposantJSON{
		Nom: c.nom,
		Contenu: composant,
		Params: c.params,
	}
}

func (composant *Composant) GetNom() string {
	return composant.nom
}

func (composant *Composant) GetContenu() []pageContent {
	return composant.contenu
}

func (composant *Composant) GetParams() any {
	return composant.params
}

func (composant *Composant) SetNom(nom string) *Composant {
	composant.nom = nom
	return composant
}

func (composant *Composant) SetContenu(contenu []pageContent) *Composant {
	composant.contenu = contenu
	return composant
}

func (composant *Composant) SetParams(params []any) *Composant {
	composant.params = params
	return composant
}