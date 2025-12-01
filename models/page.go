package models

type FooterData struct {
	Sitename string
}

type ParentPageData struct {
	Shortname      string
	Subtitle       string
	Footer         FooterData
	ChildViewModel any
}
