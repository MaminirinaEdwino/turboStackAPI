package entity

type Champs struct {
	nom          string
	type_champs         string
	default_value any
	constraint   []any
}

func (c *Champs) ToJSON() ChampsJSON  {
	return ChampsJSON{
		Nom: c.nom,
		Type: c.type_champs,
		DefaultValue: c.default_value,
		Constraint: c.constraint,
	}
}

func (c *Champs) GetNom() string  {
	return c.nom
}

func (c *Champs) GetType() string {
	return c.type_champs
}

func (c *Champs) GetDefaultValue() any {
	return c.default_value
}

func (c *Champs) GetConstraint() []any {
	return c.constraint
}

func (c *Champs) SetNom(nom string) {
	c.nom = nom
}

func (c *Champs) SetType(type_champs string) {
	c.type_champs = type_champs
}

func (c *Champs) SetDefaultValue(defaultValue any) {
	c.default_value = defaultValue
}

func (c *Champs) SetContraint(constraint []any) {
	c.constraint = constraint
}