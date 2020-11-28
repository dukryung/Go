package main

import "fmt"

type SpoonOfJam interface {
	String() string
}

type Jam interface {
	GetOneSpoon() SpoonOfJam
}
type Bread struct {
	val string
}

type StrawberryJam struct {
}
type OrangeJam struct {
}
type AppleJam struct {
}
type SpoonOfStrawberryJam struct {
}
type SpoonOfOrangeJam struct {
}
type SpoonOfAppleJam struct {
}

func (s *SpoonOfStrawberryJam) String() string {
	return "+ strawberry"
}

func (j *StrawberryJam) GetOneSpoon() SpoonOfJam {
	return &SpoonOfStrawberryJam{}
}
func (j *OrangeJam) GetOneSpoon() SpoonOfJam {
	return &SpoonOfOrangeJam{}
}

func (j *AppleJam) GetOneSpoon() SpoonOfJam {
	return &SpoonOfAppleJam{}
}
func (s *SpoonOfOrangeJam) String() string {
	return "+ Orange"
}

func (s *SpoonOfAppleJam) String() string {
	return "+ Apple"
}

func (b *Bread) PutJam(jam Jam) {
	spoon := jam.GetOneSpoon()
	b.val += spoon.String()
}

func (b *Bread) String() string {
	return "bread " + b.val
}

func main() {
	bread := &Bread{}
	//jam := &StrawberryJam{}
	//orangejam := &OrangeJam{}
	applejam := &AppleJam{}

	bread.PutJam(applejam)
	fmt.Println(bread)

}
