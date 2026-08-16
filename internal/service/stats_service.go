package service

import (
	"sort"

	"forum/internal/model"
)

// ForumStats 论坛全局统计。
type ForumStats struct {
	UserCount   int `json:"user_count"`
	BoardCount  int `json:"board_count"`
	TopicCount  int `json:"topic_count"`
	ReplyCount  int `json:"reply_count"`
	TotalViews  int `json:"total_views"`
}

// Stats 汇总全局统计。
func (s *Service) Stats() (*ForumStats, error) {
	stats := &ForumStats{
		UserCount:  len(s.store.ListUsers()),
		BoardCount: len(s.store.ListBoards()),
	}
	for _, t := range s.store.ListTopics() {
		stats.TopicCount++
		stats.TotalViews += t.ViewCount
		stats.ReplyCount += t.ReplyCount
	}
	return stats, nil
}

// BoardStats 单版块统计。
type BoardStats struct {
	BoardID    string `json:"board_id"`
	BoardName  string `json:"board_name"`
	TopicCount int    `json:"topic_count"`
	ReplyCount int    `json:"reply_count"`
}

// BoardStatsList 按版块统计主题与回帖数。
func (s *Service) BoardStatsList() ([]*BoardStats, error) {
	boards := s.store.ListBoards()
	nameByID := make(map[string]string, len(boards))
	statsByID := make(map[string]*BoardStats, len(boards))
	for _, b := range boards {
		nameByID[b.ID] = b.Name
		statsByID[b.ID] = &BoardStats{BoardID: b.ID, BoardName: b.Name}
	}
	for _, t := range s.store.ListTopics() {
		if st, ok := statsByID[t.BoardID]; ok {
			st.TopicCount++
			st.ReplyCount += t.ReplyCount
		}
	}
	result := make([]*BoardStats, 0, len(statsByID))
	for _, st := range statsByID {
		result = append(result, st)
	}
	sort.Slice(result, func(i, j int) bool {
		return result[i].TopicCount > result[j].TopicCount
	})
	return result, nil
}

// HotTopics 返回最热门的 N 个主题（按浏览数倒序）。
func (s *Service) HotTopics(limit int) ([]*model.Topic, error) {
	topics := s.store.ListTopics()
	sort.Slice(topics, func(i, j int) bool {
		return topics[i].ViewCount > topics[j].ViewCount
	})
	if limit <= 0 || limit > len(topics) {
		limit = len(topics)
	}
	return topics[:limit], nil
}
