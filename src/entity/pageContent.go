package entity

type pageContent struct {
	id                  string
	tag                 string
	content             string
	className           string
	href                string
	styles              string
	children            []pageContent
	inputType           string
	placeholder         string
	formType            string
	formAction          string
	value               string
	name                string
	checked             bool
	passwordViewToggler bool
	targetedPassword    string
	htmlId              string
}

func (p *pageContent) GetId() string {
	return p.id
}

func (p *pageContent) SetId(id string) {
	p.id = id
}

func (p *pageContent) GetTag() string {
	return p.tag
}

func (p *pageContent) SetTag(tag string) {
	p.tag = tag
}

func (p *pageContent) GetContent() string {
	return p.content
}

func (p *pageContent) SetContent(content string) {
	p.content = content
}

func (p *pageContent) GetClassName() string {
	return p.className
}

func (p *pageContent) SetClassName(className string) {
	p.className = className
}

func (p *pageContent) GetHref() string {
	return p.href
}

func (p *pageContent) SetHref(href string) {
	p.href = href
}

func (p *pageContent) GetStyles() string {
	return p.styles
}

func (p *pageContent) SetStyles(styles string) {
	p.styles = styles
}

func (p *pageContent) GetChildren() []pageContent {
	return p.children
}

func (p *pageContent) SetChildren(children []pageContent) {
	p.children = children
}

func (p *pageContent) ToJSON() PageContentJSON {
	var children []PageContentJSON
	for _, child := range p.children {
		children = append(children, child.ToJSON())
	}

	return PageContentJSON{
		Id:                  p.id,
		Tag:                 p.tag,
		Content:             p.content,
		ClassName:           p.className,
		Href:                p.href,
		Styles:              p.styles,
		Children:            children,
		InputType:           p.inputType,
		Placeholder:         p.placeholder,
		FormType:            p.formType,
		FormAction:          p.formAction,
		Value:               p.value,
		Checked:             p.checked,
		Name:                p.name,
		PasswordViewToggler: p.passwordViewToggler,
		TargetedPassword:    p.targetedPassword,
		HtmlId:              p.htmlId,
	}
}
