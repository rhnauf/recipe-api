package entity

/*
to keep it simple, the instruction will just be string to comply for the basic requirements
for the real world application might separate each of ingredient, instruction, testimonial
tab to different database table
*/
type Recipe struct {
	Id          int64
	Title       string
	Description string
	Instruction string
	Publish     *bool
}
