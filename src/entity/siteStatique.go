package entity

type SiteStatique struct {
	pages      []Page
	composants []Composant
}

type AssetJSONList []AssetJSON

func (s *SiteStatique) ToJSON() SiteStatiqueJSON {
	var pages []PageJSON
	var composants []ComposantJSON
	for _, val := range s.pages {
		pages = append(pages, val.ToJSON())
	}
	for _, val := range s.composants {
		composants = append(composants, val.ToJSON())
	}
	return SiteStatiqueJSON{
		Pages:      pages,
		Composants: composants,
	}
}

func (sitestatique *SiteStatique) GetPages() []Page {
	return sitestatique.pages
}

func (sitestatique *SiteStatique) GetComposants() []Composant {
	return sitestatique.composants
}

func (sitestatique *SiteStatique) SetPages(Pages []Page) *SiteStatique {
	sitestatique.pages = Pages
	return sitestatique
}

func (sitestatique *SiteStatique) SetComposants(Composants []Composant) *SiteStatique {
	sitestatique.composants = Composants
	return sitestatique
}
