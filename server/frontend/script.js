// init runs only once when the page loads
function init() {
    getDevices()
        .then(data => {
            if (data !== null) {
                const rooms = JSON.parse(data);
                for (let i = 0; i < rooms.length; i++) {
                    addDevice(i, rooms[i].device);
                    if ( rooms[i].state === true ) {
                        setSliderState(makeStringFancy(rooms[i].device), rooms[i].state)
                    }
                }
            }
        });

}

// getDevices get the devices and their states from the database
function getDevices() {
    return fetch('/api/getDevices')
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

// addDevice add devices to the frontend, triggered only when the page loads
function addDevice(id, deviceName) {
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

    const deviceNameElement = document.createElement("h4");
    deviceNameElement.className = 'card-title';
    deviceNameElement.innerText = makeStringFancy(deviceName);

    switchContainer.appendChild(deviceNameElement);
    switchContainer.appendChild(switchLabel);

    cardBody.appendChild(switchContainer);
    cardDiv.appendChild(cardBody);

    document.getElementById("cardContainer").appendChild(cardDiv);

    const switchInput = cardDiv.querySelector('.switch input');
    switchInput.addEventListener('click', function() {
        sendNewState(deviceName, switchInput.checked);
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
        device: name,
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

    fetch('/api/addRecord', requestOptions)
        .then(response => {
            if (!response.ok) {
                throw new Error(`HTTP error! State: ${response.state}`);
            }
            return response.json();
        })
        .then(data => {
            console.log(`Device ${makeStringFancy(data.device)} set to ${data.state}`);
        })
        .catch(error => {
            console.error('Error:', error);
        });

}
