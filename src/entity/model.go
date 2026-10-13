package entity

import (
	"fmt"
)

type Model struct {
	nom       string
	attributs []Champs
}

func (m *Model) ToJSON() ModelJSON {
	var attr []ChampsJSON

	for _, val := range m.attributs {
		attr = append(attr, val.ToJSON())
	}

	return ModelJSON{
		Nom:       m.nom,
		Attributs: attr,
	}
}

func (m *Model) GetNom() string {
	return m.nom
}

func (m *Model) GetAttributs() []Champs {
	return m.attributs
}

func (m *Model) SetNom(nom string) error {
	if nom == "" {
		return fmt.Errorf("Empty value")
	}
	m.nom = nom
	return nil
}

func (m *Model) SetAttributs(attr []Champs) {
	m.attributs = attr
}

func (m *Model) AddAttribut(attr Champs) {
	m.attributs = addChamps(m.attributs, attr)
}

func (m *Model) DeleteAttribut(attrId int) {
	m.attributs = deleteChamp(m.attributs, attrId)
}
