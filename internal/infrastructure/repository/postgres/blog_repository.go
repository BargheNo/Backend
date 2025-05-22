package repositoryimpl

import (
	"github.com/BargheNo/Backend/internal/domain/entity"
	repository "github.com/BargheNo/Backend/internal/domain/repository/postgres"
	"github.com/BargheNo/Backend/internal/infrastructure/database"
	"gorm.io/gorm"
)

type BlogRepository struct{}

func NewBlogRepository() *BlogRepository {
	return &BlogRepository{}
}

func (repo *BlogRepository) CreatePost(db database.Database, post *entity.Post) error {
	return db.GetDB().Create(post).Error
}

func (repo *BlogRepository) UpdatePost(db database.Database, post *entity.Post) error {
	return db.GetDB().Save(post).Error
}

func (repo *BlogRepository) GetCorporationPosts(db database.Database, corporationID uint, opts ...repository.QueryModifier) []entity.Post {
	var posts []entity.Post
	query := db.GetDB().Where("corporation_id = ?", corporationID)
	for _, opt := range opts {
		query = opt.Apply(query).(*gorm.DB)
	}
	result := query.Find(&posts)
	if result.Error != nil {
		panic(result.Error)
	}
	return posts
}

func (repo *BlogRepository) FindPostByID(db database.Database, postID uint) (*entity.Post, bool) {
	var post entity.Post
	result := db.GetDB().Where("id = ?", postID).First(&post)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, false
		}
		panic(result.Error)
	}
	return &post, true
}
