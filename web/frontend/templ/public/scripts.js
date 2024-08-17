let chosenAnswers = [];
let maxPossibleChosenAnswers = 1;
let timeouts = [];

function clearTimeouts() {
    timeouts.forEach(timeout => clearTimeout(timeout));
    timeouts = [];
}

function submitAnswer() {
    let overlay = document.getElementById("game-page-overlay");
    overlay.classList.remove("hidden");

    clearTimeouts()
    timeouts.push(setTimeout(() => {
        overlay.classList.add("opacity-100");
        overlay.classList.remove("opacity-0");
    }, 10));
}

function showAnswerResult() {
    let fullAnimationDuration = 1500
    let stepAnimationDuration = 300
    let skipButtonEnabled= false
    let result  = document.getElementById("game-page-answer-result");

    let readEstimation = document.getElementById("game-page-answer-read-estimation");
    if (readEstimation !== null) {
        let readEstimationValue = parseInt(readEstimation.value)

        if (readEstimationValue > fullAnimationDuration) {
            fullAnimationDuration = readEstimationValue
            skipButtonEnabled = true
        }

    }

    clearTimeouts()
    timeouts.push(setTimeout(() => {
        result.classList.add("opacity-100", "animate-pulse-fade-in");
        result.classList.remove("opacity-0");
    }, 10));

    hideAnswerResult(fullAnimationDuration, stepAnimationDuration);
    if (skipButtonEnabled) {
        document.getElementById("game-page-skip-answer").classList.remove("hidden")
    }
}

function hideAnswerResult(fullDuration, stepDuration) {
    let result  = document.getElementById("game-page-answer-result");
    let overlay = document.getElementById("game-page-overlay");

    timeouts.push(setTimeout(() => {
        result.classList.add("opacity-0");
        result.classList.remove("opacity-100", "animate-pulse-fade-in");

        overlay.classList.add("opacity-0");
        overlay.classList.remove("opacity-100");
    }, fullDuration - stepDuration));
    timeouts.push(setTimeout(() => {
        result.classList.add("hidden");
        overlay.classList.add("hidden");
    }, fullDuration));
}

function chooseAnswer(element) {
    let id = element.id
    let found = chosenAnswers.findIndex((i) => i === id)

    let hideSubmitButton = () => {
        document.getElementById("play-page-submit-button").disabled=true
        document.getElementById("play-page-submit").classList.add("hidden", "opacity-0")
        document.getElementById("play-page-submit").classList.remove("animate-fade-in-up", "sm:animate-fade-in")
    }
    let showSubmitButton = () => {
        document.getElementById("play-page-submit-button").disabled=false
        document.getElementById("play-page-submit").classList.remove("hidden", "opacity-0")
        document.getElementById("play-page-submit").classList.add("animate-fade-in-up", "sm:animate-fade-in")
    }

    if (found !== -1) {
        element.classList.remove("outline", "outline-4", "outline-green-500")
        document.getElementById("checkbox-"+id).checked=false
        chosenAnswers.splice(found, 1)

        if (chosenAnswers.length === 0) {
           hideSubmitButton()
        }
        return
    }

    if (maxPossibleChosenAnswers === 1 && chosenAnswers.length > 0) {
        chosenAnswers.forEach((i) => {
            document.getElementById(i).classList.remove("outline", "outline-4", "outline-green-500")
            document.getElementById("checkbox-"+i).checked=false
        });
        chosenAnswers = []
        hideSubmitButton()
    }

    element.classList.add("outline", "outline-4", "outline-green-500")
    document.getElementById("checkbox-"+id).checked=true
    showSubmitButton()
    chosenAnswers.push(id)
}

function writeAnswer(element) {
    let hideSubmitButton = () => {
        document.getElementById("play-page-submit-button").disabled=true
        document.getElementById("play-page-submit").classList.add("hidden", "opacity-0")
        document.getElementById("play-page-submit").classList.remove("animate-fade-in-up", "sm:animate-fade-in")
    }
    let showSubmitButton = () => {
        document.getElementById("play-page-submit-button").disabled=false
        document.getElementById("play-page-submit").classList.remove("hidden", "opacity-0")
        document.getElementById("play-page-submit").classList.add("animate-fade-in-up", "sm:animate-fade-in")
    }

    if (element.value.length > 0) {
        showSubmitButton()
        return
    }

    hideSubmitButton()
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
    window.location = "/game/play?id=" + gameId
}
