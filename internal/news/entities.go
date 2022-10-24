package news

import "time"

type Enclosure struct {
	Url    string `xml:"url,attr"`
	Length int64  `xml:"length,attr"`
	Type   string `xml:"type,attr"`
}

type Item struct {
	Title    string `xml:"title" json:"title"`
	Link     string `xml:"link" json:"link"`
	Desc     string `xml:"description" json:"description"`
	City     string `xml:"city" json:"city,omitempty"`
	Company  string `xml:"company" json:"company,omitempty"`
	Logo     string `xml:"logo" json:"logo,omitempty"`
	JobType  string `xml:"jobtype" json:"jobtype,omitempty"`
	Category string `xml:"category" json:"category,omitempty"`
	PubDate  string `xml:"pubDate" json:"pub_date"`
}

type ItemDb struct {
	Id      uint64    `db:"id"`
	Title   string    `db:"title"`
	Desc    string    `db:"description"`
	PubDate time.Time `db:"pub_date"`
	Link    string    `db:"link"`
}

type Channel struct {
	Title string `xml:"title"`
	Link  string `xml:"link"`
	Desc  string `xml:"description"`
	Items []Item `xml:"item"`
}

type Rss struct {
	Channel Channel `xml:"channel"`
}
