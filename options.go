package kredis

type Options struct {
	Config    *string
	ExpiresIn string
}

func (o *Options) GetConfig() string {
	if o.Config != nil {
		return *o.Config
	}

	return "shared"
}
