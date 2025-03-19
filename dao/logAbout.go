package dao

func SearchPassword(userID int) (password string, err error) {
	cmd := "SELECT password From users WHERE id=?"
	err = DB.QueryRow(cmd, userID).Scan(&password)
	if err != nil {
		return "", err
	}
	return password, nil
}
