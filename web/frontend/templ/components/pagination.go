package frontend_components

import (
	"github.com/a-h/templ"
	"math"
)

const (
	siblingCount = 1
	pageNumbers  = 20 + siblingCount*2
)

func Pagination(currentPage int64, totalCount int64, limit int64) templ.Component {
	totalPages := int64(math.Ceil(float64(totalCount) / float64(limit)))

	if pageNumbers >= totalPages {
		return PaginationComposition(pageButtons(1, totalPages, currentPage)...)
	}

	firstSibling := currentPage - siblingCount
	lastSibling := currentPage + siblingCount

	if firstSibling < 1 {
		lastSibling += 1 - firstSibling
		firstSibling = 1
	}
	if lastSibling > totalPages {
		firstSibling -= lastSibling - totalPages
		lastSibling = totalPages
	}

	if firstSibling == 1 {
		firstSibling = 2
	}
	if lastSibling == totalPages {
		lastSibling = totalPages - 1
	}

	buttons := make([]templ.Component, 0, pageNumbers)
	buttons = append(buttons, PageButton(1, currentPage == 1))
	if firstSibling-1 > 1 {
		buttons = append(buttons, Dots())
	}
	buttons = append(buttons, pageButtons(firstSibling, lastSibling, currentPage)...)
	if totalPages-lastSibling > 1 {
		buttons = append(buttons, Dots())
	}
	buttons = append(buttons, PageButton(int(totalPages), currentPage == totalPages))

	return PaginationComposition(buttons...)
}

func pageButtons(from int64, to int64, currentPage int64) []templ.Component {
	result := make([]templ.Component, 0, to-from+1)
	for i := from; i <= to; i++ {
		result = append(result, PageButton(int(i), i == currentPage))
	}
	return result
}
