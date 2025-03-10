let chosenQuestions = [];

function copy(element) {
    copyTextToClipboard(
        document.getElementById(element.getAttribute("data-copy-target")).value,
        "Ссылка скопирована"
    )
}

function previewImage(input) {
    const preview = document.getElementById('image-preview');
    const placeholder = document.getElementById('image-placeholder');

    if (input.files && input.files[0]) {
        const reader = new FileReader();

        reader.onload = function (e) {
            preview.src = e.target.result;
            preview.classList.remove('hidden');
            placeholder.classList.add('hidden');
        }

        reader.readAsDataURL(input.files[0]);
    } else {
        preview.src = '';
        preview.classList.add('hidden');
        placeholder.classList.remove('hidden');
    }
}

function showChoiceInput(element) {
    let id = element.getAttribute("data-id")
    document.getElementById("answer-choice-input-checkbox-" + id).checked = false
    document.getElementById("answer-choice-input-textarea-" + id).value = ""
    document.getElementById("answer-choice-input-" + id).classList.remove("hidden")
    document.getElementById("answer-choice-add-button-" + id).classList.add("hidden")
}

function showChoiceAddButton(element) {
    let id = element.getAttribute("data-id")
    document.getElementById("answer-choice-input-checkbox-" + id).checked = false
    document.getElementById("answer-choice-input-textarea-" + id).value = ""
    document.getElementById("answer-choice-add-button-" + id).classList.remove("hidden")
    document.getElementById("answer-choice-input-" + id).classList.add("hidden")
}

function selectQuestion(element) {
    let found = chosenQuestions.findIndex((i) => i === element.value)

    if (element.checked && found === -1) {
        chosenQuestions.push(element.value)
    }
    if (!element.checked && found !== -1) {
        chosenQuestions.splice(found, 1)
    }
}

function restoreSelectedQuestions() {
    chosenQuestions.forEach((id) => {
        let element = document.getElementById(id)
        if (element !== null) {
            element.checked = true
        }
    })
}

function enrichRequestByQuestions(id) {
    let element = document.getElementById(id)
    if (element === null) {
        return
    }

    element.addEventListener("submit", function (event) {
        event.preventDefault()
        chosenQuestions.forEach((id) => {
            let input = document.createElement('input')
            input.setAttribute('name', "question")
            input.setAttribute('value', id)
            input.setAttribute('type', "hidden")
            this.appendChild(input)
        })
        this.submit();
    })
}

document.addEventListener('htmx:afterRequest', function (evt) {
    if (evt.detail.target.id === 'question-list-container') {
        questionListContainerListener(evt)
    }
});

function questionListContainerListener(evt) {
    restoreSelectedQuestions()
}