package conf

type Generate struct {
	PackageName string `yaml:"package_name"`
	ServiceName string `yaml:"service_name"`
	ApiName     string `yaml:"api_name"`
	TableName   string `yaml:"table_name"`
	Sql         string `yaml:"sql"`
}

func (g *Generate) Update(packageName, serviceName, apiName, tableName string) {
	g.PackageName = packageName
	g.ServiceName = serviceName
	g.ApiName = apiName
	g.TableName = tableName
}
