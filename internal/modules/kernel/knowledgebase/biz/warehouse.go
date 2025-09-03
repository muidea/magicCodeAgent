package biz

func (s *KnowledgeBase) NewWarehouse(warehouse string) error {
	if err := s.daoPtr.CheckWarehouse(warehouse); err == nil {
		return err
	}

	return s.daoPtr.NewWarehouse(warehouse)
}

func (s *KnowledgeBase) DropWarehouse(warehouse string) error {
	if err := s.daoPtr.CheckWarehouse(warehouse); err != nil {
		return err
	}

	return s.daoPtr.DropWarehouse(warehouse)
}

func (s *KnowledgeBase) ListWarehouse() (ret []string, err error) {
	return s.daoPtr.ListWarehouse()
}
