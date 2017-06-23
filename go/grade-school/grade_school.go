package school

import (
	"sort"
	"strings"
)

const testVersion = 1

type Grade struct {
	no       int
	students []string
}

type School struct {
	grades []*Grade
}

func New() *School {
	return &School{make([]*Grade, 0, 0)}
}
func (s *School) Add(student string, sgrade int) {
	for _, grade := range s.grades {
		if grade.no == sgrade {
			grade.students = append(grade.students, student)
			return
		}
	}
	students := make([]string, 0, 0)
	students = append(students, student)
	newGrade := &Grade{sgrade, students}
	s.grades = append(s.grades, newGrade)
}

func (s *School) Grade(sgrade int) []string {
	for _, grade := range s.grades {
		if grade.no == sgrade {
			return grade.students
		}
	}
	return nil
}
func (s *School) Enrollment() []Grade {
	sort.Sort(gradeList(s.grades)) // sort grades
	copy := make([]Grade, len(s.grades), cap(s.grades))
	for i, grade := range s.grades {
		sort.Sort(studentList(grade.students)) // sort students
		copy[i] = *grade
	}
	return copy
}

// sort.Interface for School.grades
type gradeList []*Grade

func (grades gradeList) Len() int {
	return len(grades)
}

func (grades gradeList) Less(i, j int) bool {
	return grades[i].no < grades[j].no
}

func (grades gradeList) Swap(i, j int) {
	temp := grades[i]
	grades[i] = grades[j]
	grades[j] = temp
}

// sort.Interface for Grade.students
type studentList []string

func (students studentList) Len() int {
	return len(students)
}

func (students studentList) Less(i, j int) bool {
	if strings.Compare(students[i], students[j]) > 0 {
		return false
	}
	return true
}

func (students studentList) Swap(i, j int) {
	temp := students[i]
	students[i] = students[j]
	students[j] = temp
}
