function fallbackCopyTextToClipboard(text, successMessage)  {
    var textArea = document.createElement("textarea");
    textArea.value = text;

    // Avoid scrolling to bottom
    textArea.style.top = "0";
    textArea.style.left = "0";
    textArea.style.position = "fixed";

    document.body.appendChild(textArea);
    textArea.focus();
    textArea.select();

    try {
        var successful = document.execCommand('copy');
        var msg = successful ? 'successful' : 'unsuccessful';

        if (successful) {
            addToast(successMessage, 'success')
        }
    } catch (err) {
        addToast(err.message, 'error')
        console.error('Fallback: Oops, unable to copy', err);
    }

    document.body.removeChild(textArea);
}

function copyTextToClipboard(text, successMessage) {
    if (!navigator.clipboard) {
        fallbackCopyTextToClipboard(text, successMessage);
        return;
    }
    navigator.clipboard.writeText(text).then(function() {
        addToast(successMessage, 'success')
    }, function(err) {
        addToast(err.message, 'error')
        console.error('Async: Could not copy text: ', err);
    });
}

function addToast(message, messageType) {
    let globalMessages = document.getElementById("global-messages")
    let id = uuidv4()
    globalMessages.innerHTML += `
                    <div id="`+id+`" class="alert alert-`+ messageType +` text-white flex">
                        <span>`+ message +`</span>
                        <button class="btn btn-ghost btn-xs btn-square" onclick="this.parentElement.remove()">
                            <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor" class="size-6">
                              <path stroke-linecap="round" stroke-linejoin="round" d="M6 18 18 6M6 6l12 12" />
                            </svg>
                        </button>
                    </div>`
    setTimeout(function() {
        document.getElementById(id).remove()
    }, 5000)
}


function uuidv4() {
    return "10000000-1000-4000-8000-100000000000".replace(/[018]/g, c =>
        (+c ^ crypto.getRandomValues(new Uint8Array(1))[0] & 15 >> +c / 4).toString(16)
    );
}

htmx.on("htmx:responseError", function (evt) {
    addToast(evt.detail.xhr.responseText, 'error')
})