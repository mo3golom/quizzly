let timeouts = [];
let playPageQuestionFormPromise = null;

function clearTimeouts() {
    timeouts.forEach(timeout => clearTimeout(timeout));
    timeouts = [];
}

function scrollToTop() {
    window.scrollTo(0, 0);
}

function initPlayPageQuestionForm() {
    const form = document.getElementById('play-page-question-form');
    if (form != null) {
        const submitBtn = document.getElementById('play-page-submit-button');
        const submit = document.getElementById('play-page-submit');
        form.addEventListener('change', function () {
            const hasChecked = form.querySelector('input[type="checkbox"]:checked, input[type="radio"]:checked');
            if (hasChecked) {
                submitBtn.disabled = false
                submit.classList.remove("hidden", "opacity-0")
                submit.classList.add("animate-fade-in-up")
            } else {
                submitBtn.disabled = true
                submit.classList.add("hidden", "opacity-0")
                submit.classList.remove("animate-fade-in-up")
            }
        });
    }
}

function beforeRequestPlayPageQuestionForm(event) {
    let overlay = document.getElementById("game-page-overlay");

    // Remove hidden and set initial opacity-0
    overlay.classList.remove("hidden");
    overlay.classList.add("opacity-0");

    // Force browser reflow to ensure transition works
    overlay.offsetHeight;

    // Add opacity-100 to trigger the transition
    overlay.classList.add("opacity-100");
    overlay.classList.remove("opacity-0");


    playPageQuestionFormPromise = new Promise((resolve) => {
        setTimeout(() => {
            resolve();
        }, 300);
    });
}

function afterRequestPlayPageQuestionForm(event) {
    if (playPageQuestionFormPromise) {
        event.preventDefault();
        playPageQuestionFormPromise
            .then(() => {
                htmx.swap(event.detail.target, event.detail.xhr.response, { swapStyle: 'outerHTML' });
            })
            .finally(() => {
                playPageQuestionFormPromise = null;
            });
    }
}

function showAnswerResult() {
    const DEFAULT_FULL_DURATION = 1500;
    const STEP_DURATION = 300;

    const readEstimation = document.getElementById("game-page-answer-read-estimation");
    const readDuration = readEstimation ? parseInt(readEstimation.value) : DEFAULT_FULL_DURATION;

    const fullAnimationDuration = Math.max(DEFAULT_FULL_DURATION, readDuration);

    clearTimeouts();
    hideAnswerResult(fullAnimationDuration, STEP_DURATION);
}

function hideAnswerResult(fullDuration, stepDuration) {
    let result = document.getElementById("game-page-answer-result");
    let overlay = document.getElementById("game-page-overlay");
    overlay.classList.remove("opacity-100");
    overlay.classList.add("hidden", "opacity-0");

    timeouts.push(setTimeout(() => {
        result.classList.add("opacity-0");
        result.classList.remove("opacity-100", "animate-pulse-fade-in");

    }, fullDuration - stepDuration));
    timeouts.push(setTimeout(() => {
        result.classList.add("hidden");

        let resultsLink = document.getElementById("game-page-results-link")
        if (resultsLink !== null) {
            window.location = resultsLink.value
        }
    }, fullDuration));
}

function fire() {
    let count = 200
    let defaults = {
        origin: { y: 0.9 }
    };
    confetti({
        ...defaults,
        ...{
            spread: 26,
            startVelocity: 55,
        },
        particleCount: Math.floor(count * 0.25)
    });
    confetti({
        ...defaults,
        ...{
            spread: 60,
        },
        particleCount: Math.floor(count * 0.2)
    });
    confetti({
        ...defaults,
        ...{
            spread: 100,
            decay: 0.91,
            scalar: 0.8
        },
        particleCount: Math.floor(count * 0.35)
    });
    confetti({
        ...defaults,
        ...{
            spread: 120,
            startVelocity: 25,
            decay: 0.92,
            scalar: 1.2
        },
        particleCount: Math.floor(count * 0.1)
    });
    confetti({
        ...defaults,
        ...{
            spread: 120,
            startVelocity: 45,
        },
        particleCount: Math.floor(count * 0.1)
    });
}

function connectToGame() {
    let gameId = document.getElementById("game-start-page-game-id").value;
    if (gameId === "") {
        return
    }
    window.location = "/game/" + gameId
}

function copyShareResultsBlock(element) {
    let additionalText = element.getAttribute("data-additional-text")
    let text = window.location.href
    if (additionalText !== null && additionalText !== "") {
        text = additionalText + "\n" + text
    }

    copyTextToClipboard(text, "Ссылка скопирована")
}