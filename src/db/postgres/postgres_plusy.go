package postgres

func (p *Postgres) SaveOffset(offset int64) error {
	return nil
}

func (p *Postgres) GetNextOffset() (int64, error) {
	return -1, nil
}
