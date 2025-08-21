package internal

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

type IndexStrategy int

const (
	IndexStrategyNone IndexStrategy = iota
	IndexStrategyNumeric
)

func (is IndexStrategy) String() string {
	switch is {
	case IndexStrategyNumeric:
		return "Numeric"
	case IndexStrategyNone:
		return "None"
	default:
		return "None"
	}
}

func ParseIndexStrategy(s string) IndexStrategy {
	switch s {
	case "Numeric":
		return IndexStrategyNumeric
	case "None":
		return IndexStrategyNone
	default:
		return IndexStrategyNone
	}
}

type NumericConfig struct {
	DirPriority bool `json:"dir_priority"`
}

type IndexConfig struct {
	Strategy      IndexStrategy  `json:"strategy"`
	NumericConfig *NumericConfig `json:"numeric_config,omitempty"`
}

func (ic IndexConfig) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Strategy      string         `json:"strategy"`
		NumericConfig *NumericConfig `json:"numeric_config,omitempty"`
	}{
		Strategy:      ic.Strategy.String(),
		NumericConfig: ic.NumericConfig,
	})
}

func (ic *IndexConfig) UnmarshalJSON(data []byte) error {
	var aux struct {
		Strategy      string         `json:"strategy"`
		NumericConfig *NumericConfig `json:"numeric_config,omitempty"`
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	ic.Strategy = ParseIndexStrategy(aux.Strategy)
	ic.NumericConfig = aux.NumericConfig
	return nil
}

func getIndexConfig(dirPath string) (*IndexConfig, error) {
	indexFilePath := filepath.Join(dirPath, ".index")

	if _, err := os.Stat(indexFilePath); os.IsNotExist(err) {
		return &IndexConfig{Strategy: IndexStrategyNone}, nil
	}

	data, err := os.ReadFile(indexFilePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read .index file: %w", err)
	}

	var config IndexConfig
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse .index file: %w", err)
	}

	return &config, nil
}

func writeIndexConfig(dirPath string, config *IndexConfig) error {
	indexFilePath := filepath.Join(dirPath, ".index")

	if config.Strategy == IndexStrategyNone {
		if _, err := os.Stat(indexFilePath); !os.IsNotExist(err) {
			return os.Remove(indexFilePath)
		}
		return nil
	}

	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal index config: %w", err)
	}

	return os.WriteFile(indexFilePath, data, 0644)
}

type IndexedEntry struct {
	Index    int
	Name     string
	FullName string
	IsDir    bool
}

func parseIndexedName(name string) (int, string, bool) {
	re := regexp.MustCompile(`^(\d+)\s*-\s*(.+)$`)
	matches := re.FindStringSubmatch(name)
	if len(matches) == 3 {
		if index, err := strconv.Atoi(matches[1]); err == nil {
			return index, matches[2], true
		}
	}
	return 0, name, false
}

func findNextIndex(dirPath string) (int, error) {
	config, err := getIndexConfig(dirPath)
	if err != nil {
		return 1, err
	}

	if config.Strategy == IndexStrategyNone {
		return 1, nil
	}

	entries, err := os.ReadDir(dirPath)
	if err != nil {
		return 1, err
	}

	maxIndex := 0
	for _, entry := range entries {
		name := entry.Name()
		if strings.HasPrefix(name, ".") {
			continue
		}

		if index, _, isIndexed := parseIndexedName(name); isIndexed {
			if index > maxIndex {
				maxIndex = index
			}
		}
	}

	return maxIndex + 1, nil
}

func getIndexedEntries(dirPath string) ([]IndexedEntry, error) {
	entries, err := os.ReadDir(dirPath)
	if err != nil {
		return nil, err
	}

	var indexedEntries []IndexedEntry
	for _, entry := range entries {
		name := entry.Name()

		if strings.HasPrefix(name, ".") {
			continue
		}

		isDir := entry.IsDir()
		if !isDir && filepath.Ext(name) != ".md" {
			continue
		}

		index, cleanName, isIndexed := parseIndexedName(name)
		if !isIndexed {
			index = 0
			cleanName = name
		}

		indexedEntries = append(indexedEntries, IndexedEntry{
			Index:    index,
			Name:     cleanName,
			FullName: name,
			IsDir:    isDir,
		})
	}

	sort.Slice(indexedEntries, func(i, j int) bool {
		if indexedEntries[i].Index != indexedEntries[j].Index {
			return indexedEntries[i].Index < indexedEntries[j].Index
		}
		return indexedEntries[i].Name < indexedEntries[j].Name
	})

	return indexedEntries, nil
}

func applyNumericIndexing(dirPath string, dirPriority bool) error {
	entries, err := getIndexedEntries(dirPath)
	if err != nil {
		return err
	}

	if dirPriority {
		sort.Slice(entries, func(i, j int) bool {
			if entries[i].IsDir != entries[j].IsDir {
				return entries[i].IsDir
			}
			return entries[i].Name < entries[j].Name
		})
	} else {
		sort.Slice(entries, func(i, j int) bool {
			return entries[i].Name < entries[j].Name
		})
	}

	for i, entry := range entries {
		newIndex := i + 1
		oldPath := filepath.Join(dirPath, entry.FullName)

		var newName string
		if entry.IsDir {
			newName = fmt.Sprintf("%d - %s", newIndex, entry.Name)
		} else {
			newName = fmt.Sprintf("%d - %s", newIndex, entry.Name)
		}

		if newName != entry.FullName {
			newPath := filepath.Join(dirPath, newName)
			if err := os.Rename(oldPath, newPath); err != nil {
				return fmt.Errorf("failed to rename %s to %s: %w", oldPath, newPath, err)
			}
		}
	}

	return nil
}

func removeIndexing(dirPath string) error {
	entries, err := getIndexedEntries(dirPath)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		if entry.Index > 0 {
			oldPath := filepath.Join(dirPath, entry.FullName)
			newPath := filepath.Join(dirPath, entry.Name)
			if err := os.Rename(oldPath, newPath); err != nil {
				return fmt.Errorf("failed to rename %s to %s: %w", oldPath, newPath, err)
			}
		}
	}

	return nil
}

type ValidationResult struct {
	IsValid     bool
	Issues      []string
	MissingGaps []int
	Duplicates  map[int][]string
	Unindexed   []string
}

func ValidResult() *ValidationResult {
	return &ValidationResult{
		IsValid:     true,
		Issues:      []string{},
		MissingGaps: []int{},
		Duplicates:  make(map[int][]string),
		Unindexed:   []string{},
	}
}

func (dirInfo *DirInfo) validateIndexing() (*ValidationResult, error) {
	if dirInfo.IndexingStrategy == IndexStrategyNumeric {
		return validateNumericIndexing(dirInfo.AbsolutePath)
	}
	result := ValidResult()
	return result, nil
}

func validateNumericIndexing(dirPath string) (*ValidationResult, error) {
	config, err := getIndexConfig(dirPath)
	if err != nil {
		return nil, err
	}

	result := ValidResult()
	if config.Strategy != IndexStrategyNumeric {
		result.Issues = append(result.Issues, "Directory is not using numeric indexing")
		result.IsValid = false
		return result, nil
	}

	entries, err := getIndexedEntries(dirPath)
	if err != nil {
		return nil, err
	}

	indexMap := make(map[int][]string)
	maxIndex := 0

	for _, entry := range entries {
		if entry.Index == 0 {
			result.Unindexed = append(result.Unindexed, entry.FullName)
			result.IsValid = false
		} else {
			indexMap[entry.Index] = append(indexMap[entry.Index], entry.FullName)
			if entry.Index > maxIndex {
				maxIndex = entry.Index
			}
		}
	}

	for index, files := range indexMap {
		if len(files) > 1 {
			result.Duplicates[index] = files
			result.IsValid = false
		}
	}

	for i := 1; i <= maxIndex; i++ {
		if _, exists := indexMap[i]; !exists {
			result.MissingGaps = append(result.MissingGaps, i)
			result.IsValid = false
		}
	}

	if len(result.Unindexed) > 0 {
		result.Issues = append(result.Issues, fmt.Sprintf("Found %d unindexed files", len(result.Unindexed)))
	}

	if len(result.Duplicates) > 0 {
		result.Issues = append(result.Issues, fmt.Sprintf("Found %d duplicate indices", len(result.Duplicates)))
	}

	if len(result.MissingGaps) > 0 {
		result.Issues = append(result.Issues, fmt.Sprintf("Found %d gaps in indexing", len(result.MissingGaps)))
	}

	return result, nil
}

func repairNumericIndexing(dirPath string) error {
	config, err := getIndexConfig(dirPath)
	if err != nil {
		return err
	}

	if config.Strategy != IndexStrategyNumeric {
		return fmt.Errorf("directory is not using numeric indexing")
	}

	dirPriority := true
	if config.NumericConfig != nil {
		dirPriority = config.NumericConfig.DirPriority
	}

	return applyNumericIndexing(dirPath, dirPriority)
}

func moveIndexedFile(dirPath, fileName string, newIndex int) error {
	config, err := getIndexConfig(dirPath)
	if err != nil {
		return err
	}

	if config.Strategy != IndexStrategyNumeric {
		return fmt.Errorf("directory is not using numeric indexing")
	}

	entries, err := getIndexedEntries(dirPath)
	if err != nil {
		return err
	}

	var targetEntry *IndexedEntry
	for i, entry := range entries {
		if entry.FullName == fileName {
			targetEntry = &entries[i]
			break
		}
	}

	if targetEntry == nil {
		return fmt.Errorf("file %s not found in directory", fileName)
	}

	if newIndex < 1 {
		return fmt.Errorf("index must be >= 1")
	}

	maxPossibleIndex := len(entries)
	if newIndex > maxPossibleIndex {
		return fmt.Errorf("index %d exceeds maximum possible index %d", newIndex, maxPossibleIndex)
	}

	currentIndex := targetEntry.Index
	if currentIndex == newIndex {
		return nil
	}

	oldPath := filepath.Join(dirPath, targetEntry.FullName)

	var newName string
	if targetEntry.IsDir {
		newName = fmt.Sprintf("%d - %s", newIndex, targetEntry.Name)
	} else {
		newName = fmt.Sprintf("%d - %s", newIndex, targetEntry.Name)
	}

	tempPath := filepath.Join(dirPath, ".temp_"+newName)
	if err := os.Rename(oldPath, tempPath); err != nil {
		return fmt.Errorf("failed to move file to temporary location: %w", err)
	}

	if err := shiftIndicesForMove(dirPath, currentIndex, newIndex); err != nil {
		os.Rename(tempPath, oldPath)
		return fmt.Errorf("failed to shift indices: %w", err)
	}

	newPath := filepath.Join(dirPath, newName)
	if err := os.Rename(tempPath, newPath); err != nil {
		return fmt.Errorf("failed to move file to final location: %w", err)
	}

	return nil
}

func shiftIndicesForMove(dirPath string, fromIndex, toIndex int) error {
	entries, err := getIndexedEntries(dirPath)
	if err != nil {
		return err
	}

	if fromIndex == toIndex {
		return nil
	}

	var updates []struct {
		oldPath string
		newPath string
	}

	if fromIndex < toIndex {
		for _, entry := range entries {
			if entry.Index > fromIndex && entry.Index <= toIndex && !strings.HasPrefix(entry.FullName, ".temp_") {
				newIndex := entry.Index - 1
				var newName string
				if entry.IsDir {
					newName = fmt.Sprintf("%d - %s", newIndex, entry.Name)
				} else {
					newName = fmt.Sprintf("%d - %s", newIndex, entry.Name)
				}

				updates = append(updates, struct {
					oldPath string
					newPath string
				}{
					oldPath: filepath.Join(dirPath, entry.FullName),
					newPath: filepath.Join(dirPath, newName),
				})
			}
		}
	} else {
		for _, entry := range entries {
			if entry.Index >= toIndex && entry.Index < fromIndex && !strings.HasPrefix(entry.FullName, ".temp_") {
				newIndex := entry.Index + 1
				var newName string
				if entry.IsDir {
					newName = fmt.Sprintf("%d - %s", newIndex, entry.Name)
				} else {
					newName = fmt.Sprintf("%d - %s", newIndex, entry.Name)
				}

				updates = append(updates, struct {
					oldPath string
					newPath string
				}{
					oldPath: filepath.Join(dirPath, entry.FullName),
					newPath: filepath.Join(dirPath, newName),
				})
			}
		}
	}

	for _, update := range updates {
		if err := os.Rename(update.oldPath, update.newPath); err != nil {
			return fmt.Errorf("failed to rename %s to %s: %w", update.oldPath, update.newPath, err)
		}
	}

	return nil
}

func moveIndexedFileUp(dirPath, fileName string) error {
	entries, err := getIndexedEntries(dirPath)
	if err != nil {
		return err
	}

	var targetEntry *IndexedEntry
	for i, entry := range entries {
		if entry.FullName == fileName {
			targetEntry = &entries[i]
			break
		}
	}

	if targetEntry == nil {
		return fmt.Errorf("file %s not found in directory", fileName)
	}

	if targetEntry.Index <= 1 {
		return fmt.Errorf("file is already at the top position")
	}

	return moveIndexedFile(dirPath, fileName, targetEntry.Index-1)
}

func moveIndexedFileDown(dirPath, fileName string) error {
	entries, err := getIndexedEntries(dirPath)
	if err != nil {
		return err
	}

	var targetEntry *IndexedEntry
	maxIndex := 0
	for i, entry := range entries {
		if entry.FullName == fileName {
			targetEntry = &entries[i]
		}
		if entry.Index > maxIndex {
			maxIndex = entry.Index
		}
	}

	if targetEntry == nil {
		return fmt.Errorf("file %s not found in directory", fileName)
	}

	if targetEntry.Index >= maxIndex {
		return fmt.Errorf("file is already at the bottom position")
	}

	return moveIndexedFile(dirPath, fileName, targetEntry.Index+1)
}

