package tagcloud

import (
	"sort"
)

// TagCloud aggregates statistics about used tags
type TagCloud struct {
	Tags map[string]*TagStat
}

// TagStat represents statistics regarding single tag
type TagStat struct {
	Tag             string
	OccurrenceCount int
}

// New should create a valid TagCloud instance
func New() *TagCloud {
	return &TagCloud{
		Tags: make(map[string]*TagStat),
	}
}

// AddTag should add a tag to the cloud if it wasn't present and increase tag occurrence count
// thread-safety is not needed
func (tagCloud *TagCloud) AddTag(tag string) {
	tagStat, isExisting := tagCloud.Tags[tag]
	if isExisting {
		tagStat.OccurrenceCount += 1
	} else {
		tagCloud.Tags[tag] = &TagStat{Tag: tag, OccurrenceCount: 1}
	}
}

// TopN should return top N most frequent tags ordered in descending order by occurrence count
// if there are multiple tags with the same occurrence count then the order is defined by implementation
// if n is greater that TagCloud size then all elements should be returned
// thread-safety is not needed
// there are no restrictions on time complexity
func (tagCloud *TagCloud) TopN(n int) []TagStat {
	tagStats := make([]TagStat, 0, len(tagCloud.Tags))
	for _, tagStat := range tagCloud.Tags {
		tagStats = append(tagStats, *tagStat)
	}
	sort.Slice(tagStats, func(i, j int) bool {
		return tagStats[i].OccurrenceCount > tagStats[j].OccurrenceCount
	})
	sliceLimit := n
	if sliceLimit > len(tagStats) {
		sliceLimit = len(tagStats)
	}
	return tagStats[:sliceLimit]
}
