// init runs only once when the page loads
function init() {
    getLamps()
        .then(data => {
            if (data !== null) {
                const rooms = JSON.parse(data);
                for (let i = 0; i < rooms.length; i++) {
                    addLamp(i, rooms[i].lamp);
                    if ( rooms[i].state === true ) {
                        setSliderState(makeStringFancy(rooms[i].lamp), rooms[i].state)
                    }
                }
            }
        });

}

// getLamps get the lamps and their states from the database
function getLamps() {
    return fetch('http://backend:8088/api/getLamps')
        .then(response => {
            if (!response.ok) {
                throw new Error(`HTTP error! State: ${response.state}`);
            }
            return response.text();
        })
        .then(data => {
            return data;
        })
        .catch(error => {
            console.error(`Error making the GET request: ${error.message}`);
            return null;
        });
}

// makeStringFancy replace underscore with a space and make every first character uppercase
function makeStringFancy(str) {
    const words = str.split('_');
    const capitalizedWords = words.map(word => word.charAt(0).toUpperCase() + word.slice(1));
    return capitalizedWords.join(' ');
}

// addLamp add lamps to the frontend, triggered only when the page loads
function addLamp(id, lampName) {
    const cardDiv = document.createElement("div");
    cardDiv.className = "card";
    cardDiv.style.width = "auto";

    const cardBody = document.createElement("div");
    cardBody.className = "card-body";

    const switchContainer = document.createElement("div");
    switchContainer.className = "switch-container";

    const switchLabel = document.createElement("label");
    switchLabel.className = "switch";
    switchLabel.innerHTML = "<input type=\"checkbox\"><span class=\"slider\"></span>";

    const lampNameElement = document.createElement("h4");
    lampNameElement.className = 'card-title';
    lampNameElement.innerText = makeStringFancy(lampName);

    switchContainer.appendChild(lampNameElement);
    switchContainer.appendChild(switchLabel);

    cardBody.appendChild(switchContainer);
    cardDiv.appendChild(cardBody);

    document.getElementById("cardContainer").appendChild(cardDiv);

    const switchInput = cardDiv.querySelector('.switch input');
    switchInput.addEventListener('click', function() {
        sendNewState(lampName, switchInput.checked);
    });
}

// setSliderState set the state of the given toggle switch based on the given value
function setSliderState(cardTitle, newState) {
    const cardTitles = document.querySelectorAll('.card-title');
    cardTitles.forEach(titleElement => {
        if (titleElement.textContent === cardTitle) {
            const slider = titleElement.nextElementSibling.querySelector('.switch input');
            if (slider) {
                slider.checked = newState;

                const changeEvent = new Event('change', { bubbles: true });
                slider.dispatchEvent(changeEvent);
            }
        }
    });
}

// sendNewState send the new state of a switch, triggered by a toggle switch
function sendNewState(name, state) {
    const data = {
        lamp: name,
        state: state
    };
    const jsonData = JSON.stringify(data);

    const requestOptions = {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: jsonData
    };

    fetch('http://backend:8088/api/addRecord', requestOptions)
        .then(response => {
            if (!response.ok) {
                throw new Error(`HTTP error! State: ${response.state}`);
            }
            return response.json();
        })
        .then(data => {
            console.log(`Lamp ${makeStringFancy(data.lamp)} set to ${data.state}`);
        })
        .catch(error => {
            console.error('Error:', error);
        });

}