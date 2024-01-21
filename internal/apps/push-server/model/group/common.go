package group

type Paging struct {
	PageSize int `default:"1"`
	PageNum  int `default:"20"`
}

func VerifyPage(page *Paging) (limit int, offset int) {
	if page == nil {
		return 20, 0
	}

	if page.PageNum <= 0 {
		page.PageNum = 1
	}

	if page.PageSize <= 0 {
		page.PageSize = 20
	}

	return page.PageSize, (page.PageNum - 1) * page.PageSize
}
