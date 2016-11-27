package main

import (
	"strings"
)

// Class 職業類別
type Class string

const (
	Ranger  Class = "遊俠"
	Warrior Class = "戰士"
	Magic   Class = "魔導士"
)

// Element 屬性
type Element string

const (
	Fire  Element = "火"
	Water Element = "水"
	Earth Element = "土"
	Wind  Element = "風"
	Light Element = "光"
	Dark  Element = "黯"
)

type Job struct {
	Name       string `json:"name" gorm:"column:name;primary_key"`
	Class      `json:"class" gorm:"column:class;index"`
	Elements   []Element `json:"elements" gorm:"-"`
	ElementsDB string    `json:"-" gorm:"column:elements"`
}

func (Job) TableName() string {
	return "jobs"
}

func (j *Job) BeforeSave() error {

	if j.Elements == nil {
		j.Elements = []Element{}
	}

	for _, el := range j.Elements {
		j.ElementsDB = "," + string(el)
	}

	j.ElementsDB = strings.TrimLeft(j.ElementsDB, ",")

	return nil
}

func (j *Job) AfterFind() error {
	els := []Element{}
	for _, el := range strings.Split(j.ElementsDB, ",") {
		els = append(els, Element(el))
	}

	j.Elements = els

	return nil
}

type Card struct {
	Name     string `json:"name" gorm:"column:name;primary_key"`
	Class    `json:"class" gorm:"column:class;index"`
	Element  `json:"element" gorm:"column:element;index"`
	Describe string `json:"describe" gorm:"column:describe;type:text"`
}

func (Card) TableName() string {
	return "cards"
}
