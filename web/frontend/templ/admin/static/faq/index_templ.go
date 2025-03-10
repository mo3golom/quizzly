// Code generated by templ - DO NOT EDIT.

// templ: version: v0.2.663
package static_faq

//lint:file-ignore SA4006 This context is only used if a nested component is present.

import "github.com/a-h/templ"
import "context"
import "io"
import "bytes"

func HowToCreateQuestion() templ.Component {
	return templ.ComponentFunc(func(ctx context.Context, templ_7745c5c3_W io.Writer) (templ_7745c5c3_Err error) {
		templ_7745c5c3_Buffer, templ_7745c5c3_IsBuffer := templ_7745c5c3_W.(*bytes.Buffer)
		if !templ_7745c5c3_IsBuffer {
			templ_7745c5c3_Buffer = templ.GetBuffer()
			defer templ.ReleaseBuffer(templ_7745c5c3_Buffer)
		}
		ctx = templ.InitializeContext(ctx)
		templ_7745c5c3_Var1 := templ.GetChildren(ctx)
		if templ_7745c5c3_Var1 == nil {
			templ_7745c5c3_Var1 = templ.NopComponent
		}
		ctx = templ.ClearChildren(ctx)
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<div class=\"pt-4 leading-relaxed\"><p class=\"mb-3\">Чтобы создать новый вопрос, перейдите в раздел <a href=\"/admin/question/list\" target=\"_blank\" class=\"link link-primary no-underline\"><b>\"Список вопросов\"</b></a> и нажмите кнопку <a href=\"/admin/question/new\" target=\"_blank\" class=\"btn btn-success rounded-2xl btn-sm align-middle\"><svg xmlns=\"http://www.w3.org/2000/svg\" fill=\"none\" viewBox=\"0 0 24 24\" stroke-width=\"1.5\" stroke=\"currentColor\" class=\"size-6\"><path stroke-linecap=\"round\" stroke-linejoin=\"round\" d=\"M12 9v6m3-3H9m12 0a9 9 0 1 1-18 0 9 9 0 0 1 18 0Z\"></path></svg> Добавить новый вопрос</a></p><div class=\"mb-3\">Перед вами откроется форма для создания вопроса. Сначала нужно выбрать тип вопроса:<ul class=\"list-inside list-decimal\"><li class=\"pl-4\"><b>Один ответ</b> - у вопроса есть только один правильный ответ.</li><li class=\"pl-4\"><b>Несколько ответов</b> - в вопросе может быть несколько правильных ответов. Используйте селектор \"Считать ответ верным если\", чтобы определить, как засчитывать правильный ответ:<ul class=\"list-inside list-disc\"><li class=\"pl-6\"><i>\"Выбраны ВСЕ правильные варианты ответа\"</i> - ответ будет засчитан, если выбраны все правильные варианты.</li><li class=\"pl-6\"><i>\"Выбран ЛЮБОЙ из правильных вариантов ответа\"</i> - ответ считается верным, если выбран хотя бы один из правильных вариантов. (в случае выбора неверного ответа, весь ответ не засчитывается)</li></ul></li><li class=\"pl-4\"><b>Ввод слова</b> -  здесь нужно ввести ответ самостоятельно. Правильность ответа проверяется по точному совпадению с ответом, который вы укажете. Не учитывается регистр букв. Если нужно задать вопрос с пропуском слова, используйте символ \"_\" для указания места пропуска (по желанию).</li></ul></div><p>Каждый тип вопроса требует заполнения текста вопроса, также можно добавить изображение, если это необходимо. <br>Ввод вариантов ответа зависит от выбранного типа вопроса.</p></div>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		if !templ_7745c5c3_IsBuffer {
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteTo(templ_7745c5c3_W)
		}
		return templ_7745c5c3_Err
	})
}

func HowToCreateGame() templ.Component {
	return templ.ComponentFunc(func(ctx context.Context, templ_7745c5c3_W io.Writer) (templ_7745c5c3_Err error) {
		templ_7745c5c3_Buffer, templ_7745c5c3_IsBuffer := templ_7745c5c3_W.(*bytes.Buffer)
		if !templ_7745c5c3_IsBuffer {
			templ_7745c5c3_Buffer = templ.GetBuffer()
			defer templ.ReleaseBuffer(templ_7745c5c3_Buffer)
		}
		ctx = templ.InitializeContext(ctx)
		templ_7745c5c3_Var2 := templ.GetChildren(ctx)
		if templ_7745c5c3_Var2 == nil {
			templ_7745c5c3_Var2 = templ.NopComponent
		}
		ctx = templ.ClearChildren(ctx)
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<div class=\"pt-4 leading-relaxed\"><p class=\"mb-3\">Чтобы создать игру, у вас должен быть список заранее подготовленных вопросов. Если у вас еще нет списка вопросов, ознакомьтесь с руководством <a href=\"#how-to-create-question\" class=\"link link-primary no-underline\">\"Как создать вопрос?\"</a> для их создания.</p><p class=\"mb-3\">Если у вас уже есть готовые вопросы, в боковом меню нажмите <a href=\"/admin/game/new\" target=\"_blank\" class=\"btn btn-sm text-amber-500 bg-transparent hover:text-white hover:bg-amber-500 border-0 align-middle shadow-none\"><svg xmlns=\"http://www.w3.org/2000/svg\" fill=\"none\" viewBox=\"0 0 24 24\" stroke-width=\"1.5\" stroke=\"currentColor\" class=\"size-6\"><path stroke-linecap=\"round\" stroke-linejoin=\"round\" d=\"M12 9v6m3-3H9m12 0a9 9 0 1 1-18 0 9 9 0 0 1 18 0Z\"></path></svg> <span>Новая игра</span></a> или зайдите в раздел <a href=\"/admin/game/list\" target=\"_blank\" class=\"link link-primary no-underline\"><b>\"Список игр\"</b></a> и нажмите кнопку <a href=\"/admin/game/new\" target=\"_blank\" class=\"btn btn-sm bg-success hover:bg-green-600 border-0 text-white align-middle\"><svg xmlns=\"http://www.w3.org/2000/svg\" fill=\"none\" viewBox=\"0 0 24 24\" stroke-width=\"1.5\" stroke=\"currentColor\" class=\"size-6\"><path stroke-linecap=\"round\" stroke-linejoin=\"round\" d=\"M12 9v6m3-3H9m12 0a9 9 0 1 1-18 0 9 9 0 0 1 18 0Z\"></path></svg> <span>Создать новую игру</span></a></p><div>Вы попадете на форму создания новой игры. Здесь вам нужно выбрать вопросы, которые войдут в вашу игру. Придумайте название игры и настройте параметры, если это необходимо. Вот какие настройки сейчас доступны:<ul class=\"list-inside list-decimal\"><li class=\"pl-4\"><b>Перемешать вопросы</b> – если эта опция включена, каждый игрок будет видеть вопросы в случайном порядке. Если игрок не ответил на вопрос, он увидит его снова при следующем входе в игру.</li><li class=\"pl-4\"><b>Перемешать ответы</b> – при включении этой функции ответы на каждый вопрос будут перемешиваться для каждого игрока. Это помогает избежать запоминания игроками правильного порядка ответов.</li><li class=\"pl-4\"><b>Показывать правильный ответ в случае неудачи</b> – когда этот параметр включен, при неправильном ответе игрока на экран результатов будет выводиться правильный ответ. Обратите внимание, что кнопка \"играть снова\" всегда активна, так что игрок может запомнить правильные ответы и пройти викторину без ошибок во второй раз.</li></ul></div></div>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		if !templ_7745c5c3_IsBuffer {
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteTo(templ_7745c5c3_W)
		}
		return templ_7745c5c3_Err
	})
}

func HowToStartGame() templ.Component {
	return templ.ComponentFunc(func(ctx context.Context, templ_7745c5c3_W io.Writer) (templ_7745c5c3_Err error) {
		templ_7745c5c3_Buffer, templ_7745c5c3_IsBuffer := templ_7745c5c3_W.(*bytes.Buffer)
		if !templ_7745c5c3_IsBuffer {
			templ_7745c5c3_Buffer = templ.GetBuffer()
			defer templ.ReleaseBuffer(templ_7745c5c3_Buffer)
		}
		ctx = templ.InitializeContext(ctx)
		templ_7745c5c3_Var3 := templ.GetChildren(ctx)
		if templ_7745c5c3_Var3 == nil {
			templ_7745c5c3_Var3 = templ.NopComponent
		}
		ctx = templ.ClearChildren(ctx)
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<div class=\"pt-4 leading-relaxed\"><p class=\"mb-3\">Если вы уже создали игру (инструкции по созданию игры смотрите в разделе <a href=\"#how-to-create-game\" class=\"link link-primary no-underline\">\"Как создать игру?\"</a>), вы можете начать её, нажав на кнопку <button class=\"btn btn-sm bg-success hover:bg-green-600 border-0 text-white align-middle\"><svg xmlns=\"http://www.w3.org/2000/svg\" fill=\"none\" viewBox=\"0 0 24 24\" stroke-width=\"1.5\" stroke=\"currentColor\" class=\"size-6\"><path stroke-linecap=\"round\" stroke-linejoin=\"round\" d=\"M5.25 5.653c0-.856.917-1.398 1.667-.986l11.54 6.347a1.125 1.125 0 0 1 0 1.972l-11.54 6.347a1.125 1.125 0 0 1-1.667-.986V5.653Z\"></path></svg> <span>Начать</span></button> на странице игры.</p><p>Когда игра начнётся, вы сможете скопировать ссылку на неё, нажав <button class=\"btn btn-sm bg-amber-500 hover:bg-amber-600 text-white border-0 align-middle\">скопировать</button> рядом с адресом ссылки. После этого вы сможете поделиться этой ссылкой с друзьями, чтобы они тоже могли присоединиться к игре.</p></div>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		if !templ_7745c5c3_IsBuffer {
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteTo(templ_7745c5c3_W)
		}
		return templ_7745c5c3_Err
	})
}

func HowToShareGame() templ.Component {
	return templ.ComponentFunc(func(ctx context.Context, templ_7745c5c3_W io.Writer) (templ_7745c5c3_Err error) {
		templ_7745c5c3_Buffer, templ_7745c5c3_IsBuffer := templ_7745c5c3_W.(*bytes.Buffer)
		if !templ_7745c5c3_IsBuffer {
			templ_7745c5c3_Buffer = templ.GetBuffer()
			defer templ.ReleaseBuffer(templ_7745c5c3_Buffer)
		}
		ctx = templ.InitializeContext(ctx)
		templ_7745c5c3_Var4 := templ.GetChildren(ctx)
		if templ_7745c5c3_Var4 == nil {
			templ_7745c5c3_Var4 = templ.NopComponent
		}
		ctx = templ.ClearChildren(ctx)
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<div class=\"pt-4 leading-relaxed\">Чтобы поделиться игрой с друзьями, скопируйте ссылку на страницу игры и отправьте её им :)</div>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		if !templ_7745c5c3_IsBuffer {
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteTo(templ_7745c5c3_W)
		}
		return templ_7745c5c3_Err
	})
}

func HowToExploreStatistics() templ.Component {
	return templ.ComponentFunc(func(ctx context.Context, templ_7745c5c3_W io.Writer) (templ_7745c5c3_Err error) {
		templ_7745c5c3_Buffer, templ_7745c5c3_IsBuffer := templ_7745c5c3_W.(*bytes.Buffer)
		if !templ_7745c5c3_IsBuffer {
			templ_7745c5c3_Buffer = templ.GetBuffer()
			defer templ.ReleaseBuffer(templ_7745c5c3_Buffer)
		}
		ctx = templ.InitializeContext(ctx)
		templ_7745c5c3_Var5 := templ.GetChildren(ctx)
		if templ_7745c5c3_Var5 == nil {
			templ_7745c5c3_Var5 = templ.NopComponent
		}
		ctx = templ.ClearChildren(ctx)
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<div class=\"pt-4 leading-relaxed\">Чтобы узнать, сколько людей уже сыграли в игру, перейдите на страницу игры и откройте вкладку <b>\"Участники\"</b>. Там вы найдёте следующую информацию о каждом игроке:<ul class=\"list-inside list-decimal\"><li class=\"pl-4\"><b>Имя игрока</b> – сейчас это случайно сгенерированное имя, но в будущем планируется возможность вводить своё имя.</li><li class=\"pl-4\"><b>Процент прохождения игры</b> – количество правильных ответов, выраженное в процентах.</li><li class=\"pl-4\"><b>Дата старта</b> – день, когда игроку был показан первый вопрос.</li><li class=\"pl-4\"><b>Дата последнего ответа</b> – день, когда игрок дал последний ответ. Эту дату можно считать окончанием игры, если статус прохождения – <span class=\"badge bg-orange-500 text-white align-middle\">\"Завершено\"</span>.</li><li class=\"pl-4\"><b>Статус прохождения</b> – возможные статусы: <span class=\"badge bg-success text-white align-middle\">\"В процессе\"</span> – игрок ещё проходит игру; <span class=\"badge bg-orange-500 text-white align-middle\">\"Завершено\"</span> – игра пройдена.</li></ul></div>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		if !templ_7745c5c3_IsBuffer {
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteTo(templ_7745c5c3_W)
		}
		return templ_7745c5c3_Err
	})
}

func WhatAboutEndGame() templ.Component {
	return templ.ComponentFunc(func(ctx context.Context, templ_7745c5c3_W io.Writer) (templ_7745c5c3_Err error) {
		templ_7745c5c3_Buffer, templ_7745c5c3_IsBuffer := templ_7745c5c3_W.(*bytes.Buffer)
		if !templ_7745c5c3_IsBuffer {
			templ_7745c5c3_Buffer = templ.GetBuffer()
			defer templ.ReleaseBuffer(templ_7745c5c3_Buffer)
		}
		ctx = templ.InitializeContext(ctx)
		templ_7745c5c3_Var6 := templ.GetChildren(ctx)
		if templ_7745c5c3_Var6 == nil {
			templ_7745c5c3_Var6 = templ.NopComponent
		}
		ctx = templ.ClearChildren(ctx)
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<div class=\"pt-4 leading-relaxed\">Когда игра завершена, новые игроки не смогут присоединиться по ссылке и играть в неё. Однако, страницы с результатами участников останутся доступными и будут продолжать работать.</div>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		if !templ_7745c5c3_IsBuffer {
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteTo(templ_7745c5c3_W)
		}
		return templ_7745c5c3_Err
	})
}
