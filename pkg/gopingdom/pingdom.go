package gopingdom

import "github.com/jsdidierlaurent/go-pingdom/pingdom"

type PingdomCheckApi interface {
	List(params ...map[string]string) ([]pingdom.CheckResponse, error)
	Read(id int) (*pingdom.CheckResponse, error)
}
