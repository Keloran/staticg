package generate

import (
	"fmt"
)

func Pages(root string) (IndexContent, error) {
	ic := IndexContent{}

	pages, err := BlogPages(root)
	if err != nil {
		return ic, fmt.Errorf("blog err: %w", err)
	}
	ic.Blog = pages
	fmt.Printf("blogs gathered\n")

	pages, err = CurrentProjects(root)
	if err != nil {
		return ic, fmt.Errorf("current projects: %w", err)
	}
	ic.Current = pages
	fmt.Printf("current projects gathered\n")

	pages, err = PastProjects(root)
	if err != nil {
		return ic, fmt.Errorf("past projects: %w", err)
	}
	ic.Past = pages
	fmt.Printf("past projects gathered\n")

	pages, err = getOthers()
	if err != nil {
		return ic, fmt.Errorf("others: %w", err)
	}
	ic.Other = pages
	fmt.Printf("other links\n")

	return ic, nil
}
