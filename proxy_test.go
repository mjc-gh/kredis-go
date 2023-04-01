package kredis

import "time"

func (s *KredisTestSuite) TestNewProxyPrefixesNamespace() {
	p, e := NewProxy("key", Options{})

	s.NoError(e)
	s.Equal("ns:key", p.key)
}

func (s *KredisTestSuite) TestNewProxyWithExpiresIn() {
	p, e := NewProxy("key", Options{ExpiresIn: "1m15s"})

	s.NoError(e)
	s.Equal(75*time.Second, p.expiresIn)
}
