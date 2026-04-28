//go:generate mockgen -source=good_interface.go -destination=good_interface_mock.go -package=example

package example

type GoodService interface {
	Do() error
}
