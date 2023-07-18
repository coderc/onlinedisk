package user_handler

func checkUsername(username string) bool {
	usernameLen := len(username)
	if usernameLen < 6 || usernameLen > 20 {
		return false
	}

	return true
}

func checkPasswordInSignup(password, confirmPassword string) bool {
	if password != confirmPassword {
		return false
	}

	passwordLen := len(password)
	if passwordLen < 3 || passwordLen > 20 {
		return false
	}

	return true
}
