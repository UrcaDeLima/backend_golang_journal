package store

// Store ...
type Store interface {
	News() NewsRepository
	Header() HeaderRepository
	Article() ArticleRepository
	Post() PostRepository
	InnerDescription() InnerDescriptionRepository
	Image() ImageRepository
	Recommendation() RecommendationRepository
	Interaction() InteractionRepository
}
