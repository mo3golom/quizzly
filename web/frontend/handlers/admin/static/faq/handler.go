package faq

import (
	"github.com/a-h/templ"
	"net/http"
	staticFAQ "quizzly/web/frontend/templ/admin/static/faq"
	frontendComponents "quizzly/web/frontend/templ/components"

	frontend "quizzly/web/frontend/templ"
)

const (
	faqTitle = "FAQ"
)

type (
	StaticFAQHandler struct{}
)

func NewStaticFAQHandler() *StaticFAQHandler {
	return &StaticFAQHandler{}
}

func (h *StaticFAQHandler) Handle(_ http.ResponseWriter, _ *http.Request, _ struct{}) (templ.Component, error) {
	return frontend.AdminPageComponent(
		faqTitle,
		frontendComponents.Composition(
			frontendComponents.Header(faqTitle),
			frontendComponents.Anchor(
				"how-to-create-question",
				frontendComponents.Accordion(
					"Как создать вопрос?",
					staticFAQ.HowToCreateQuestion(),
				),
			),
			frontendComponents.Anchor(
				"how-to-create-game",
				frontendComponents.Accordion(
					"Как создать игру?",
					staticFAQ.HowToCreateGame(),
				),
			),
			frontendComponents.Anchor(
				"how-to-start-game",
				frontendComponents.Accordion(
					"Как начать игру?",
					staticFAQ.HowToStartGame(),
				),
			),
			frontendComponents.Anchor(
				"how-to-share-game",
				frontendComponents.Accordion(
					"Как поделиться игрой?",
					staticFAQ.HowToShareGame(),
				),
			),
			frontendComponents.Anchor(
				"how-to-explore-statistics",
				frontendComponents.Accordion(
					"Как узнать, сколько человек сыграло в игру?",
					staticFAQ.HowToExploreStatistics(),
				),
			),
			frontendComponents.Anchor(
				"what-about-end-game",
				frontendComponents.Accordion(
					"Что будет, если завершить игру?",
					staticFAQ.WhatAboutEndGame(),
				),
			),
		),
	), nil
}
