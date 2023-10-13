package source

type Source interface {
	Load() (string, error)
}
