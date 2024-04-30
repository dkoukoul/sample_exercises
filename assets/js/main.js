function speak(text) {
    var msg = new SpeechSynthesisUtterance(" "+text);
    msg.lang = 'el-GR';
    window.speechSynthesis.speak(msg);
}
    
function submitAnswer(question, answer) {
    fetch('/answer', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/x-www-form-urlencoded',
        },
        body: new URLSearchParams({
            'question': question,
            'answer': answer
        })
    }).then(response => {
        if (!response.ok) {
            throw new Error('Network response was not ok');
        }
        return response.json();
    }).then(data => {
        // Handle the response data here
    }).catch(error => {
        console.error('There has been a problem with your fetch operation:', error);
    });
}

function hideExercises() {
    let selectedIndex = 3;
    let templateDivs = document.querySelectorAll('.template');
    for(let i = 0; i < templateDivs.length; i++) {
        if(i === selectedIndex) {
            templateDivs[i].style.display = 'block';
        } else {
            templateDivs[i].style.display = 'none';
        }
    }
}

hideExercises();

function nextExercise() {
    console.log('Next exercise');
    let templateDivs = document.querySelectorAll('.template');
    let currentTemplate = 0;
    for(let i = 0; i < templateDivs.length; i++) {
        if (templateDivs[i].style.display === 'block') {
            currentTemplate = i;
            console.log(currentTemplate);
            break;
        }
    }
    templateDivs[currentTemplate].style.display = 'none';
    templateDivs[(currentTemplate + 1) % templateDivs.length].style.display = 'block';
}

let currentQuestion = 0;
function updateQuestion(back = false, templateClass) {
    let questions = document.querySelectorAll(templateClass);
    console.log(`Template: ${templateClass}`);
    let curQ, nexQ = 0;
    if (back) {
        curQ = currentQuestion + 1;
        nexQ = currentQuestion;
    } else {
        curQ = currentQuestion - 1;
        nexQ = currentQuestion;
    }
    console.log(`Questions: ${questions.length} next:${nexQ}`);
    let currQuestion = document.querySelector(`${templateClass}.i_${curQ}`);
    let nextQuestion = document.querySelector(`${templateClass}.i_${nexQ}`);

    //Hide the current question
    if (back) {
        currQuestion.style.animation = 'slideOutToRight 0.6s forwards';
        setTimeout(() => {
            nextQuestion.style.display = 'block';
            nextQuestion.style.animation = 'slideInFromLeft 0.6s backwards';
        }, 600);
    } else {
        currQuestion.style.animation = 'slideOutToLeft 0.6s forwards';
        setTimeout(() => {
            nextQuestion.style.display = 'block';
            nextQuestion.style.animation = 'slideInFromRight 0.6s backwards';
        }, 600);
    }
    
    setTimeout(() => {
        currQuestion.style.display = 'none';
    }, 600);
    
    //progressIndicator.textContent = `${currentQuestion + 1} / ${questions.length}`;
}

function nextQuestion(templateClass) {
    console.log(templateClass);
    let questions = document.querySelectorAll(templateClass);
    if (currentQuestion === questions.length - 1) {
        currentQuestion = 0;
        nextExercise();
        return;
    }
    currentQuestion = (currentQuestion + 1) % questions.length;
    updateQuestion(false, templateClass);
}

function selectAnswer(element) {
    let parent = element.parentNode;
    let sameClassElements = parent.getElementsByClassName(element.className);
    for (let i = 0; i < sameClassElements.length; i++) {
        sameClassElements[i].classList.remove('selected');
    }
    element.classList.add('selected');
}

function drag(ev) {
    ev.dataTransfer.setData("text", ev.target.id);
}

function allowDrop(ev) {
    ev.preventDefault();
}

function drop(ev) {
    ev.preventDefault();
    var data = ev.dataTransfer.getData("text"); 
    ev.target.appendChild(document.getElementById(data));    
}

// Add event listeners to the dropzones
var dropzones = document.querySelectorAll('.dropzone');
dropzones.forEach(function(dropzone) {
    dropzone.addEventListener('dragover', allowDrop);
    dropzone.addEventListener('drop', drop);
});