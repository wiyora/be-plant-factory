package config

type FrontendConfig struct {
	AllowedUrl StringList `env:"ALLOWED_URL" validate:"required,unique,dive,required,url" default:"[\"localhost:3000\"]"`

	mapAllowedUrl map[string]struct{}
}

func (c *FrontendConfig) Init() {
	c.mapAllowedUrl = make(map[string]struct{})
	for _, url := range c.AllowedUrl {
		c.mapAllowedUrl[url] = struct{}{}
	}
}

func (c *FrontendConfig) IsValidAllowedUrl(fullUrl string) bool {
	_, ok := c.mapAllowedUrl[fullUrl]
	return ok
}

func (c *FrontendConfig) DefaultAllowedUrl() string {
	if len(c.AllowedUrl) > 0 {
		return c.AllowedUrl[0]
	}

	return ""
}
