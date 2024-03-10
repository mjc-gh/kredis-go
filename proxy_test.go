package kredis

import "time"

func (s *KredisTestSuite) TestNewProxyPrefixesNamespace() {
	p, e := NewProxy("key")

	s.NoError(e)
	s.Equal("ns:key", p.key)
}

func (s *KredisTestSuite) TestNewProxyWithConfigName() {
	p, e := NewProxy("key", WithConfigName("badconn"))

	s.NoError(e)
	s.NotNil(p)
	s.Equal("key", p.key)
}

func (s *KredisTestSuite) TestNewProxyWithExpiresIn() {
	p, e := NewProxy("key")
	s.NoError(e)
	s.Equal(0*time.Second, p.expiresIn)

	p, e = NewProxy("key", WithExpiry("1m15s"))
	s.NoError(e)
	s.Equal(75*time.Second, p.expiresIn)
}

func (s *KredisTestSuite) TestProxyKey() {
	p, _ := NewProxy("key")

	s.Equal("ns:key", p.Key())
}
