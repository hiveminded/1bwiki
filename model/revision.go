package model

type Revision struct {
	ID        int64
	PageTitle string
	TextID    int64
	Comment   string
	UserID    int64
	UserText  string
	Minor     bool
	Deleted   bool
	Len       int
	ParentID  int
	TimeStamp int64
	LenDiff   int
}

func (r Revision) Verify() error {
	if r.PageTitle == "" || r.UserText == "" || r.Len < 1 || r.TimeStamp < 1 {
		return logger.Error("Invalid revision", "pageTitle", r.PageTitle, "userText", r.UserText, "len", r.Len, "time", r.TimeStamp)
	}
	return nil
}

func GetRevisions() ([]*Revision, error) {
	var revs []*Revision
	err := db.Select(&revs, `SELECT * FROM revision ORDER BY id DESC`)
	if err != nil {
		return revs, logger.Error("Unable to select all revisions", "err", err)
	}
	return revs, nil
}

func GetPageRevisions(title string) ([]*Revision, error) {
	var revs []*Revision
	err := db.Select(&revs, `SELECT * FROM revision WHERE pagetitle=$1 ORDER BY id DESC`, title)
	if err != nil {
		return revs, logger.Error("Unable to get revisions for page", "page", title, "err", err)
	}
	return revs, nil
}
