package datadir

type Basis struct {
	Dir
}

func NewBasis(path string) (*Basis, error) {
	dir, err := newRootDir(path)
	if err != nil {
		return nil, err
	}

	return &Basis{Dir: dir}, nil
}
