package kredis

type ScalarTime struct{ Proxy }

//func NewTime(key string, options Options) (*ScalarTime, error) {
//proxy, err := NewProxy(key, options)

//if err != nil {
//return nil, err
//}

//return &ScalarTime{Proxy: *proxy}, nil
//}

//func (s *ScalarTime) Value() time.Time {
//val, err := s.client.Get(s.ctx, s.key).Result()

//if err != nil {
//return time.Time{}
//}

//return val
//}

//func (s *ScalarTime) SetValue(v time.Time) error {
//s.client.Set(s.ctx, s.key, v, s.expiresIn)

//return nil
//}

//func (s *ScalarTime) DefaultValue() Time {
//switch s.defaultValue.(type) {
//case time.Time:
//return s.defaultValue.(time.Time)
//default:
//return time.Time{}
//}
//}
