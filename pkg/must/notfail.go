package must

// NotFail panics on error
func NotFail(err error) {
	if err != nil {
		panic(err.Error())
	}
}
