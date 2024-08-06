
function showChoiceInput(element) {
    let id = element.getAttribute("data-id")
    document.getElementById("answer-choice-input-checkbox-"+id).checked = false
    document.getElementById("answer-choice-input-textarea-"+id).value = ""
    document.getElementById("answer-choice-input-"+id).classList.remove("hidden")
    document.getElementById("answer-choice-add-button-"+id).classList.add("hidden")
}

function showChoiceAddButton(element) {
    let id = element.getAttribute("data-id")
    document.getElementById("answer-choice-input-checkbox-"+id).checked = false
    document.getElementById("answer-choice-input-textarea-"+id).value = ""
    document.getElementById("answer-choice-add-button-"+id).classList.remove("hidden")
    document.getElementById("answer-choice-input-"+id).classList.add("hidden")
}