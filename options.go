package kredis

type Options struct {
	Config       *string
	ExpiresIn    string
	DefaultValue any // TODO deprecate this field
}

func (o *Options) GetConfig() string {
	if o.Config != nil {
		return *o.Config
	}

	return "shared"
}
