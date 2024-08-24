document.addEventListener("DOMContentLoaded", function () {
    let users = [];
    let currentIndex = 0;

    fetch('/memories')
        .then(response => response.json())
        .then(data => {
            users = data;
            const currentUser = document.querySelector('#memory').dataset.user;
            currentIndex = users.indexOf(currentUser);
            
        })
        .catch(error => {
            console.error("Error loading memories:", error);
        });

    function loadMemory(index) {
        const memoryId = users[index];
        fetch(`/memories/${memoryId}.json`)
            .then(response => {
                if (!response.ok) {
                    throw new Error("Network response was not ok");
                }
                return response.json();
            })
            .then(memory => {
                const memoryElement = document.querySelector('#memory');
                memoryElement.style.opacity = 0;
                
                const memoryText = document.querySelector('#memory .text');
                const memoryUsername = document.querySelector('#memory .username')
                const memoryUserpic = document.querySelector('#memory .userpic')

                setTimeout(() => {
                    let textParagraphs = memory.text.split("\n");
                    memoryText.innerHTML = textParagraphs.map(line => `<p>${line}</p>`).join("");
                    memoryUsername.innerText = memoryId;
                    memoryUsername.setAttribute("href", "/~" + memoryId);
                    if (typeof(memory.image) === "undefined") {
                        memoryUserpic.setAttribute("src", "/static/img/empty.webp");
                    } else {
                        memoryUserpic.setAttribute("src", memory.image);
                    }
                    history.pushState(null, "", "/");
                }, 250);
                
                setTimeout(() => {
                    memoryElement.style.opacity = 1;
                }, 500);

                memoryElement.setAttribute('data-user', memoryId);
            })
            .catch(error => {
                console.error("Error loading memory:", error);
            });
    }
    
    const memoryBlock = document.getElementById("memory");

    let touchStartX = 0;
    let touchEndX = 0;

    function showNextMemory() {
        currentIndex = (currentIndex + 1) % users.length;
        loadMemory(currentIndex);
    }
    
    function showPreviousMemory() {
        currentIndex = (currentIndex - 1 + users.length) % users.length;
        loadMemory(currentIndex);
    }
    
    memoryBlock.addEventListener("touchstart", function(event) {
        touchStartX = event.changedTouches[0].screenX;
    });

    memoryBlock.addEventListener("touchend", function(event) {
        touchEndX = event.changedTouches[0].screenX;
        handleGesture();
    });

    function handleGesture() {
        if (touchEndX < touchStartX - 50) {
            showNextMemory();
            return;
        }
        
        if (touchEndX > touchStartX + 50) {
            showPreviousMemory();
        }
    }

    document.querySelector('.arrow.next').addEventListener('click', showNextMemory);

    document.querySelector('.arrow.previous').addEventListener('click', showPreviousMemory);
    
    const overlay = document.getElementById('overlay');
    const popups = document.querySelectorAll('.popup');
    const popupToggles = document.querySelectorAll('.popup-toggle');

    function showPopup(popupId) {
        overlay.style.display = 'flex'; // Показываем затемнение
        popups.forEach(popup => {
            if (popup.getAttribute('data-popup') === popupId) {
                popup.style.display = 'block'; // Показываем нужный поп-ап
            }
        });
    }

    function hidePopup() {
        overlay.style.display = 'none'; // Скрываем затемнение
        popups.forEach(popup => {
            popup.style.display = 'none'; // Скрываем все поп-апы
        });
    }

    popupToggles.forEach(toggle => {
        toggle.addEventListener('click', function() {
            const popupId = this.getAttribute('data-popup');
            showPopup(popupId);
        });
    });

    overlay.addEventListener('click', function(e) {
        if (e.target === overlay || e.target.classList.contains('close')) {
            hidePopup();
        }
    });
});
