package goapimaker

func DbCallerPG() string {
	return DBCallerTemplate()
}
func DbCallerPGWebApp() string {
	return DBCallerTemplateWebAPp()
}

func DBCallerHandler(sgbd string) string {
	if sgbd == "pg" {
		return DbCallerPG()
	}
	return ""
}
func DBCallerHandlerWebApp(sgbd string) string {
	if sgbd == "pg" {
		return DbCallerPG()
	}
	return ""
}
