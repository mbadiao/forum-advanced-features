function showSignup() {
    document.getElementById('login-left-content').style.display = 'none';
    document.getElementById('Signup-left-content').style.display = 'flex';
}

function showSignin() {
    document.getElementById('Signup-left-content').style.display = 'none';
    document.getElementById('login-left-content').style.display = 'flex';
}


const modalContainer = document.querySelector(".modal-container");
const modalTrigger = document.querySelectorAll(".modal-trigger");
const commentContainer = document.querySelector(".main-comment");
const commentTrigger = document.querySelectorAll(".comment-trigger");

commentTrigger.forEach(trigger => trigger.addEventListener("click",commentModal));
modalTrigger.forEach(trigger => trigger.addEventListener("click",toggleModal));
function commentModal(){
    commentContainer.classList.toggle("active1");
}
function toggleModal(){
    modalContainer.classList.toggle("active");
}

// Fonction pour stocker la position du défilement de la div centerside
function saveScrollPosition() {
    var centerside = document.getElementById('centerside');
    localStorage.setItem('centersideScrollPosition', centerside.scrollTop);
}

// Fonction pour charger et appliquer la position de défilement de la div centerside
function restoreScrollPosition() {
    var centerside = document.getElementById('centerside');
    var scrollPosition = localStorage.getItem('centersideScrollPosition');
    if (scrollPosition) {
        centerside.scrollTop = scrollPosition;
    }
}

// Écouteur d'événement pour enregistrer la position de défilement lors du défilement
document.getElementById('centerside').addEventListener('scroll', saveScrollPosition);

// Appel de la fonction pour restaurer la position de défilement lors du chargement de la page
window.addEventListener('load', restoreScrollPosition);

document.getElementById('file-input').addEventListener('change', function() {
    document.querySelector('.Photo').style.border = '2px solid #519e7a';
});


// .....................................................................................................


var textElements = document.getElementsByClassName('voirplus');

for (var i =  0; i < textElements.length; i++) {
    var textElement = textElements[i];
    var maxLength =  500;
    var originalText = textElement.textContent;
    var isExpanded = false; // Cette variable doit être réinitialisée pour chaque élément

    if (originalText.length > maxLength) {
        textElement.textContent = originalText.slice(0, maxLength) + '... ';
        var moreLink = document.createElement('a');
        moreLink.href = '#';
        moreLink.textContent = 'Voir plus';
        moreLink.style.color = '#87B29E'; // Ajout de la couleur au lien
        moreLink.onclick = (function (originalText, isExpanded) {
            return function () {
                if (isExpanded) {
                    this.parentNode.textContent = originalText.slice(0, maxLength) + '... ';
                    this.textContent = 'Voir plus';
                } else {
                    this.parentNode.textContent = originalText;
                    this.textContent = 'Voir moins';
                }
                // Basculer l'état de isExpanded
                isExpanded = !isExpanded;
                return false;
            };
        })(originalText, isExpanded);
        textElement.appendChild(moreLink);
    }
}

// Modifier le deuxième morceau de code en renommant les fonctions
var modalNotif = document.querySelector(".modalnotif");
var triggersNotif = document.querySelectorAll(".triggernotif");
var closeButtonNotif = document.querySelector(".close-buttonnotif");

function toggleModalNotif() {
  modalNotif.classList.toggle("show-modalnotif");
}

function windowOnClickNotif(event) {
  if (event.target === modalNotif) {
    toggleModalNotif();
  }
}

for (var i = 0, len = triggersNotif.length; i < len; i++) {
  triggersNotif[i].addEventListener("click", toggleModalNotif);
}
closeButtonNotif.addEventListener("click", toggleModalNotif);
window.addEventListener("click", windowOnClickNotif);


function view(event) {
    event.currentTarget.nextElementSibling.classList.toggle("style");
    event.currentTarget.nextElementSibling.classList.toggle("active");
}