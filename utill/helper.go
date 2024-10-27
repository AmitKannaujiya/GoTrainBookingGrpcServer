package utill

import (
	s "train-book/pkg/models"
)

func GetSectionType(section string) s.TRAIN_SECTION {
	switch section {
	case SECTION_A:
		return s.TRAIN_SECTION_A
	case SECTION_B:
		return s.TRAIN_SECTION_B
	}
	return s.TRAIN_SECTION_UNDEFINED
}
