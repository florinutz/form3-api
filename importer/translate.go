package importer

import "form3/business"

// converter converts incoming model to business model
type converter struct {
	jsonEmployees []*Employee
	jsonGifts     []*Gift
	Categories    []*business.Category
	Employees     []*business.Employee
	Gifts         []*business.Gift
}

// NewConverter does everything in the constructor
func NewConverter(jsonEmployees []*Employee, jsonGifts []*Gift) (c *converter) {
	categories := getUniqueCategories(jsonEmployees, jsonGifts)
	c = &converter{
		jsonEmployees: jsonEmployees,
		jsonGifts:     jsonGifts,
		Categories:    categories,
		Employees:     getEmployees(jsonEmployees, categories),
		Gifts:         getGifts(jsonGifts, categories),
	}

	return c
}

func getEmployees(jsonEmployees []*Employee, uniqueCategories []*business.Category) (employees []*business.Employee) {
	for _, je := range jsonEmployees {
		employee := &business.Employee{
			Name:       je.Name,
			Categories: matchCategories(je.Categories, uniqueCategories),
		}
		employees = append(employees, employee)
	}
	return
}

func getGifts(jsonGifts []*Gift, uniqueCategories []*business.Category) (gifts []*business.Gift) {
	for _, jg := range jsonGifts {
		gift := &business.Gift{
			Name:       jg.Name,
			Categories: matchCategories(jg.Categories, uniqueCategories),
		}
		gifts = append(gifts, gift)
	}
	return
}

func matchCategories(jsonCategories []Category, uniqueCategories []*business.Category) (categories []*business.Category) {
	for _, jec := range jsonCategories {
		for _, bc := range uniqueCategories {
			if bc.Name == string(jec) {
				categories = append(categories, bc)
			}
		}
	}
	return
}

func getUniqueCategories(jsonEmployees []*Employee, jsonGifts []*Gift) (categories []*business.Category) {
	var names []string

	for _, e := range jsonEmployees {
		for _, c := range e.Categories {
			names = appendUnique(names, string(c))
		}
	}
	for _, g := range jsonGifts {
		for _, c := range g.Categories {
			names = appendUnique(names, string(c))
		}
	}

	for _, n := range names {
		categories = append(categories, &business.Category{Name: n})
	}

	return
}

func appendUnique(names []string, name string) []string {
	for _, n := range names {
		if n == name {
			return names
		}
	}
	return append(names, name)
}
