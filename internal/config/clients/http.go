package clients

type Http struct {
	Wiki *Wiki `yaml:"wiki" env-required:"true"`
}

type Wiki struct {
	Host        string `yaml:"host" env-required:"true"`
	Port        int    `yaml:"port" env-default:"0"`
	Timeout     string `yaml:"timeout" env-default:"5s"`
	TokenId     string `yaml:"token_id" env-required:"true"`
	TokenSecret string `yaml:"token_secret" env-required:"true"`
	Paths       *Paths `yaml:"paths" env-required:"true"`
}

type Paths struct {
	PagesExportMarkdown string `yaml:"pages_export_markdown" env-required:"true"`
}
