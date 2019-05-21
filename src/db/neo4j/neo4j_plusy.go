package neo4j

func (n *Neo4j) SaveOffset(offset int64) error {
	return nil
}

func (n *Neo4j) GetNextOffset() (int64, error) {
	return -1, nil
}
