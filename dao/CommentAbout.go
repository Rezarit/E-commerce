package dao

import (
	"database/sql"
	"errors"
	"github.com/Rezarit/E-commerce/domain"
)

func GetCommentbyProductID(productID int) (*sql.Rows, error) {
	cmd := "SELECT * FROM comments WHERE product_id=?"
	rows, err := DB.Query(cmd, productID)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }() //此处忽略了错误信息

	return rows, nil
}

func IsPraised(commentID, userID int) (int, error) {
	var LikeStatu int
	cmd := "SELECT like_status FROM comment_likes WHERE user_id = ? AND product_id = ?"
	err := DB.QueryRow(cmd, commentID, userID).Scan(&LikeStatu)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, nil
		} else {
			return -1, err
		}
	}

	if LikeStatu == 1 {
		return 1, nil
	} else if LikeStatu == 2 {
		return 2, nil
	} else if LikeStatu == 0 {
		return 0, nil
	}

	return 0, nil
}

func InsertComment(comment domain.Comment, productID int) error {
	cmd := "INSERT INTO comments (post_id, publish_time, content, user_id, avatar, nickname, praise_count, product_id) VALUES (?, ?, ?, ?, ?, ?, ?, ?)"
	_, err := DB.Exec(cmd,
		comment.CommentID,
		comment.PublishTime,
		comment.Content,
		comment.UserID,
		comment.Avatar,
		comment.Nickname,
		comment.PraiseCount,
		productID)
	if err != nil {
		return err
	}
	return nil
}

func DeleteComment(commentID int) error {
	cmd := "DELETE FROM comments WHERE post_id=?"
	_, err := DB.Exec(cmd, commentID)
	if err != nil {
		return err
	}
	return nil
}

func UpdataComment(comment domain.Comment, commentID int) error {
	cmd := "UPDATE comments SET content=?,publish_time=? WHERE post_id=?"
	_, err := DB.Exec(cmd, comment.Content, comment.PublishTime, commentID)
	if err != nil {
		return err
	}
	return nil
}

func PraiseComment(comment domain.Comment) error {
	cmd := "INSERT INTO comments (user_id,post_id,like_status) VALUES (?, ?, ?)"
	_, err := DB.Exec(cmd, comment.UserID, comment.CommentID, comment.IsPraised)
	if err != nil {
		return err
	}
	return nil
}
