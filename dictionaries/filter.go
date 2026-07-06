package filter

import (
	"fmt"

	sensitive "github.com/LuYongwang/go-sensitive-word"
	"go.uber.org/zap"
)

func InitFilter(logger *zap.Logger) (*sensitive.Manager, error) {
	// Инициализация фильтра
	filter, err := sensitive.NewFilter(
		sensitive.StoreOption{Type: sensitive.StoreMemory},
		sensitive.FilterOption{Type: sensitive.FilterDfa},
	)
	if err != nil {
		logger.Error(fmt.Sprintf("Error creating filter: %v", err))
		return &sensitive.Manager{}, err
	}
	// Загрузка словаря Русских ругательств из файла
	err = filter.LoadDictPath("dictionaries/russian-bad-words.txt")
	if err != nil {
		logger.Error(fmt.Sprintf("Error loading dictionary: %v", err))
		return &sensitive.Manager{}, err
	}
	return filter, nil
}
