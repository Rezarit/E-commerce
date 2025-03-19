package service

import (
	"database/sql"
	"github.com/Rezarit/E-commerce/dao"
	"github.com/Rezarit/E-commerce/domain"
)

func PutCommentMsg(rows *sql.Rows) ([]domain.Comment, error) {
	var comments []domain.Comment

	for rows.Next() {
		var comment domain.Comment

		err := rows.Scan(
			&comment.CommentID,
			&comment.PublishTime,
			&comment.Content,
			&comment.UserID,
			&comment.Avatar,
			&comment.Nickname,
			&comment.PraiseCount,
			&comment.ProductID)
		if err != nil {
			return comments, err
		}

		comment.IsPraised, err = dao.IsPraised(comment.CommentID, comment.UserID)
		if err != nil {
			return comments, err
		}

		comments = append(comments, comment)
	}
	return comments, nil
}
